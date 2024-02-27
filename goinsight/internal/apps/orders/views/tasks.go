package views

import (
	"encoding/json"
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/orders/forms"
	ordersModels "goInsight/internal/apps/orders/models"
	"goInsight/internal/apps/orders/services"
	"goInsight/internal/pkg/response"
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/requestid"
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
	requestID := requestid.Get(c)
	// 判断下载记录是否存在
	var task ordersModels.InsightOrderTasks
	tx := global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
		Where("task_id=?", task_id).Scan(&task)
	if tx.RowsAffected == 0 {
		global.App.Log.WithField("request_id", requestID).WithField("username", username).Errorf("下载记录关联的任务不存在，任务ID：%s", task_id)
		c.JSON(http.StatusNotFound, map[string]interface{}{})
		return
	}
	// 判断用户是否有权限下载
	var record ordersModels.InsightOrderRecords
	global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
		Where("order_id=?", task.OrderID).Scan(&record)
	if record.Applicant != username {
		global.App.Log.WithField("request_id", requestID).WithField("username", username).Errorf("用户%s无权限下载任务%s的导出文件", username, task_id)
		c.JSON(http.StatusForbidden, map[string]interface{}{})
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
		global.App.Log.WithField("request_id", requestID).WithField("username", username).Errorf("解析下载文件信息异常，错误：%s", err.Error())
		c.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	// 检查本地文件是否存在
	_, err = os.Stat(file.FilePath)
	if os.IsNotExist(err) {
		global.App.Log.WithField("request_id", requestID).WithField("username", username).Error(fmt.Sprintf("下载的文件%s不存在", file.FilePath))
		c.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}

	// Open the file
	f, err := os.Open(file.FilePath)
	if err != nil {
		global.App.Log.WithField("request_id", requestID).WithField("username", username).Error(fmt.Sprintf("打开下载文件%s失败", file.FilePath))
		c.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	defer f.Close()

	// Set headers for file download
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Header("Accept-Length", fmt.Sprintf("%d", file.FileSize))
	c.File(file.FilePath)
}
