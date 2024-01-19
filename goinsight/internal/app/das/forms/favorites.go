package forms

import (
	"goInsight/internal/pkg/pagination"
)

type GetFavoritesForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}

type CreateFavoritesForm struct {
	Title   string `form:"title"  json:"title" binding:"required"`
	Sqltext string `form:"sqltext" json:"sqltext" binding:"required" `
}

type UpdateFavoritesForm struct {
	Title   string `form:"title"  json:"title"`
	Sqltext string `form:"sqltext" json:"sqltext"`
}
