/*
@Time    :   2023/10/10 17:50:04
@Author  :   xff
*/

package views

import (
	"goInsight/internal/orders/forms"
	"goInsight/internal/orders/services"
	"goInsight/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetListView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.GetListForm = &forms.GetListForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetListServices{
			GetListForm: form,
			C:           c,
			Username:    username,
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

func GetDetailView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	service := services.GetDetailServices{
		OrderID:  c.Param("order_id"),
		C:        c,
		Username: username,
	}
	returnData, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, returnData, "success")
	}
}

func GetOpLogsView(c *gin.Context) {
	var form *forms.GetOpLogsForm = &forms.GetOpLogsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetOpLogsServices{
			GetOpLogsForm: form,
			C:             c,
		}
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
