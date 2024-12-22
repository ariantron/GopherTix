package models

import (
	"gopher_tix/packages/common/types"
	"time"
)

type User struct {
	types.SoftDeleteModel
	Email      string     `gorm:"uniqueIndex;not null" json:"email"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}
