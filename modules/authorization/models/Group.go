package models

import (
	"github.com/google/uuid"
	"gopher_tix/packages/common/types"
)

type Group struct {
	types.SoftDeleteModel
	Name          string     `gorm:"type:varchar(50)" json:"name"`
	ParentGroupID *uuid.UUID `gorm:"type:uuid" json:"parent_group_id"`
	ParentGroup   *Group     `gorm:"foreignKey:ParentGroupID" json:"parent_group"`
	SubGroups     []Group    `gorm:"foreignKey:ParentGroupID" json:"sub_groups"`
	UserRoles     []UserRole `gorm:"foreignKey:GroupID" json:"user_roles"`
}
