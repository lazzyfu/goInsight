package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	"github.com/gin-gonic/gin"
)

// 语法审核
func InspectOrderSyntaxView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.InspectOrderSyntaxForm = &forms.InspectOrderSyntaxForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.InspectOrderSyntaxService{
			InspectOrderSyntaxForm: form,
			C:                      c,
			Username:               username,
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
