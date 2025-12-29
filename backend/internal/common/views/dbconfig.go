package views

import (
	"strconv"

	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/common/forms"
	"github.com/lazzyfu/goinsight/internal/common/services"

	"github.com/gin-gonic/gin"
)

func AdminGetInstancesView(c *gin.Context) {
	var form *forms.AdminGetInstancesForm = &forms.AdminGetInstancesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetInstancesServices{
			AdminGetInstancesForm: form,
			C:                     c,
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

func AdminCreateInstancesView(c *gin.Context) {
	var form forms.AdminCreateInstancesForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminCreateInstancesService{
		AdminCreateInstancesForm: &form,
		C:                        c,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func AdminUpdateInstancesView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form forms.AdminUpdateInstancesForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminUpdateInstancesService{
		AdminUpdateInstancesForm: &form,
		C:                        c,
		ID:                       uint64(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func AdminDeleteInstances(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.AdminDeleteInstanceConfigsService{
		C:  c,
		ID: uint64(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func AdminGetInstanceInspectParamsView(c *gin.Context) {
	var form *forms.AdminGetInstanceInspectParamsForm = &forms.AdminGetInstanceInspectParamsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetInstanceInspectParamsService{
			AdminGetInstanceInspectParamsForm: form,
			C:                                 c,
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
func AdminCreateInstanceInspectParamsView(c *gin.Context) {
	var form forms.AdminCreateInstanceInspectParamsForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminCreateInstanceInspectParamsService{
		AdminCreateInstanceInspectParamsForm: &form,
		C:                                    c,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func AdminUpdateInstanceInspectParamsView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form forms.AdminUpdateInstanceInspectParamsForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminUpdateInstanceInspectParamsService{
		AdminUpdateInstanceInspectParamsForm: &form,
		C:                                    c,
		ID:                                   uint64(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func AdminDeleteInstanceInspectParamsView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.AdminDeleteInstanceInspectParamsService{
		C:  c,
		ID: uint64(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}
