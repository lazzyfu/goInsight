package views

import (
	"goInsight/internal/app/orders/forms"
	"goInsight/internal/app/orders/services"
	"goInsight/internal/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 提交工单
func GenerateTasksView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	// 解析表单数据
	var form forms.GenerateTasksForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	// 执行
	service := services.GenerateTasksService{
		GenerateTasksForm: form,
		C:                 c,
		Username:          username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 成功响应
	response.Success(c, nil, "success")
}

// 获取任务列表
func GetTasksView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.GetTasksForm = &forms.GetTasksForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetTasksServices{
			GetTasksForm: form,
			OrderID:      c.Param("order_id"),
			C:            c,
			Username:     username,
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

// 预览任务
func PreviewTasksView(c *gin.Context) {
	var form *forms.PreviewTasksForm = &forms.PreviewTasksForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.PreviewTasksServices{
			PreviewTasksForm: form,
			C:                c,
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

// 执行单个任务
func ExecuteSingleTaskView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	// 解析表单数据
	var form forms.ExecuteSingleTaskForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	// 执行
	service := services.ExecuteSingleTaskService{
		ExecuteSingleTaskForm: form,
		C:                     c,
		Username:              username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 成功响应
	response.Success(c, nil, "success")
}

// 批量执行任务
func ExecuteAllTaskView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.ExecuteAllTaskForm = &forms.ExecuteAllTaskForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ExecuteAllTaskService{
			ExecuteAllTaskForm: form,
			C:                  c,
			Username:           username,
		}
		err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, nil, "执行完成")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}
