/*
@Time    :   2023/06/09 14:26:16
@Author  :   xff
*/

package views

import (
	"goInsight/internal/apps/das/forms"
	"goInsight/internal/apps/das/services"
	"goInsight/internal/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 获取环境
func GetEnvironmentsView(c *gin.Context) {
	service := services.GetEnvironmentsService{C: c}
	returnData, err := service.Run()
	if err == nil {
		response.Success(c, returnData, "success")
	} else {
		response.Fail(c, err.Error())
	}
}

// 获取用户授权的库
func GetSchemasView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	service := services.GetSchemasService{C: c, Username: username}
	returnData, err := service.Run()
	if err == nil {
		response.Success(c, returnData, "success")
	} else {
		response.Fail(c, err.Error())
	}
}

// 获取用户授权的库表
func GetTablesView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.GetTablesForm = &forms.GetTablesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetTablesService{GetTablesForm: form, C: c, Username: username}
		returnData, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, returnData, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 获取表元信息
func GetTableInfoView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.GetTableInfoForm = &forms.GetTableInfoForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetTableInfoService{GetTableInfoForm: form, C: c, Username: username}
		returnData, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, returnData, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 获取db字典
func GetDbDictView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.GetDbDictForm = &forms.GetDbDictForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetDbDictService{GetDbDictForm: form, C: c, Username: username}
		returnData, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, returnData, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 获取历史SQL
func GetHistoryView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.GetHistoryForm = &forms.GetHistoryForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetHistoryService{GetHistoryForm: form, C: c, Username: username}
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
