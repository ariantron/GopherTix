package prof_models

import (
	"github.com/google/uuid"
	autnmodels "gopher_tix/modules/authentication/models"
	"gopher_tix/packages/common/types"
)

type Profile struct {
	types.BaseModel
	Name   string          `gorm:"type:varchar(50)" json:"name"`
	UserID uuid.UUID       `gorm:"type:uuid" json:"user_id"`
	User   autnmodels.User `gorm:"foreignKey:UserID" json:"user"`
	Fields []ProfileField  `gorm:"foreignKey:ProfileID" json:"fields"`
}
