/*
@Time    :   2023/03/08 14:54:37
@Author  :   xff
*/

package forms

import "goInsight/internal/pkg/pagination"

type GetHistoryForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}
