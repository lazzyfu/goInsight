package views

import (
	"goInsight/internal/app/orders/forms"
	"goInsight/internal/app/orders/services"
	"goInsight/internal/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 提交工单
func SyntaxInspectView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.SyntaxInspectForm = &forms.SyntaxInspectForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.SyntaxInspectService{
			SyntaxInspectForm: form,
			C:                 c,
			Username:          username,
		}
		returnData, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, returnData, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}
