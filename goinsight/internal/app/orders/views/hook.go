package views

import (
	"goInsight/internal/app/orders/forms"
	"goInsight/internal/app/orders/services"
	"goInsight/internal/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// hook工单
func HookOrdersView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.HookOrdersForm = &forms.HookOrdersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.HookOrdersService{
			HookOrdersForm: form,
			C:              c,
			Username:       username,
		}
		err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, nil, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}
