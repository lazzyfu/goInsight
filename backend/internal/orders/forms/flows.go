package forms

import (
	"github.com/lazzyfu/goinsight/pkg/pagination"
	"gorm.io/datatypes"
)

type AdminGetApprovalFlowsForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}

type AdminUpdateApprovalFlowsForm struct {
	Definition datatypes.JSON `form:"definition" json:"definition"`
	Name       string         `form:"name"  json:"name" binding:"required,min=3,max=256"`
}

type AdminCreateApprovalFlowsForm struct {
	Definition datatypes.JSON `form:"definition" json:"definition"`
	Name       string         `form:"name"  json:"name" binding:"required,min=3,max=256"`
}
