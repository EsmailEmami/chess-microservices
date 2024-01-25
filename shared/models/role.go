package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

var (
	ROLE_ROOT  = uuid.MustParse("83a3d49f-f161-4b72-be5c-6f3903ba18be")
	ROLE_ADMIN = uuid.MustParse("27d8272f-95fe-4e5a-8033-a1a731c05f23")
	ROLE_USER  = uuid.MustParse("8c252b4f-81bf-4d66-be10-96bed059a807")
)

type Role struct {
	Model

	Name        string          `gorm:"column:name"        json:"name"`
	Code        string          `gorm:"column:code"        json:"code"`
	IsSystem    bool            `gorm:"column:is_system"   json:"isSystem"`
	Permissions RolePermissions `gorm:"column:permissions" json:"permissions"`
}

func (Role) TableName() string {
	return "role"
}

func (model Role) Permitted(action string) bool {
	for _, p := range model.Permissions {
		if p == action {
			return true
		}
	}

	return false
}

type RolePermissions []string

func (p RolePermissions) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *RolePermissions) Scan(value interface{}) error {
	if value == nil {
		*j = nil
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
