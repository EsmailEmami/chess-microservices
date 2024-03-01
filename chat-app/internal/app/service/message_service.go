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

func (m *MessageService) Get(ctx context.Context, id uuid.UUID) (*appModels.MessageOutPutDto, error) {
	db := psql.DBContext(ctx)
	return m.getMessage(db, id)
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

	qry := m.messageQry(db).Where("m.room_id = ?", roomID)

	var totalCounts int64

	if err := qry.Count(&totalCounts).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	skip := int(totalCounts) - cacheMessagesCount
	if skip < 0 {
		skip = 0
	}

	if err := qry.Offset(skip).Limit(cacheMessagesCount).Find(&messages).Error; err != nil {
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

	if len(lastMessages) > cacheMessagesCount {
		lastMessages = lastMessages[len(lastMessages)-cacheMessagesCount:]
	}

	// cache the result
	if err := m.cache.Set(m.roomMessagesCacheKey(roomID), &lastMessages, lastMessagesCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return message, nil
}

func (m *MessageService) NewFileMessage(ctx context.Context, roomID, userID, messageID uuid.UUID, content string, fileType string) (*appModels.MessageOutPutDto, error) {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	dbMsg := models.Message{
		Content: content,
		RoomID:  roomID,
		Type:    fileType,
	}
	dbMsg.CreatedByID = &userID
	dbMsg.ID = messageID

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

	if len(lastMessages) > cacheMessagesCount {
		lastMessages = lastMessages[len(lastMessages)-cacheMessagesCount:]
	}

	// cache the result
	if err := m.cache.Set(m.roomMessagesCacheKey(roomID), &lastMessages, lastMessagesCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return message, nil
}

func (m *MessageService) EditMessage(ctx context.Context, id, roomID uuid.UUID, content string) (*appModels.MessageOutPutDto, error) {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	dbMsg, err := m.getMessage(db, id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&models.Message{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"content":   content,
		"is_edited": true,
	}).Error; err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	// delete message cache
	if err := m.cache.Delete(m.messagesCacheKey(dbMsg.ID)); err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	message, err := m.getMessage(tx, dbMsg.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//delete room cache
	if err := m.cache.Delete(m.roomMessagesCacheKey(roomID)); err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return message, nil
}

func (m *MessageService) DeleteMessage(ctx context.Context, id, roomID uuid.UUID) error {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	dbMsg, err := m.getMessage(db, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.Message{}).Where("id = ?", id).Delete(&models.Message{}).Error; err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	// delete message cache
	if err := m.cache.Delete(m.messagesCacheKey(dbMsg.ID)); err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	//delete room cache
	if err := m.cache.Delete(m.roomMessagesCacheKey(roomID)); err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return nil
}

func (m *MessageService) SeenMessage(ctx context.Context, id, roomID uuid.UUID) error {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	dbMsg, err := m.getMessage(db, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.Message{}).Where("room_id = ? AND is_seen = ? AND created_at <= ?", roomID, false, dbMsg.CreatedAt).UpdateColumn("is_seen", true).Error; err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	// delete message cache
	if err := m.cache.Delete(m.messagesCacheKey(dbMsg.ID)); err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	//delete room cache
	if err := m.cache.Delete(m.roomMessagesCacheKey(roomID)); err != nil {
		tx.Rollback()
		return errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return nil
}

func (m *MessageService) DeleteCache(id uuid.UUID) error {
	return m.cache.Delete(m.messagesCacheKey(id))
}

func (m *MessageService) DeleteRoomMessagesCache(roomID uuid.UUID) error {
	return m.cache.Delete(m.roomMessagesCacheKey(roomID))
}

func (m *MessageService) messageQry(db *gorm.DB) *gorm.DB {
	return db.Table("chat.message m").
		Joins("INNER JOIN public.user u ON u.id = m.created_by_id").
		Joins("LEFT JOIN chat.message rm ON rm.id = m.reply_to_id").
		Joins("LEFT JOIN public.user ur ON ur.id = rm.created_by_id").
		Order("m.created_at ASC").
		Where("m.deleted_at IS NULL").
		Select("m.id, m.created_at, m.content, m.is_pin, m.created_by_id as user_id, m.is_edited, m.is_seen, u.first_name, u.last_name, m.reply_to_id, rm.content as reply_content, ur.first_name as reply_first_name, ur.last_name as reply_last_name, m.type")
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

	// message not found, we cannot use First in top query
	if message.ID == uuid.Nil {
		return nil, errs.NotFoundErr()
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
