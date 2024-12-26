package requests

import "github.com/google/uuid"

type GroupCreateRequest struct {
	Name          string     `json:"name" validate:"required,min=2,max=55"`
	ParentGroupID *uuid.UUID `json:"parent_group_id" validate:"omitempty"`
	OwnerUserID   uuid.UUID  `json:"owner_user_id" validate:"required"`
}
