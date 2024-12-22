package prof_models

import (
	"github.com/google/uuid"
	autnmodels "gopher_tix/modules/authentication/models"
	"gopher_tix/packages/common/types"
)

type Profile struct {
	types.BaseModel
	UserID uuid.UUID       `gorm:"type:uuid" json:"user_id"`
	User   autnmodels.User `gorm:"foreignKey:UserID" json:"user"`
	Fields []ProfileField  `gorm:"foreignKey:ProfileID" json:"fields"`
}
