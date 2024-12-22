package models

import "gopher_tix/packages/common/types"

type Permission struct {
	types.BaseModel
	Name string `gorm:"type:varchar(50)" json:"name"`
}
