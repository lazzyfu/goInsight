package views

import (
	"goInsight/internal/app/inspect/forms"
	"goInsight/internal/app/inspect/services"
	"goInsight/internal/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func SyntaxInspectView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)

	// 解析表单数据
	var form forms.SyntaxInspectForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	// 执行
	service := services.SyntaxInspectService{
		Form:     form,
		C:        c,
		Username: username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 成功响应
	response.Success(c, nil, "success")
}
