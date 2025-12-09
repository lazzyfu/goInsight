package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/notifier"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
	"github.com/lazzyfu/goinsight/internal/orders/api/execute"
	"github.com/lazzyfu/goinsight/internal/orders/forms"
	ordersModels "github.com/lazzyfu/goinsight/internal/orders/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 检查工单所有任务是否完成，如果所有子任务已完成，更新工单状态为已完成
func updateOrderStatusToFinish(order_id string) {
	// 判断所有任务是否都完成
	type TaskCount struct {
		Count int64
	}
	var taskCount TaskCount
	global.App.DB.Table("`insight_order_tasks`").
		Select("count(*) as count").
		Where("order_id=? and progress not in ('已完成')", order_id).
		Scan(&taskCount)
	if taskCount.Count == 0 {
		// 更新工单为已完成
		global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
			Where("order_id=?", order_id).
			Update("progress", "已完成")

		// 发送通知消息
		var record ordersModels.InsightOrderRecords
		global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
			Where("order_id=?", order_id).Scan(&record)
		receiver := []string{record.Applicant}
		msg := fmt.Sprintf(
			"您好，工单已经执行完成，请悉知\n"+
				">工单标题：%s",
			record.Title,
		)
		notifier.SendMessage(record.Title, order_id, receiver, msg)
	}
}

// 判断当前工单是否没有执行中的任务
func noTaskExecuting(orderID string) bool {
	var count int64
	global.App.DB.Table("`insight_order_tasks`").
		Where("order_id=? AND progress='EXECUTING'", orderID).
		Count(&count)
	return count == 0
}

// 检查工单状态
func checkOrderStatusAndPerm(order_id string, username string) error {
	var record ordersModels.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", order_id).Take(&record)
	if tx.Error != nil {
		return fmt.Errorf("查询工单记录失败: %v", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单记录`%s`不存在", order_id)
	}
	// 检查是否有执行权限
	if record.Claimer != username {
		return fmt.Errorf("您不是工单认领人，没有执行工单权限")
	}
	// 检查状态是否允许执行
	if !utils.IsContain([]string{"CLAIMED", "EXECUTING"}, string(record.Progress)) {
		return fmt.Errorf("当前工单状态不为已认领或执行中，禁止执行")
	}
	return nil
}

func sendExportFileInfoToApplicant(task_id uuid.UUID) {
	var task ordersModels.InsightOrderTasks
	global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
		Where("task_id=?", task_id).Scan(&task)

	var record ordersModels.InsightOrderRecords
	global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
		Where("order_id=?", task.OrderID).Scan(&record)

	if record.SQLType != "EXPORT" {
		return
	}

	var file base.ExportFile
	_ = json.Unmarshal([]byte(task.Result), &file)

	receiver := []string{record.Applicant}
	msg := fmt.Sprintf(
		"您好，导出文件信息如下，请查收\n"+
			">工单标题：%s\n"+
			">任务ID：%s\n"+
			">文件名：%s\n"+
			">文件大小：%d字节\n"+
			">数据行数：%d\n"+
			">文件解密密码：%s\n"+
			">文件格式：%s\n"+
			">文件下载路径：%s",
		record.Title, task_id.String(),
		file.FileName,
		file.FileSize,
		file.ExportRows,
		file.EncryptionKey,
		file.ContentType,
		file.DownloadUrl,
	)
	notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
}

// 执行任务
func executeTask(task ordersModels.InsightOrderTasks) (string, error) {
	// 获取DB配置信息
	type Record struct {
		Hostname         string
		Port             uint16
		User             string
		Password         string
		Schema           string
		DBType           string
		SQLType          string
		ExportFileFormat string
	}
	var record Record
	tx := global.App.DB.Table("`insight_order_records` a").
		Select("a.db_type,a.sql_type,a.schema,a.export_file_format,b.hostname,b.port,b.user,b.password").
		Joins("join `insight_db_config` b on a.instance_id=b.instance_id").
		Where("a.order_id=?", task.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		returnData := base.ReturnData{Error: "执行失败，没有发现工单关联的数据库信息"}
		data, _ := json.Marshal(returnData)
		return string(data), errors.New("执行失败，没有发现工单关联的数据库信息")
	}
	// 解密密码
	plainPassword, err := utils.Decrypt(record.Password)
	if err != nil {
		return "", err
	}
	config := base.DBConfig{
		Hostname:         record.Hostname,
		Port:             record.Port,
		UserName:         record.User,
		Password:         plainPassword,
		Schema:           record.Schema,
		DBType:           record.DBType,
		SQLType:          record.SQLType,
		ExportFileFormat: record.ExportFileFormat,
		SQL:              task.SQL,
		OrderID:          task.OrderID.String(),
		TaskID:           task.TaskID.String(),
	}
	// 执行工单
	executor := execute.NewExecuteSQLAPI(&config)
	returnData, err := executor.Run()
	if err != nil {
		base.PublishMessageToChannel(task.OrderID.String(), err.Error(), "")
	}
	// 转换为json
	data, _ := json.Marshal(returnData)
	return string(data), err
}

// ---------- 执行单个任务 ----------
type ExecuteSingleTaskService struct {
	*forms.ExecuteSingleTaskForm
	C        *gin.Context
	Username string
}

func (s *ExecuteSingleTaskService) Run() (err error) {
	// 当工单的状态不为已批准或执行中的时候，禁止执行任务
	if err = checkOrderStatusAndPerm(s.OrderID, s.Username); err != nil {
		return err
	}
	// 获取任务记录
	var task ordersModels.InsightOrderTasks
	tx := global.App.DB.Table("`insight_order_tasks`").Where("id=? and order_id=?", s.ID, s.OrderID).Take(&task)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("任务ID为`%d`的记录不存在", s.ID)
	}
	// 跳过已完成的任务
	if task.Progress == "已完成" {
		return errors.New("当前任务已完成，请勿重复执行")
	}
	// 跳过执行中的任务
	if task.Progress == "执行中" {
		return errors.New("当前任务正在执行中，请勿重复执行")
	}
	// 判断当前工单的所有任务是否存在执行中的任务，避免跳过执行中的任务执行其他的任务
	if !noTaskExecuting(s.OrderID) {
		return errors.New("当前有任务正在执行中，请先等待执行完成")
	}
	// 更新当前任务进度为执行中，工单状态为执行中
	if err := func() error {
		return global.App.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&ordersModels.InsightOrderTasks{}).
				Where("id=? and order_id=?", s.ID, s.OrderID).
				Update("progress", "EXECUTING").Error; err != nil {
				global.App.Log.Error(err)
				return err
			}
			if err := tx.Model(&ordersModels.InsightOrderRecords{}).
				Where("order_id=?", s.OrderID).
				Update("progress", "EXECUTING").Error; err != nil {
				global.App.Log.Error(err)
				return err
			}
			return nil
		})
	}(); err != nil {
		return err
	}

	// 执行任务
	data, err := executeTask(task)

	// 返回错误，更新任务状态
	if err != nil {
		var taskProgress string
		// 错误类型断言，可以添加更多状态
		switch err.(type) {
		case base.SQLExecuteError:
			taskProgress = "FAILED"
		case base.RollbackSQLError:
			taskProgress = "COMPLETED"
		default:
			taskProgress = "FAILED"
		}
		global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
			Where("id=? and order_id=?", s.ID, s.OrderID).
			Updates(map[string]any{"progress": taskProgress, "result": data})
		return err
	}

	// 没有错误返回，更新任务状态为已完成
	global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
		Where("id=? and order_id=?", s.ID, s.OrderID).
		Updates(map[string]any{"progress": "COMPLETED", "result": data})

	// 导出工单需要发送导出文件信息给申请人、抄送人
	go sendExportFileInfoToApplicant(task.TaskID)

	// 更新工单状态为已完成
	updateOrderStatusToFinish(s.OrderID)
	return nil
}

