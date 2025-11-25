package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/middleware"
)

func Routers(r *gin.Engine) {
	adminV1 := r.Group("/api/v1/admin/inspect")
	adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
	{
		RegisterAdminRoutes(adminV1)
	}
}
