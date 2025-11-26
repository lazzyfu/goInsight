package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/orders/views"
)

func RegisterAdminRoutes(v1 *gin.RouterGroup) {
	v1.GET("", views.AdminGetApprovalFlowsView)
	v1.POST("", views.AdminCreateApprovalFlowsView)
	v1.PUT(":id", views.AdminUpdateApprovalFlowsView)
	v1.DELETE(":id", views.AdminDeleteApprovalFlowsView)
	v1.POST("/bind-users", views.AdminBindUsersToApprovalFlowView)
	v1.GET("/bind-users/:approval_id", views.AdminGetApprovalFlowUsersView)
}
