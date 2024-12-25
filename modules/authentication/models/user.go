package models

import (
	"gopher_tix/packages/common/types"
)

type User struct {
	types.SoftDeleteModel
	Name     string `gorm:"type:varchar(50);not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
}
