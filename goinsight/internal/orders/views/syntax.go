package views

import (
	"goInsight/internal/orders/forms"
	"goInsight/internal/orders/services"
	"goInsight/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 语法审核
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
