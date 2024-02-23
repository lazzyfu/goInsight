package views

import (
	"encoding/json"
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/orders/forms"
	ordersModels "goInsight/internal/apps/orders/models"
	"goInsight/internal/apps/orders/services"
	"goInsight/internal/pkg/response"
	"io"
	"os"

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

// 下载导出文件
func DownloadExportFileView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	task_id := c.Param("task_id")
	// 判断下载记录是否存在
	var task ordersModels.InsightOrderTasks
	tx := global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
		Where("task_id=?", task_id).Scan(&task)
	if tx.RowsAffected == 0 {
		response.Fail(c, "记录不存在")
		return
	}
	// 判断用户是否有权限下载
	var record ordersModels.InsightOrderRecords
	global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
		Where("order_id=?", task.OrderID).Scan(&record)
	if record.Applicant != username {
		response.Fail(c, "无权限下载")
		return
	}
	// 获取下载文件信息
	type exportFile struct {
		FileName      string `json:"file_name"`
		FileSize      int64  `json:"file_size"`
		FilePath      string `json:"file_path"`
		ContentType   string `json:"content_type"`
		EncryptionKey string `json:"encryption_key"`
		ExportRows    int64  `json:"export_rows"`
	}
	var file exportFile
	err := json.Unmarshal([]byte(task.Result), &file)
	if err != nil {
		response.Fail(c, fmt.Sprintf("解析下载文件信息异常，错误：%s", err.Error()))
		return
	}
	// 检查本地文件是否存在
	_, err = os.Stat(file.FilePath)
	if os.IsNotExist(err) {
		response.Fail(c, "下载文件不存在")
		return
	}

	// Open the file
	f, err := os.Open(file.FilePath)
	if err != nil {
		response.Fail(c, fmt.Sprintf("下载文件读取失败，错误：%s", err.Error()))
		return
	}
	defer f.Close()

	// Set headers for file download
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))

	// Copy the file to the response writer
	_, err = io.Copy(c.Writer, f)
	if err != nil {
		response.Fail(c, fmt.Sprintf("下载文件发送失败，错误：%s", err.Error()))
		return
	}
}
