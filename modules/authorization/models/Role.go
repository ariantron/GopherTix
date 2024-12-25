package models

import (
	"gopher_tix/packages/common/types"
)

type Role struct {
	types.BaseModel
	Name string `gorm:"type:varchar(50)" json:"name"`
}
