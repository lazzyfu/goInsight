package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/das/views"
)

func RegisterAdminRoutes(admin *gin.RouterGroup) {
	admin.GET("/schemas/grant", views.AdminGetSchemasGrantView)
	admin.POST("/schemas/grant", views.AdminCreateSchemasGrantView)
	admin.DELETE("/schemas/grant/:id", views.AdminDeleteSchemasGrantView)

	admin.GET("/tables/grant", views.AdminGetTablesGrantView)
	admin.POST("/tables/grant", views.AdminCreateTablesGrantView)
	admin.DELETE("/tables/grant/:id", views.AdminDeleteTablesGrantView)

	admin.GET("/instances/list", views.AdminGetInstancesListView)
	admin.GET("/schemas/list", views.AdminGetSchemasListView)
	admin.GET("/tables/list", views.AdminGetTablesListView)
}
