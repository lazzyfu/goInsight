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

type ExecuteTaskForm struct {
	TaskID  string  `form:"task_id" json:"task_id" binding:"required,uuid"`
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}

type ExecuteBatchTasksForm struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required,uuid"`
}
