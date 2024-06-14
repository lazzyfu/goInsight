package forms

import (
	"goInsight/pkg/pagination"
)

type AdminInspectParamsForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}

type AdminUpdateInspectParamsForm struct {
	Params map[string]interface{} `form:"params" json:"params"`
	Remark string                 `form:"remark"  json:"remark" binding:"required,min=3,max=256"`
}
