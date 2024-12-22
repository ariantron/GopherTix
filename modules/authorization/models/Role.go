package models

import (
	"gopher_tix/packages/common/types"
)

type Role struct {
	types.BaseModel
	Name        string           `json:"name"`
	Permissions []RolePermission `gorm:"foreignKey:RoleID" json:"permissions"`
}
