package requests

import "github.com/google/uuid"

type CreateTicketRequest struct {
	Title      string    `json:"title" validate:"required,min=1,max=50"`
	Text       string    `json:"text" validate:"required,min=1"`
	GroupID    uuid.UUID `json:"group_id" validate:"required"`
	CreatedBy  uuid.UUID `json:"created_by" validate:"required"`
	AssignedTo uuid.UUID `json:"assigned_to,omitempty"`
}
