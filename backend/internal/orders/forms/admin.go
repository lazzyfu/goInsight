package forms

import (
	"github.com/google/uuid"
	"github.com/lazzyfu/goinsight/pkg/pagination"
	"gorm.io/datatypes"
)

type AdminGetApprovalFlowUnboundUsersForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}

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

type AdminBindUsersToApprovalFlowForm struct {
	ApprovalID uuid.UUID `form:"approval_id"  json:"approval_id" binding:"required"`
	Users      []string  `form:"users"  json:"users" binding:"required"`
}

type AdminGetApprovalFlowUsersForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
	ApprovalID  string `form:"approval_id"  json:"approval_id" binding:"required"`
}
