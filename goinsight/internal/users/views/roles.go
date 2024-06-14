/*
@Time    :   2023/08/31 15:46:42
@Author  :   xff
@Desc    :   角色管理
*/

package views

import (
	"goInsight/internal/users/forms"
	"goInsight/internal/users/services"
	"goInsight/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRolesView(c *gin.Context) {
	var form *forms.GetRolesForm = &forms.GetRolesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetRolesServices{
			GetRolesForm: form,
			C:            c,
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

func CreateRolesView(c *gin.Context) {
	var form *forms.CreateRolesForm = &forms.CreateRolesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CreateRolesService{
			CreateRolesForm: form,
			C:               c,
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

func UpdateRolesView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form *forms.UpdateRolesForm = &forms.UpdateRolesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.UpdateRolesService{
			UpdateRolesForm: form,
			C:               c,
			ID:              uint64(id),
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

func DeleteRolesView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.DeleteRolesService{
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
