package router

import (
	"goInsight/global"
	"goInsight/internal/app/common/views"
	"goInsight/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(v1 *gin.RouterGroup) {
	admin := v1.Group("/admin")
	admin.Use(middleware.HasAdminPermission())

	admin.GET("/environment", views.AdminGetEnvironmentView)
	admin.POST("/environment", views.AdminCreateEnvironmentView)
	admin.PUT("/environment/:id", views.AdminUpdateEnvironmentView)
	admin.DELETE("/environment/:id", views.AdminDeleteEnvironmentView)

	admin.GET("/dbconfig", views.AdminGetDBConfigView)
	admin.POST("/dbconfig", views.AdminCreateDBConfigView)
	admin.PUT("/dbconfig/:id", views.AdminUpdateDBConfigView)
	admin.DELETE("/dbconfig/:id", views.AdminDeleteDBConfigView)
}

func Routers(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.Use(global.App.JWT.MiddlewareFunc())

	AdminRoutes(v1)
}
