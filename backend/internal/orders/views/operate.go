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

// 更新状态未执行中或已完成
func FeedbackView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.FeedbackForm = &forms.FeedbackForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.FeedbackService{
			FeedbackForm: form,
			C:            c,
			Username:     username,
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

// 更新状态未执行中或已完成
func ReviewView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.ReviewForm = &forms.ReviewForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ReviewService{
			ReviewForm: form,
			C:          c,
			Username:   username,
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

// 关闭工单
func CloseOrderView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.CloseOrderForm = &forms.CloseOrderForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CloseOrderService{
			CloseOrderForm: form,
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
