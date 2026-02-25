package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/users/views"
)

func RegisterApiRoutes(v1 *gin.RouterGroup) {
	v1.GET("", views.GetUserInfoView)
	v1.PUT("/:uid", views.UpdateUserInfoView)
	v1.POST("/change/avatar", views.ChangeUserAvatarView)
	v1.POST("/change/password", views.ChangeUserPasswordView)
}
