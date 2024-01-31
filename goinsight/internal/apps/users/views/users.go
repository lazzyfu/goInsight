/*
@Time    :   2023/08/30 15:46:42
@Author  :   lazzyfu
@Desc    :   用户管理
*/

package views

import (
	"goInsight/internal/apps/users/forms"
	"goInsight/internal/apps/users/services"
	"goInsight/internal/pkg/response"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetUsersView(c *gin.Context) {
	var form *forms.GetUsersForm = &forms.GetUsersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetUsersServices{
			GetUsersForm: form,
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

func CreateUsersView(c *gin.Context) {
	var form *forms.CreateUsersForm = &forms.CreateUsersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CreateUsersService{
			CreateUsersForm: form,
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

func UpdateUsersView(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	var form *forms.UpdateUsersForm = &forms.UpdateUsersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.UpdateUsersService{
			UpdateUsersForm: form,
			C:               c,
			UID:             uint64(uid),
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

func DeleteUsersView(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	service := services.DeleteUsersService{
		C:   c,
		UID: uint64(uid),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}

func ChangeUsersPasswordView(c *gin.Context) {
	var form *forms.ChangeUsersPasswordForm = &forms.ChangeUsersPasswordForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ChangeUsersPasswordService{ChangeUsersPasswordForm: form, C: c}
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

func ChangeUserAvatarView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	service := services.ChangeUserAvatarService{
		C:        c,
		Username: username,
		File:     file,
	}
	err = service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}
