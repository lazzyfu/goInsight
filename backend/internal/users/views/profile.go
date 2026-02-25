package views

import (
	"net/http"

	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/users/forms"
	"github.com/lazzyfu/goinsight/internal/users/services"

	"github.com/gin-gonic/gin"
)

func GetUserInfoView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	service := services.GetUserInfoServices{C: c, Username: username}
	returnData, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, returnData, "success")
	}
}

func UpdateUserInfoView(c *gin.Context) {
	var form *forms.UpdateUserInfoForm = &forms.UpdateUserInfoForm{}
	uid, ok := parseUint64Param(c, "uid")
	if !ok {
		return
	}
	if err := c.ShouldBind(&form); err == nil {
		service := services.UpdateUserInfoService{
			UpdateUserInfoForm: form,
			C:                  c,
			UID:                uint32(uid),
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

// 用户修改密码
func ChangeUserPasswordView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.ChangeUserPasswordForm = &forms.ChangeUserPasswordForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ChangeUserPasswordService{ChangeUserPasswordForm: form, C: c, Username: username}
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

func GetOTPAuthURLView(c *gin.Context) {
	var form *forms.GetOTPAuthURLForm = &forms.GetOTPAuthURLForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetOTPAuthURLService{GetOTPAuthURLForm: form, C: c}
		data, err := service.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, data)
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

func GetOTPAuthCallbackView(c *gin.Context) {
	var form *forms.GetOTPAuthCallbackForm = &forms.GetOTPAuthCallbackForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetOTPAuthCallbackService{GetOTPAuthCallbackForm: form, C: c}
		err := service.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "0001", "message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": "0000", "message": "OTP绑定成功"})
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}
