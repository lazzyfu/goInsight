package forms

import (
	"goInsight/pkg/pagination"
)

type GetListForm struct {
	PaginationQ  pagination.Pagination
	OnlyMyOrders int    `form:"only_my_orders"`
	Search       string `form:"search"`
	Progress     string `form:"progress" json:"progress"`
	Environment  int    `form:"environment" json:"environment" `
}

type GetOpLogsForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}
