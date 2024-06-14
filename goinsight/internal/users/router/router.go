package router

import (
	"goInsight/global"
	"goInsight/internal/users/views"
	"goInsight/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(v1 *gin.RouterGroup) {
	admin := v1.Group("/admin")
	admin.Use(middleware.HasAdminPermission())

	// 用户
	admin.GET("/users", views.GetUsersView)
	admin.POST("/users", views.CreateUsersView)
	admin.PUT("/users/:uid", views.UpdateUsersView)
	admin.DELETE("/users/:uid", views.DeleteUsersView)
	admin.POST("/users/change/password", views.ChangeUsersPasswordView)
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

func Routers(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// app
		v1.GET("/app/title", views.GetAppTitleView)
		// user auth
		v1.POST("/user/login", middleware.OTPMiddleware(), global.App.JWT.LoginHandler)
		v1.POST("/user/logout", global.App.JWT.LogoutHandler)
		v1.GET("/user/otp-auth-url", views.GetOTPAuthURLView)
		v1.GET("/user/otp-auth-callback", views.GetOTPAuthCallbackView)
		v1.GET("/user/refresh_token", global.App.JWT.RefreshHandler)
	}
	// 下面接口需要认证
	v1.Use(global.App.JWT.MiddlewareFunc())
	{
		// 个人
		v1.GET("/user", views.GetUserInfoView)
		v1.PUT("/user/:uid", views.UpdateUserInfoView)
		v1.POST("/user/change/avatar", views.ChangeUserAvatarView)
		v1.POST("/user/change/password", views.ChangeUserPasswordView)

		AdminRoutes(v1)
	}
}
