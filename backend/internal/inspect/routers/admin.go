package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/inspect/views"
)

func RegisterAdminRoutes(v1 *gin.RouterGroup) {
	v1.GET("/params", views.AdminGetGlobalInspectParamsView)
	v1.PUT("/params/:id", views.AdminUpdateGlobalInspectParamsView)
}
