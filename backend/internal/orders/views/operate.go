package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 审批
func ApprovalOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.ApprovalOrderForm = &forms.ApprovalOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ApprovalOrderService{
			ApprovalOrderForm: form,
			C:                 c,
			Username:          username,
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

// 认领
func ClaimOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.ClaimOrderForm = &forms.ClaimOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ClaimOrderService{
			ClaimOrderForm: form,
			C:              c,
			Username:       username,
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

// 转交工单给其他执行人
func TransferOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.TransferOrderForm = &forms.TransferOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.TransferOrderService{
			TransferOrderForm: form,
			C:                 c,
			Username:          username,
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

// 撤销工单
func RevokeOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.RevokeOrderForm = &forms.RevokeOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.RevokeOrderService{
			RevokeOrderForm: form,
			C:               c,
			Username:        username,
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

// 手动更新工单为已完成
func CompleteOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.CompleteOrderForm = &forms.CompleteOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CompleteOrderService{
			CompleteOrderForm: form,
			C:                 c,
			Username:          username,
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

// 手动更新工单为失败
func FailOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.FailOrderForm = &forms.FailOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.FailOrderService{
			FailOrderForm: form,
			C:             c,
			Username:      username,
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

// 复核
func ReviewOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.ReviewOrderForm = &forms.ReviewOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ReviewOrderService{
			ReviewOrderForm: form,
			C:               c,
			Username:        username,
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
