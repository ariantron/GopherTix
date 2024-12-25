package models

import (
	"github.com/google/uuid"
	autnmodels "gopher_tix/modules/authentication/models"
	"gopher_tix/packages/common/types"
)

type UserRole struct {
	types.BaseModel
	UserID  uuid.UUID       `gorm:"type:uuid" json:"user_id"`
	User    autnmodels.User `gorm:"foreignKey:UserID" json:"user"`
	RoleID  uuid.UUID       `gorm:"type:uuid" json:"role_id"`
	Role    Role            `gorm:"foreignKey:RoleID" json:"role"`
	GroupID uuid.UUID       `gorm:"type:uuid" json:"group_id"`
	Group   Group           `gorm:"foreignKey:GroupID" json:"group"`
}
