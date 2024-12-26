package requests

import "github.com/google/uuid"

type GroupUpdateRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=100"`
	OwnerUserID *uuid.UUID `json:"owner_user_id" validate:"omitempty"`
}
