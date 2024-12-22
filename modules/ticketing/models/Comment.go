package models

import (
	"github.com/google/uuid"
	autnmodels "gopher_tix/modules/authentication/models"
	"gopher_tix/packages/common/types"
)

type Comment struct {
	types.BaseModel
	TicketID    uuid.UUID       `gorm:"type:uuid" json:"ticket_id"`
	Ticket      Ticket          `gorm:"foreignKey:TicketID" json:"ticket"`
	CommentedBy uuid.UUID       `gorm:"type:uuid" json:"commented_by"`
	Commenter   autnmodels.User `gorm:"foreignKey:CommentedBy" json:"commenter"`
	Text        string          `json:"text"`
}
