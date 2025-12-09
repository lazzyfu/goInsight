package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	ordersModels "github.com/lazzyfu/goinsight/internal/orders/models"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// 生成执行任务
func GenOrderTasksView(c *gin.Context) {
	username := jwt.ExtractClaims(c)["id"].(string)
	var form forms.GenOrderTasksForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.GenOrderTasksService{
		GenOrderTasksForm: &form,
		C:                 c,
		Username:          username,
	}
	if err := service.Run(); err != nil {
		response.Fail(c, err.Error())
		return
	}
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
			return
		}
		response.PaginationSuccess(c, total, returnData)
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
			return
		}
		response.Success(c, returnData, "success")
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 执行单个任务
func ExecuteSingleTaskView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.ExecuteSingleTaskForm = &forms.ExecuteSingleTaskForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.ExecuteSingleTaskService{
			ExecuteSingleTaskForm: form,
			C:                     c,
			Username:              username,
		}
		if err := service.Run(); err != nil {
			response.Fail(c, err.Error())
			return
		}
		response.Success(c, nil, "success")
	} else {
		response.ValidateFail(c, err.Error())
	}

}

// 批量执行任务
func ExecuteBatchTasksView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form forms.ExecuteBatchTasksForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.ExecuteBatchTasksService{
		ExecuteBatchTasksForm: &form,
		C:                     c,
		Username:              username,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "执行完成")
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
