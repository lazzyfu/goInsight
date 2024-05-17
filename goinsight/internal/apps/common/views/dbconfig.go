/*
@Author  :   xff
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

func AdminGetDBConfigView(c *gin.Context) {
	var form *forms.AdminGetDBConfigForm = &forms.AdminGetDBConfigForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetDBConfigServices{
			AdminGetDBConfigForm: form,
			C:                    c,
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

func AdminCreateDBConfigView(c *gin.Context) {
	var form *forms.AdminCreateDBConfigForm = &forms.AdminCreateDBConfigForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminCreateDBConfigService{
			AdminCreateDBConfigForm: form,
			C:                       c,
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

func AdminUpdateDBConfigView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form *forms.AdminUpdateDBConfigForm = &forms.AdminUpdateDBConfigForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminUpdateDBConfigService{
			AdminUpdateDBConfigForm: form,
			C:                       c,
			ID:                      uint64(id),
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

func AdminDeleteDBConfigView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.AdminDeleteDBConfigService{
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
