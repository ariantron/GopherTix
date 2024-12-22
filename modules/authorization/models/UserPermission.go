package models

import (
	"github.com/google/uuid"
	autnmodels "gopher_tix/modules/authentication/models"
	"gopher_tix/packages/common/types"
)

type UserPermission struct {
	types.BaseModel
	UserID       uuid.UUID       `gorm:"type:uuid" json:"user_id"`
	User         autnmodels.User `gorm:"foreignKey:UserID" json:"user"`
	PermissionID uuid.UUID       `gorm:"type:uuid" json:"permission_id"`
	Permission   Permission      `gorm:"foreignKey:PermissionID" json:"permission"`
}
