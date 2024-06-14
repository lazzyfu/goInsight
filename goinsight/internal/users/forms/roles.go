package forms

import (
	"goInsight/pkg/pagination"
)

type GetRolesForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}

type CreateRolesForm struct {
	Name string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}

type UpdateRolesForm struct {
	Name string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}
