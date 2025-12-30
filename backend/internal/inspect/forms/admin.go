package forms

import (
	"github.com/lazzyfu/goinsight/pkg/pagination"
)

type AdminGlobalInspectParamsForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}

type AdminUpdateGlobalInspectParamsForm struct {
	Title string `form:"title" json:"title" binding:"required,min=3,max=256"`
	Type  string `form:"type"  json:"type"  binding:"required,oneof=string number boolean"`
	Key   string `form:"key" json:"key" binding:"required,min=3,max=256"`
	Value string `form:"value" json:"value" binding:"required,min=1,max=256"`
}
