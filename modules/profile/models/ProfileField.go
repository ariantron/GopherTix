package prof_models

import (
	"github.com/google/uuid"
	"gopher_tix/packages/common/types"
)

type ProfileField struct {
	types.BaseModel
	ProfileID uuid.UUID `gorm:"type:uuid" json:"profile_id"`
	Profile   Profile   `gorm:"foreignKey:ProfileID" json:"profile"`
	FieldID   uuid.UUID `gorm:"type:uuid" json:"field_id"`
	Field     Field     `gorm:"foreignKey:FieldID" json:"field"`
	Value     string    `json:"value"`
}
