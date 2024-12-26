package requests

import (
	"github.com/google/uuid"
)

type RoleManagementRequest struct {
	UserID  uuid.UUID `json:"user_id" validate:"required"`
	GroupID uuid.UUID `json:"group_id" validate:"required"`
	RoleID  uint8     `json:"role_id" validate:"required,min=1,max=255"`
}
