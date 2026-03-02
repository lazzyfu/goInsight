package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/users/views"
)

func RegisterAdminRoutes(v1 *gin.RouterGroup) {
	// 用户
	v1.GET("/users", views.GetUsersView)
	v1.POST("/users", views.CreateUsersView)
	v1.PUT("/users/:uid", views.UpdateUsersView)
	v1.DELETE("/users/:uid", views.DeleteUsersView)
	v1.GET("/users/:uid/organizations", views.GetUserOrganizationsView)
	v1.POST("/users/reset-password", views.ResetUsersPasswordView)
	// 角色
	v1.GET("/roles", views.GetRolesView)
	v1.POST("/roles", views.CreateRolesView)
	v1.PUT("/roles/:id", views.UpdateRolesView)
	v1.DELETE("/roles/:id", views.DeleteRolesView)
	// 组织
	v1.GET("/organizations", views.GetOrganizationsView)
	v1.PUT("/organizations", views.UpdateOrganizationsView)
	v1.DELETE("/organizations", views.DeleteOrganizationsView)
	v1.POST("/organizations/root-node", views.CreateRootOrganizationsView)
	v1.POST("/organizations/child-node", views.CreateChildOrganizationsView)
	v1.GET("/organizations/users", views.GetOrganizationsUsersView)
	v1.POST("/organizations/users", views.BindOrganizationsUsersView)
	v1.DELETE("/organizations/users", views.DeleteOrganizationsUsersView)
}
