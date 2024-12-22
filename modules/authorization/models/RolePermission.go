package models

import (
	"github.com/google/uuid"
	"gopher_tix/packages/common/types"
)

type RolePermission struct {
	types.BaseModel
	RoleID       uuid.UUID  `gorm:"type:uuid" json:"role_id"`
	Role         Role       `gorm:"foreignKey:RoleID" json:"role"`
	PermissionID uuid.UUID  `gorm:"type:uuid" json:"permission_id"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}
