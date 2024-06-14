package forms

import (
	"goInsight/pkg/pagination"
)

type AdminGetEnvironmentForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}

type AdminCreateEnvironmentForm struct {
	Name string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}

type AdminUpdateEnvironmentForm struct {
	Name string `form:"name"  json:"name" binding:"required,min=2,max=32"`
}
