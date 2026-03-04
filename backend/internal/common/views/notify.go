package views

import (
	"github.com/lazzyfu/goinsight/internal/common/forms"
	"github.com/lazzyfu/goinsight/internal/common/services"
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/gin-gonic/gin"
)

func AdminGetNotifyConfigView(c *gin.Context) {
	service := services.AdminGetNotifyConfigService{C: c}
	data, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

func AdminUpdateNotifyConfigView(c *gin.Context) {
	var form forms.AdminUpdateNotifyConfigForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminUpdateNotifyConfigService{
		AdminUpdateNotifyConfigForm: &form,
		C:                           c,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func AdminTestNotifyConfigView(c *gin.Context) {
	var form forms.AdminTestNotifyConfigForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminTestNotifyConfigService{
		AdminTestNotifyConfigForm: &form,
		C:                         c,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}
