package models

import (
	"github.com/google/uuid"
	autnmodels "gopher_tix/modules/authentication/models"
	autzmodels "gopher_tix/modules/authorization/models"
	"gopher_tix/packages/common/types"
)

type Ticket struct {
	types.BaseModel
	Title      string           `gorm:"type:varchar(50)" json:"title"`
	Text       string           `json:"text"`
	GroupID    uuid.UUID        `gorm:"type:uuid" json:"group_id"`
	Group      autzmodels.Group `gorm:"foreignKey:GroupID" json:"group"`
	CreatedBy  uuid.UUID        `gorm:"type:uuid" json:"created_by"`
	Creator    autnmodels.User  `gorm:"foreignKey:CreatedBy" json:"creator"`
	AssignedTo *uuid.UUID       `gorm:"type:uuid" json:"assigned_to"`
	Assignee   *autnmodels.User `gorm:"foreignKey:AssignedTo" json:"assignee"`
	Comments   []Comment        `gorm:"foreignKey:TicketID" json:"comments"`
}
