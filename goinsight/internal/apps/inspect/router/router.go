package router

import (
	"goInsight/global"
	"goInsight/internal/apps/inspect/views"
	"goInsight/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(v1 *gin.RouterGroup) {
	admin := v1.Group("/admin/inspect")
	admin.Use(middleware.HasAdminPermission())

	admin.GET("/params", views.AdminGetInspectParamsView)
	admin.PUT("/params/:id", views.AdminUpdateInspectParamsView)
}

func Routers(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.Use(global.App.JWT.MiddlewareFunc())

	AdminRoutes(v1)
}
