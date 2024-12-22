package models

import "gopher_tix/packages/common/types"

type Permission struct {
	types.BaseModel
	Name string `json:"name"`
}
