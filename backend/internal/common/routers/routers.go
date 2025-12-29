package routers

import (
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/middleware"

	"github.com/lazzyfu/goinsight/internal/common/views"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(v1 *gin.RouterGroup) {
	admin := v1.Group("/admin")
	admin.Use(middleware.HasAdminPermission())
	// 环境
	admin.GET("/environments", views.AdminGetEnvironmentView)
	admin.POST("/environment", views.AdminCreateEnvironmentView)
	admin.PUT("/environment/:id", views.AdminUpdateEnvironmentView)
	admin.DELETE("/environment/:id", views.AdminDeleteEnvironmentView)
	// 实例配置
	admin.GET("/instances", views.AdminGetInstancesView)
	admin.POST("/instances", views.AdminCreateInstancesView)
	admin.PUT("/instances/:id", views.AdminUpdateInstancesView)
	admin.DELETE("/instances/:id", views.AdminDeleteInstances)
	admin.GET("/instances/inspect/params", views.AdminGetInstanceInspectParamsView)
	admin.POST("/instances/inspect/params", views.AdminCreateInstanceInspectParamsView)
	admin.PUT("/instances/inspect/params/:id", views.AdminUpdateInstanceInspectParamsView)
	admin.DELETE("/instances/inspect/params/:id", views.AdminDeleteInstanceInspectParamsView)
}

func Routers(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.Use(global.App.JWT.MiddlewareFunc())

	AdminRoutes(v1)
}
