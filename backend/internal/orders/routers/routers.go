package routers

import (
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/middleware"

	"github.com/lazzyfu/goinsight/internal/orders/views"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(v1 *gin.RouterGroup) {
	admin := v1.Group("/admin")
	admin.Use(middleware.HasAdminPermission())

	admin.GET("/approval-flows", views.AdminGetApprovalFlowsView)
	admin.PUT("/approval-flows/:id", views.AdminUpdateApprovalFlowsView)
	admin.POST("/approval-flows", views.AdminCreateApprovalFlowsView)
	admin.DELETE("/approval-flows/:id", views.AdminDeleteApprovalFlowsView)
}

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
		// 获取工单操作记录
		v1.GET("logs/:order_id", views.GetOrderLogsView)
		// 审批
		v1.PUT("approval", views.ApprovalOrderView)
		// 认领
		v1.PUT("claim", views.ClaimOrderView)
		// 转交
		v1.PUT("transfer", views.TransferOrderView)
		// 撤销工单
		v1.PUT("revoke", views.RevokeOrderView)
		// 手动更新工单为已完成
		v1.PUT("complete", views.CompleteOrderView)
		// 手动更新工单为失败
		v1.PUT("fail", views.FailOrderView)
		// 复核
		v1.PUT("review", views.ReviewOrderView)

		v1.POST("hook", views.HookOrdersView)
		v1.POST("generate-tasks", views.GenerateTasksView)
		v1.GET("tasks/:order_id", views.GetTasksView)
		v1.GET("tasks/preview", views.PreviewTasksView)
		v1.POST("tasks/execute-single", views.ExecuteSingleTaskView)
		v1.POST("tasks/execute-all", views.ExecuteAllTaskView)
		v1.GET("download/exportfile/:task_id", views.DownloadExportFileView)

		AdminRoutes(v1)
	}
}
