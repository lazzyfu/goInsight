package router

import (
	"goInsight/global"
	"goInsight/internal/app/inspect/views"

	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	v1 := r.Group("/api/v1/inspect")
	v1.Use(global.App.JWT.MiddlewareFunc())
	{
		v1.GET("/check", views.SyntaxInspectView)
	}
}
