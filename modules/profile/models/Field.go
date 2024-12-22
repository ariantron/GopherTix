package prof_models

import "gopher_tix/packages/common/types"

type Field struct {
	types.BaseModel
	Name string `gorm:"type:varchar(50)" json:"name"`
}
