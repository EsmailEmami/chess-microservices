package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

var GlobalRoomID = uuid.MustParse("9b7af5f3-cc90-4127-96ca-1e8b32e8bb75")

type Room struct {
	models.Model

	Name        string      `gorm:"name" json:"name"`
	IsPrivate   bool        `gorm:"is_private" json:"isPrivate"`
	Users       []UserRoom  `gorm:"foreignKey:room_id;references:id;" json:"users"`
	Messages    []Message   `gorm:"foreignKey:room_id;references:id;" json:"messages"`
	Avatar      string      `gorm:"column:avatar" json:"avatar"`
	PinMessages PinMessages `gorm:"column:pin_messages" json:"pinMessages"`
}

func (Room) TableName() string {
	return "chat.room"
}

type PinMessage struct {
	MessageID uuid.UUID
	Content   string
	Type      string
	PinDate   time.Time
}

type PinMessages []PinMessage

func (p PinMessages) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *PinMessages) Scan(value interface{}) error {
	if value == nil {
		j = nil
		return nil
	}
	var bts []byte
	switch v := value.(type) {
	case []byte:
		bts = v
	case string:
		bts = []byte(v)
	case nil:
		*j = nil
		return nil
	}
	return json.Unmarshal(bts, &j)
}
