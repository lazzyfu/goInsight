package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/users/views"
)

func RegisterAdminRoutes(admin *gin.RouterGroup) {
	// 用户
	admin.GET("/users", views.GetUsersView)
	admin.POST("/users", views.CreateUsersView)
	admin.PUT("/users/:uid", views.UpdateUsersView)
	admin.DELETE("/users/:uid", views.DeleteUsersView)
	admin.POST("/users/reset-password", views.ResetUsersPasswordView)
	// 角色
	admin.GET("/roles", views.GetRolesView)
	admin.POST("/roles", views.CreateRolesView)
	admin.PUT("/roles/:id", views.UpdateRolesView)
	admin.DELETE("/roles/:id", views.DeleteRolesView)
	// 组织
	admin.GET("/organizations", views.GetOrganizationsView)
	admin.PUT("/organizations", views.UpdateOrganizationsView)
	admin.DELETE("/organizations", views.DeleteOrganizationsView)
	admin.POST("/organizations/root-node", views.CreateRootOrganizationsView)
	admin.POST("/organizations/child-node", views.CreateChildOrganizationsView)
	admin.GET("/organizations/users", views.GetOrganizationsUsersView)
	admin.POST("/organizations/users", views.BindOrganizationsUsersView)
	admin.DELETE("/organizations/users", views.DeleteOrganizationsUsersView)
}
