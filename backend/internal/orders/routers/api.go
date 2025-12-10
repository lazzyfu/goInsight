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
	// 获取工单
	v1.POST("", views.CreateOrderView)
	v1.GET("", views.GetOrderListView)
	v1.GET(":order_id", views.GetOrderDetailView)
	// 获取审批流和日志
	v1.GET("approval/:order_id", views.GetOrderApprovalView)
	v1.GET("logs/:order_id", views.GetOrderLogsView)
	// 操作
	v1.PUT("approval", views.ApprovalOrderView)
	v1.PUT("claim", views.ClaimOrderView)
	v1.PUT("transfer", views.TransferOrderView)
	v1.PUT("revoke", views.RevokeOrderView)
	v1.PUT("complete", views.CompleteOrderView)
	v1.PUT("fail", views.FailOrderView)
	v1.PUT("review", views.ReviewOrderView)
	// 新增生成执行任务接口
	v1.POST("generate-tasks", views.GenOrderTasksView)
	// 获取执行任务列表
	v1.GET("tasks/:order_id", views.GetTasksView)
	// 执行任务
	v1.POST("tasks/execute", views.ExecuteTaskView)
	// 批量执行任务
	v1.POST("tasks/execute-batch", views.ExecuteBatchTasksView)

	v1.POST("hook", views.HookOrdersView)

	v1.GET("tasks/preview", views.PreviewTasksView)

	v1.GET("download/exportfile/:task_id", views.DownloadExportFileView)
}
