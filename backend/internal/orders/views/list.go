package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetOrderListView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form forms.GetOrderListForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.GetOrderListServices{
		GetOrderListForm: &form,
		C:                c,
		Username:         username,
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)
}

func GetOrderDetailView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	service := services.GetOrderDetailServices{
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

func GetOrderApprovalView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	service := services.GetOrderApprovalServices{
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

func GetOrderLogsView(c *gin.Context) {
	service := services.GetOrderLogsServices{
		OrderID: c.Param("order_id"),
		C:       c,
	}
	returnData, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, returnData, "success")
	}
}
