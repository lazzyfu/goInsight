package forms

import (
	"github.com/lazzyfu/goinsight/pkg/pagination"
)

type GetOrderListForm struct {
	PaginationQ  pagination.Pagination
	OnlyMyOrders int    `form:"only_my_orders"`
	Search       string `form:"search"`
	Progress     string `form:"progress" json:"progress"`
	Environment  int    `form:"environment" json:"environment" `
}
