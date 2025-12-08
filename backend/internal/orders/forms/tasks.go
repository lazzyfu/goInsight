package forms

import "github.com/lazzyfu/goinsight/pkg/pagination"

type GenOrderTasksForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}

type GetTasksForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
	Progress    string `form:"progress" json:"progress"`
}

type PreviewTasksForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}

type ExecuteSingleTaskForm struct {
	ID      int64  `form:"id" json:"id" binding:"required"`
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}

type ExecuteBatchTasksForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}
