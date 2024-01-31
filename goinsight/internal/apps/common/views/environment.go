/*
@Author  :   lazzyfu
@Desc    :
*/

package views

import (
	"goInsight/internal/apps/common/forms"
	"goInsight/internal/apps/common/services"
	"goInsight/internal/pkg/response"
	"strconv"

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
	id, _ := strconv.Atoi(c.Param("id"))
	var form *forms.AdminUpdateEnvironmentForm = &forms.AdminUpdateEnvironmentForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminUpdateEnvironmentService{
			AdminUpdateEnvironmentForm: form,
			C:                          c,
			ID:                         uint64(id),
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
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.AdminDeleteEnvironmentService{
		C:  c,
		ID: uint64(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}
