package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/services"

	"github.com/gin-gonic/gin"
)

// 获取用户权限
func GetUserGrantsView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.UserGrantsForm = &forms.UserGrantsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetUserGrantsService{UserGrantsForm: form, C: c, Username: username}
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
