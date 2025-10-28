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
		v1.GET("environments", views.GetEnvironmentsView)
		v1.GET("instances", views.GetInstancesView)
		v1.POST("inspect-sql-syntax", views.SyntaxInspectView)
		v1.GET("schemas", views.GetSchemasView)
		v1.GET("users", views.GetUsersView)
		v1.POST("commit", views.CreateOrdersView)
		v1.GET("history", views.GetListView)
		v1.GET("history/:order_id", views.GetDetailView)
		v1.GET("approval/:order_id", views.GetOrderApprovalView)
		v1.PUT("operate/approve", views.ApproveView)
		v1.PUT("operate/feedback", views.FeedbackView)
		v1.PUT("operate/review", views.ReviewView)
		v1.PUT("operate/close", views.CloseView)
		v1.POST("hook", views.HookOrdersView)
		v1.POST("generate-tasks", views.GenerateTasksView)
		v1.GET("tasks/:order_id", views.GetTasksView)
		v1.GET("tasks/preview", views.PreviewTasksView)
		v1.POST("tasks/execute-single", views.ExecuteSingleTaskView)
		v1.POST("tasks/execute-all", views.ExecuteAllTaskView)
		v1.GET("download/exportfile/:task_id", views.DownloadExportFileView)
	}
}
