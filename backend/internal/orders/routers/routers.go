package routers

import (
	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/internal/orders/views"

	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	r.GET("/ws/:channel", views.WebSocketHandler)

	v1 := r.Group("/api/v1/orders")
	v1.Use(global.App.JWT.MiddlewareFunc())
	{
		// 工单上下文选项查询
		v1.GET("environments", views.GetOrderEnvironmentsView)
		v1.GET("instances", views.GetOrderInstancesView)
		v1.GET("schemas", views.GetOrderSchemasView)
		v1.GET("users", views.GetOrderUsersView)
		// SQL 语法检查
		v1.POST("inspect-syntax", views.InspectOrderSyntaxView)
		// 工单主流程
		v1.POST("", views.CreateOrderView)
		v1.GET("", views.GetOrderListView)
		v1.GET(":order_id", views.GetOrderDetailView)
		// 获取工单审批状态
		v1.GET("approval/:order_id", views.GetOrderApprovalView)
		// 审批
		v1.PUT("approval", views.ApprovalOrderView)
		// 认领
		v1.PUT("claim", views.ClaimOrderView)
		// 关闭工单
		v1.PUT("close", views.CloseOrderView)

		v1.PUT("operate/feedback", views.FeedbackView)
		v1.PUT("operate/review", views.ReviewView)
		v1.POST("hook", views.HookOrdersView)
		v1.POST("generate-tasks", views.GenerateTasksView)
		v1.GET("tasks/:order_id", views.GetTasksView)
		v1.GET("tasks/preview", views.PreviewTasksView)
		v1.POST("tasks/execute-single", views.ExecuteSingleTaskView)
		v1.POST("tasks/execute-all", views.ExecuteAllTaskView)
		v1.GET("download/exportfile/:task_id", views.DownloadExportFileView)
	}
}
