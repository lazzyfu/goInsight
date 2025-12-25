package routers

import (
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/middleware"

	"github.com/lazzyfu/goinsight/internal/orders/views"

	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	r.GET("/ws/:channel", views.WebSocketHandler)

	apiV1 := r.Group("/api/v1/orders")
	apiV1.Use(global.App.JWT.MiddlewareFunc())
	{
		RegisterApiRoutes(apiV1)
	}

	adminV1 := r.Group("/api/v1/admin/approval-flows")
	adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
	{
		RegisterAdminRoutes(adminV1)
	}
}
