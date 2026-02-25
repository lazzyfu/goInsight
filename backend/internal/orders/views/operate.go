package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	"github.com/gin-gonic/gin"
)

// 审批
func ApprovalOrderView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form forms.ApprovalOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.ApprovalOrderService{
		ApprovalOrderForm: &form,
		C:                 c,
		Username:          username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

// 认领
func ClaimOrderView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form forms.ClaimOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.ClaimOrderService{
		ClaimOrderForm: &form,
		C:              c,
		Username:       username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

// 转交工单给其他执行人
func TransferOrderView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form forms.TransferOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.TransferOrderService{
		TransferOrderForm: &form,
		C:                 c,
		Username:          username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

// 撤销工单
func RevokeOrderView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form forms.RevokeOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.RevokeOrderService{
		RevokeOrderForm: &form,
		C:               c,
		Username:        username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

// 手动更新工单为已完成
func CompleteOrderView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form forms.CompleteOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.CompleteOrderService{
		CompleteOrderForm: &form,
		C:                 c,
		Username:          username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

// 手动更新工单为失败
func FailOrderView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form forms.FailOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.FailOrderService{
		FailOrderForm: &form,
		C:             c,
		Username:      username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

// 复核
func ReviewOrderView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form forms.ReviewOrderForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.ReviewOrderService{
		ReviewOrderForm: &form,
		C:               c,
		Username:        username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}
