package models

import (
	"github.com/google/uuid"
	"gopher_tix/packages/common/types"
	"net"
)

type Login struct {
	types.BaseModel
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User    User      `gorm:"foreignKey:UserID" json:"user"`
	IP      net.IP    `gorm:"type:inet;not null" json:"ip"`
	Succeed bool      `gorm:"not null" json:"succeed"`
}
