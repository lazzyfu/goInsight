package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/users/forms"
	"github.com/lazzyfu/goinsight/internal/users/services"

	"github.com/gin-gonic/gin"
)

func GetUsersView(c *gin.Context) {
	var form forms.GetUsersForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.GetUsersServices{
		GetUsersForm: &form,
		C:            c,
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)
}

func CreateUsersView(c *gin.Context) {
	var form forms.CreateUsersForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.CreateUsersService{
		CreateUsersForm: &form,
		C:               c,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func UpdateUsersView(c *gin.Context) {
	uid, ok := parseUint64Param(c, "uid")
	if !ok {
		return
	}
	var form forms.UpdateUsersForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.UpdateUsersService{
		UpdateUsersForm: &form,
		C:               c,
		UID:             uid,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func DeleteUsersView(c *gin.Context) {
	uid, ok := parseUint64Param(c, "uid")
	if !ok {
		return
	}
	service := services.DeleteUsersService{
		C:   c,
		UID: uid,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func ResetUsersPasswordView(c *gin.Context) {
	var form forms.ResetUsersPasswordForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.ResetUsersPasswordService{
		ResetUsersPasswordForm: &form,
		C:                      c,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func ChangeUserAvatarView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
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
		return
	}
	response.Success(c, nil, "success")
}

func GetUserOrganizationsView(c *gin.Context) {
	uid, ok := parseUint64Param(c, "uid")
	if !ok {
		return
	}
	service := services.GetUserOrganizationsService{
		UID: uid,
	}
	returnData, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, returnData, "success")
}
