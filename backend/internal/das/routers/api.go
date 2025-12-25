package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/das/views"
)

func RegisterDasRoutes(das *gin.RouterGroup) {
	das.GET("environments", views.GetEnvironmentsView)
	das.GET("schemas", views.GetSchemasView)
	das.GET("schema/tables", views.GetTablesView)
	das.GET("schema/grants", views.GetUserGrantsView)
	das.POST("query/mysql", views.ExecuteMySQLQueryView)
	das.POST("execute/query/clickhouse", views.ExecuteClickHouseQueryView)
	das.GET("table-info", views.GetTableInfoView)
	das.GET("dbdict", views.GetDbDictView)
	das.GET("history", views.GetHistoryView)
	das.GET("favorites", views.GetFavoritesView)
	das.POST("favorites", views.CreateFavoritesView)
	das.PUT("favorites/:id", views.UpdateFavoritesView)
	das.DELETE("favorites/:id", views.DeleteFavoritesView)
}
