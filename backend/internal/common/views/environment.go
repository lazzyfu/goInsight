package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/common/forms"
	"github.com/lazzyfu/goinsight/internal/common/services"

	"github.com/gin-gonic/gin"
)

func AdminGetEnvironmentView(c *gin.Context) {
	var form *forms.AdminGetEnvironmentForm = &forms.AdminGetEnvironmentForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetEnvironmentServices{
			AdminGetEnvironmentForm: form,
			C:                       c,
		}
		returnData, total, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.PaginationSuccess(c, total, returnData)
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

func AdminCreateEnvironmentView(c *gin.Context) {
	var form *forms.AdminCreateEnvironmentForm = &forms.AdminCreateEnvironmentForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminCreateEnvironmentService{
			AdminCreateEnvironmentForm: form,
			C:                          c,
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

func AdminUpdateEnvironmentView(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var form *forms.AdminUpdateEnvironmentForm = &forms.AdminUpdateEnvironmentForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminUpdateEnvironmentService{
			AdminUpdateEnvironmentForm: form,
			C:                          c,
			ID:                         id,
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

func AdminDeleteEnvironmentView(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	service := services.AdminDeleteEnvironmentService{
		C:  c,
		ID: id,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}
