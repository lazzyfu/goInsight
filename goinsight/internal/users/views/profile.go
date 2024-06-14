/*
@Time    :   2023/08/30 15:46:42
@Author  :   xff
@Desc    :   个人中心
*/
package views

import (
	"goInsight/internal/users/forms"
	"goInsight/internal/users/services"
	"goInsight/pkg/response"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetUserInfoView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
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
	uid, _ := strconv.Atoi(c.Param("uid"))
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
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
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
