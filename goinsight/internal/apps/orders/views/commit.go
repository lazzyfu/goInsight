/*
@Time    :   2023/10/10 17:50:04
@Author  :   lazzyfu
*/

package views

import (
	"goInsight/internal/apps/orders/forms"
	"goInsight/internal/apps/orders/services"
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

// 获取指定环境的Schemas
func GetInstancesView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.GetInstancesForm = &forms.GetInstancesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetInstancesService{
			GetInstancesForm: form,
			C:                c,
			Username:         username,
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

// 获取指定环境的Schemas
func GetSchemasView(c *gin.Context) {
	var form *forms.GetSchemasForm = &forms.GetSchemasForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetSchemasService{
			GetSchemasForm: form,
			C:              c,
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

// 获取审核/复核/抄送人
func GetUsersView(c *gin.Context) {
	var form *forms.GetUsersForm = &forms.GetUsersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetUsersService{
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

// 提交工单
func CreateOrdersView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	// 解析表单数据
	var form *forms.CreateOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	// 执行创建
	service := services.CreateOrdersService{
		CreateOrderForm: form,
		C:               c,
		Username:        username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 成功响应
	response.Success(c, nil, "success")
}
