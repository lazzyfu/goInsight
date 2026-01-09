package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/orders/views"
)

func RegisterApiRoutes(v1 *gin.RouterGroup) {
	v1.GET("environments", views.GetOrderEnvironmentsView)
	v1.GET("instances", views.GetOrderInstancesView)
	v1.GET("schemas", views.GetOrderSchemasView)
	v1.GET("users", views.GetOrderUsersView)
	// 语法检查
	v1.POST("inspect-syntax", views.InspectOrderSyntaxView)
	// 创建工单
	v1.POST("", views.CreateOrderView)
	// 获取工单
	v1.GET("", views.GetOrderListView)
	v1.GET(":order_id", views.GetOrderDetailView)
	// 获取审批流和日志
	v1.GET("approvals/:order_id", views.GetOrderApprovalsView)
	v1.GET("logs/:order_id", views.GetOrderLogsView)
	// 操作
	v1.PUT("actions/approval", views.ApprovalOrderView)
	v1.PUT("actions/claim", views.ClaimOrderView)
	v1.PUT("actions/transfer", views.TransferOrderView)
	v1.PUT("actions/revoke", views.RevokeOrderView)
	v1.PUT("actions/complete", views.CompleteOrderView)
	v1.PUT("actions/fail", views.FailOrderView)
	v1.PUT("actions/review", views.ReviewOrderView)
	// 新增生成执行任务接口
	v1.POST("tasks", views.GenOrderTasksView)
	// 获取执行任务列表
	v1.GET("tasks/:order_id", views.GetTasksView)
	// 执行任务
	v1.POST("tasks/execute", views.ExecuteTaskView)
	// 批量执行任务
	v1.POST("tasks/execute-batch", views.ExecuteBatchTasksView)
	// 预览
	v1.GET("tasks/preview", views.PreviewTasksView)
	// 下载导出文件
	v1.GET("tasks/exports/:task_id", views.DownloadExportFileView)
}