// ---------- 批量执行任务 ----------
type ExecuteBatchTasksService struct {
	*forms.ExecuteBatchTasksForm
	C        *gin.Context
	Username string
}

func (s *ExecuteBatchTasksService) Run() (err error) {
	// 检查工单状态和执行权限
	if err = checkOrderStatusAndPerm(s.OrderID, s.Username); err != nil {
		return err
	}
	// 判断当前工单的所有任务是否存在执行中的任务，避免跳过执行中的任务执行其他的任务
	if !noTaskExecuting(s.OrderID) {
		return errors.New("当前有任务正在执行中，请先等待执行完成")
	}
	// 更新当前工单进度为执行中
	if err := global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
		Where("order_id=?", s.OrderID).
		Update("progress", "EXECUTING").Error; err != nil {
		global.App.Log.Error(err)
		return err
	}
	// 获取工单所有的任务
	var tasks []ordersModels.InsightOrderTasks
	if err := global.App.DB.Where("order_id=?", s.OrderID).Find(&tasks).Error; err != nil {
		return err
	}
	// 串行执行任务
	for _, task := range tasks {
		// 跳过已完成的任务
		if task.Progress == "COMPLETED" {
			continue
		}

		// 更新当前任务进度为执行中
		if err := global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
			Where("task_id=?", task.TaskID).
			Update("progress", "EXECUTING").Error; err != nil {
			global.App.Log.Error(err)
			return err
		}

		// 执行任务
		data, err := executeTask(task)

		// 返回错误，更新任务状态
		if err != nil {
			var taskProgress string
			// 错误类型断言，可以添加更多状态
			switch err.(type) {
			case base.SQLExecuteError:
				taskProgress = "FAILED"
			case base.RollbackSQLError:
				taskProgress = "COMPLETED"
			default:
				taskProgress = "FAILED"
			}
			global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
				Where("task_id=?", task.TaskID).
				Updates(map[string]any{"progress": taskProgress, "result": data})
		} else {
			global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
				Where("task_id=?", task.TaskID).
				Updates(map[string]any{"progress": "COMPLETED", "result": data})

			// 导出工单需要发送导出文件信息给申请人、抄送人
			go sendExportFileInfoToApplicant(task.TaskID)
		}
	}
	// 更新工单状态为已完成
	updateOrderStatusToFinish(s.OrderID)
	return nil
}
