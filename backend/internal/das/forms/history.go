package forms

import "github.com/lazzyfu/goinsight/pkg/pagination"

type GetHistoryForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}
