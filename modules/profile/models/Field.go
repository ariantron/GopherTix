package prof_models

import "gopher_tix/packages/common/types"

type Field struct {
	types.BaseModel
	Name string `json:"name"`
}
