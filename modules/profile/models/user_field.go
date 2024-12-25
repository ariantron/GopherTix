package prof_models

import (
	"github.com/google/uuid"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/packages/common/types"
)

type UserField struct {
	types.BaseModel
	UserID  uuid.UUID   `gorm:"type:uuid" json:"user_id"`
	User    models.User `gorm:"foreignKey:UserID" json:"user"`
	FieldID uuid.UUID   `gorm:"type:uuid" json:"field_id"`
	Field   Field       `gorm:"foreignKey:FieldID" json:"field"`
	Value   string      `json:"value"`
}
