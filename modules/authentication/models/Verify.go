package models

import (
	"github.com/google/uuid"
	"gopher_tix/packages/common/types"
)

type Verify struct {
	types.BaseModel
	UserID uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
	Token  string    `json:"token"`
	Type   TokenType `json:"type"`
}

type TokenType string

const (
	VerifyAccount TokenType = "VERIFY_ACCOUNT"
	ResetPassword TokenType = "RESET_PASSWORD"
)
