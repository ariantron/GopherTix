package requests

import "github.com/google/uuid"

type LeaveCommentRequest struct {
	TicketID    uuid.UUID `json:"ticket_id" validate:"required"`
	CommentedBy uuid.UUID `json:"commented_by" validate:"required"`
	Text        string    `json:"text" validate:"required,min=1"`
}
