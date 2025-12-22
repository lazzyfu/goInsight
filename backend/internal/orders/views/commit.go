package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 获取环境
func GetOrderEnvironmentsView(c *gin.Context) {
	service := services.GetOrderEnvironmentsService{C: c}
	returnData, err := service.Run()
	if err == nil {
		response.Success(c, returnData, "success")
	} else {
		response.Fail(c, err.Error())
	}
}

// 获取指定环境的Schemas
func GetOrderInstancesView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form forms.GetOrderInstancesForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.GetOrderInstancesService{
		GetOrderInstancesForm: &form,
		C:                     c,
		Username:              username,
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)

}

// 获取指定环境的Schemas
func GetOrderSchemasView(c *gin.Context) {
	var form forms.GetOrderSchemasForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.GetOrderSchemasService{
		GetOrderSchemasForm: &form,
		C:                   c,
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)
}

// 获取审核/复核/抄送人
func GetOrderUsersView(c *gin.Context) {
	var form forms.GetOrderUsersForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.GetOrderUsersService{
		GetOrderUsersForm: &form,
		C:                 c,
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)
}

// 提交工单
func CreateOrderView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form forms.CreateOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.CreateOrderService{
		CreateOrderForm: &form,
		C:               c,
		Username:        username,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}
