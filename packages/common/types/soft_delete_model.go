package types

import (
	"gorm.io/gorm"
)

type SoftDeleteModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
