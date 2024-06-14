package forms

import "goInsight/pkg/pagination"

type GenerateTasksForm struct {
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

type ExecuteAllTaskForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}
