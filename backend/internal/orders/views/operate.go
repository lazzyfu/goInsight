package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 审批
func ApprovalView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.ApprovalForm = &forms.ApprovalForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ApprovalService{
			ApprovalForm: form,
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
func CloseView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form *forms.CloseForm = &forms.CloseForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CloseService{
			CloseForm: form,
			C:         c,
			Username:  username,
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
