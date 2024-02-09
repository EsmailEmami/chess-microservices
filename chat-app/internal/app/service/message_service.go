package service

import (
	"context"
	"time"

	appModels "github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/models"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	cacheMessagesCount        = 300
	lastMessagesCacheDuration = 1 * time.Hour
	messageCacheDuration      = 30 * time.Second
)

type MessageService struct {
	cache *redis.Redis
}

func NewMessageService(cache *redis.Redis) *MessageService {
	return &MessageService{
		cache: cache,
	}
}

func (m *MessageService) GetLastMessages(ctx context.Context, roomID uuid.UUID) ([]appModels.MessageOutPutDto, error) {
	db := psql.DBContext(ctx)

	return m.getLastMessages(ctx, db, roomID)
}

func (m *MessageService) getLastMessages(ctx context.Context, db *gorm.DB, roomID uuid.UUID) ([]appModels.MessageOutPutDto, error) {
	var messages []appModels.MessageOutPutDto

	// try get from cache
	if err := m.cache.UnmarshalToObject(m.roomMessagesCacheKey(roomID), &messages); err == nil {
		return messages, nil
	}

	qry := m.messageQry(db).Where("m.room_id = ?", roomID).Limit(cacheMessagesCount)

	if err := qry.Find(&messages).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	// cache the result
	if err := m.cache.Set(m.roomMessagesCacheKey(roomID), &messages, lastMessagesCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return messages, nil
}

func (m *MessageService) NewMessage(ctx context.Context, roomID, userID uuid.UUID, content string, replyTo *uuid.UUID) (*appModels.MessageOutPutDto, error) {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	dbMsg := models.Message{
		Content:   content,
		RoomID:    roomID,
		ReplyToID: replyTo,
	}
	dbMsg.CreatedByID = &userID
	dbMsg.ID = uuid.New()

	if err := tx.Create(&dbMsg).Error; err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	message, err := m.getMessage(tx, dbMsg.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	lastMessages, err := m.GetLastMessages(ctx, roomID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	lastMessages = append(lastMessages, *message)

	if len(lastMessages) > 300 {
		lastMessages = lastMessages[len(lastMessages)-300:]
	}

	// cache the result
	if err := m.cache.Set(m.roomMessagesCacheKey(roomID), &lastMessages, lastMessagesCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return message, nil
}

func (m *MessageService) EditMessage(ctx context.Context, id uuid.UUID, content string) (*appModels.MessageOutPutDto, error) {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	var dbMsg models.Message

	if err := db.Model(&models.Message{}).Find(&dbMsg, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return nil, errs.NotFoundErr().WithError(err)
	}
	dbMsg.Content = content

	if err := tx.Model(&models.Message{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"content":   content,
		"is_edited": true,
	}).Error; err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	message, err := m.getMessage(tx, dbMsg.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	lastMessages, err := m.getLastMessages(ctx, tx, dbMsg.RoomID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// cache the result
	if err := m.cache.Set(m.roomMessagesCacheKey(dbMsg.RoomID), &lastMessages, lastMessagesCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return message, nil
}

func (m *MessageService) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	var dbMsg models.Message

	if err := db.Model(&models.Message{}).Find(&dbMsg, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return errs.NotFoundErr().WithError(err)
	}

	if err := tx.Model(&models.Message{}).Where("id = ?", id).Delete(&dbMsg).Error; err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	lastMessages, err := m.getLastMessages(ctx, tx, dbMsg.RoomID)

	if err != nil {
		tx.Rollback()
		return err
	}

	// cache the result
	if err := m.cache.Set(m.roomMessagesCacheKey(dbMsg.RoomID), &lastMessages, lastMessagesCacheDuration); err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return nil
}

func (m *MessageService) SeenMessage(ctx context.Context, id uuid.UUID) error {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	var dbMsg models.Message

	if err := db.Model(&models.Message{}).Find(&dbMsg, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return errs.NotFoundErr().WithError(err)
	}

	if err := tx.Model(&models.Message{}).Where("room_id = ? AND is_seen = ? AND created_at <= ", dbMsg.RoomID, false, dbMsg.CreatedAt).UpdateColumn("is_seen", true).Error; err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	lastMessages, err := m.getLastMessages(ctx, tx, dbMsg.RoomID)

	if err != nil {
		tx.Rollback()
		return err
	}

	// cache the result
	if err := m.cache.Set(m.roomMessagesCacheKey(dbMsg.RoomID), &lastMessages, lastMessagesCacheDuration); err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return nil
}

func (m *MessageService) messageQry(db *gorm.DB) *gorm.DB {
	return db.Table("chat.message m").
		Joins("INNER JOIN public.user u ON u.id = m.created_by_id").
		Joins("LEFT JOIN chat.message rm ON rm.id = m.reply_to_id").
		Joins("LEFT JOIN public.user ur ON ur.id = rm.created_by_id").
		Order("m.created_at DESC").
		Select("m.id, m.created_at, m.content, m.created_by_id as user_id, m.is_edited, m.is_seen, u.first_name, u.last_name, m.reply_to_id, rm.content as reply_content, ur.first_name as reply_first_name, ur.last_name as reply_last_name")
}

func (m *MessageService) getMessage(db *gorm.DB, id uuid.UUID) (*appModels.MessageOutPutDto, error) {
	var message appModels.MessageOutPutDto

	// try get from cache
	if err := m.cache.UnmarshalToObject(m.messagesCacheKey(id), &message); err == nil {
		return &message, nil
	}

	if err := m.messageQry(db).Where("m.id = ?", id).Find(&message).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	if err := m.cache.Set(m.messagesCacheKey(id), &message, messageCacheDuration); err != nil {
		return nil, err
	}

	return &message, nil
}

func (m *MessageService) roomMessagesCacheKey(roomID uuid.UUID) string {
	return "chat_room_last_messages_" + roomID.String()
}

func (m *MessageService) messagesCacheKey(id uuid.UUID) string {
	return "chat_room_messages_" + id.String()
}
