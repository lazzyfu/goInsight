package router

import (
	"goInsight/global"
	"goInsight/internal/das/views"
	"goInsight/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(v1 *gin.RouterGroup) {
	admin := v1.Group("/admin")
	admin.Use(middleware.HasAdminPermission())

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

func Routers(r *gin.Engine) {
	v1 := r.Group("/api/v1/das")
	v1.Use(global.App.JWT.MiddlewareFunc())
	{
		v1.GET("environments", views.GetEnvironmentsView)
		v1.GET("schemas", views.GetSchemasView)
		v1.GET("tables", views.GetTablesView)
		v1.POST("execute/query/mysql", views.ExecuteMySQLQueryView)
		v1.POST("execute/query/clickhouse", views.ExecuteClickHouseQueryView)
		v1.GET("table-info", views.GetTableInfoView)
		v1.GET("dbdict", views.GetDbDictView)
		v1.GET("history", views.GetHistoryView)
		v1.GET("favorites", views.GetFavoritesView)
		v1.POST("favorites", views.CreateFavoritesView)
		v1.PUT("favorites/:id", views.UpdateFavoritesView)
		v1.DELETE("favorites/:id", views.DeleteFavoritesView)
		v1.GET("user/grants", views.GetUserGrantsView)

		AdminRoutes(v1)
	}
}
