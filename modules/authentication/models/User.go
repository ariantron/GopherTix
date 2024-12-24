package models

import (
	"gopher_tix/packages/common/types"
)

type User struct {
	types.SoftDeleteModel
	Email      string     `gorm:"uniqueIndex;not null" json:"email"`
	Password   string     `gorm:"not null" json:"-"`
}
