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

// 尝试获取锁，如果返回 false 表示已经有用户在执行
func acquireOrderLock(ctx *gin.Context, orderID string) (bool, error) {
	key := fmt.Sprintf("order:lock:%s", orderID)
	ok, err := global.App.Redis.SetNX(ctx, key, "1", 0).Result()
	return ok, err
}

// 释放锁
func releaseOrderLock(ctx *gin.Context, orderID string) error {
	key := fmt.Sprintf("order:lock:%s", orderID)
	_, err := global.App.Redis.Del(ctx, key).Result()
	return err
}

// 所有任务均为 COMPLETED -> 将工单置为 COMPLETED 并通知申请人
func updateOrderStatusToFinish(orderID string) error {
	err := global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 统计状态
		var counts struct {
			NotDone int64
		}

		if err := tx.Table("`insight_order_tasks`").
			Select(
				"SUM(CASE WHEN progress != ? THEN 1 ELSE 0 END) AS not_done",
				"COMPLETED").
			Where("order_id = ?", orderID).
			Scan(&counts).Error; err != nil {
			return err
		}

		// 如果没有未完成任务（NotDone==0），则全部完成
		if counts.NotDone == 0 {
			if err := tx.Model(&ordersModels.InsightOrderRecords{}).
				Where("order_id = ?", orderID).
				Update("progress", "COMPLETED").Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 事务提交成功后，如果全部完成则发送消息
	var record ordersModels.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", orderID).Take(&record).Error; err == nil {
		if record.Progress == "COMPLETED" {
			// 发送消息，发送给申请人
			receiver := []string{record.Applicant}
			notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderExecutionCompleted, notifier.MessageParams{
				Order: &record,
			})
		}
	}

	return nil
}

// 判断当前工单是否没有执行中的任务
func noTaskExecuting(orderID string) bool {
	var count int64
	global.App.DB.Table("`insight_order_tasks`").
		Where("order_id = ? AND progress = 'EXECUTING'", orderID).
		Count(&count)
	return count == 0
}

func CheckOrderExecutable(record *ordersModels.InsightOrderRecords) error {
	// 检查工单状态是否允许执行
	// 'PENDING','APPROVED','REJECTED','CLAIMED','EXECUTING','COMPLETED', 'FAILED','REVIEWED','REVOKED'
	if !utils.IsContain([]string{"CLAIMED", "EXECUTING"}, string(record.Progress)) {
		progressMap := map[string]string{
			"PENDING":   "待审批",
			"APPROVED":  "已批准",
			"REJECTED":  "已驳回",
			"CLAIMED":   "已认领",
			"EXECUTING": "执行中",
			"COMPLETED": "已完成",
			"FAILED":    "已失败",
			"REVIEWED":  "已复核",
			"REVOKED":   "已撤销",
		}
		progressCN := progressMap[string(record.Progress)]
		if progressCN == "" {
			progressCN = string(record.Progress)
		}
		return fmt.Errorf("当前工单%s，禁止执行", progressCN)
	}
	return nil
}

func sendExportFileInfoToApplicant(orderID uuid.UUID) {
	var task ordersModels.InsightOrderTasks
	global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
		Where("task_id=?", orderID).Scan(&task)

	var record ordersModels.InsightOrderRecords
	global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
		Where("order_id=?", task.OrderID).Scan(&record)

	if record.SQLType != "EXPORT" {
		return
	}

	var file base.ExportFile
	_ = json.Unmarshal([]byte(task.Result), &file)

	receiver := []string{record.Applicant}
	notifier.SendOrderMessage(receiver, notifier.MsgTypeExportFileInfo, notifier.MessageParams{
		Order: &record,
		Task:  &task,
	})
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
		Joins("join `insight_instances` b on a.instance_id=b.instance_id").
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
	// 转换为json
	data, _ := json.Marshal(returnData)
	return string(data), err
}

// 执行单个任务
type ExecuteTaskService struct {
	*forms.ExecuteTaskForm
	C        *gin.Context
	Username string
}

func (s *ExecuteTaskService) Run() (err error) {
	// 检查工单记录是否存在
	var record ordersModels.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.Error != nil {
		return fmt.Errorf("查询工单记录失败: %v", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单记录`%s`不存在", s.OrderID)
	}
	// 检查是否有工单执行权限
	if record.Claimer != s.Username {
		return fmt.Errorf("您不是工单认领人，没有执行工单权限")
	}
	// 检查工单状态是否允许执行
	if err := CheckOrderExecutable(&record); err != nil {
		return err
	}
	// 获取锁
	locked, err := acquireOrderLock(s.C, s.OrderID)
	if err != nil {
		return fmt.Errorf("获取工单锁失败: %v", err)
	}
	if !locked {
		return fmt.Errorf("工单正在执行中，请稍后再试")
	}
	defer releaseOrderLock(s.C, s.OrderID)
	// 获取任务记录
	var task ordersModels.InsightOrderTasks
	tx = global.App.DB.Table("`insight_order_tasks`").Where("task_id=? and order_id=?", s.TaskID, s.OrderID).Take(&task)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("任务记录不存在")
	}
	// 跳过已完成的任务
	if task.Progress == "COMPLETED" {
		return errors.New("当前任务已完成，请勿重复执行")
	}
	// 跳过执行中的任务
	if task.Progress == "EXECUTING" {
		return errors.New("当前任务正在执行中，请勿重复执行")
	}
	// 更新当前任务进度为执行中，工单状态为执行中
	if err := func() error {
		return global.App.DB.Transaction(func(tx *gorm.DB) error {
			if tx = tx.Model(&ordersModels.InsightOrderTasks{}).
				Where("task_id=? and order_id=?", s.TaskID, s.OrderID).
				Updates(map[string]any{"progress": "EXECUTING", "executor": s.Username}); tx.Error != nil {
				global.App.Log.Error(tx.Error)
				return tx.Error
			}
			if tx = tx.Model(&ordersModels.InsightOrderRecords{}).
				Where("order_id=?", s.OrderID).
				Updates(map[string]any{"progress": "EXECUTING", "executor": s.Username}); tx.Error != nil {
				global.App.Log.Error(tx.Error)
				return tx.Error
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
			Where("task_id=? and order_id=?", s.TaskID, s.OrderID).
			Updates(map[string]any{"progress": taskProgress, "result": data})
		return err
	}

	// 没有错误返回，更新任务状态为已完成
	global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
		Where("task_id=? and order_id=?", s.TaskID, s.OrderID).
		Updates(map[string]any{"progress": "COMPLETED", "result": data})

	// 导出工单需要发送导出文件信息给申请人、抄送人
	go sendExportFileInfoToApplicant(task.TaskID)

	// 更新工单状态为已完成
	updateOrderStatusToFinish(s.OrderID)
	return nil
}

// 批量执行任务
type ExecuteBatchTasksService struct {
	*forms.ExecuteBatchTasksForm
	C        *gin.Context
	Username string
}

func (s *ExecuteBatchTasksService) Run() (err error) {
	// 检查工单状态和执行权限
	var record ordersModels.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.Error != nil {
		return fmt.Errorf("查询工单记录失败: %v", tx.Error)
	}
	// 检查工单记录是否存在
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单记录`%s`不存在", s.OrderID)
	}
	// 检查是否有工单执行权限
	if record.Claimer != s.Username {
		return fmt.Errorf("您不是工单认领人，没有执行工单权限")
	}
	// 检查工单状态是否允许执行
	if err := CheckOrderExecutable(&record); err != nil {
		return err
	}
	// 判断当前工单是否存在执行中的任务，避免跳过执行中的任务执行其他的任务
	if !noTaskExecuting(s.OrderID) {
		return errors.New("当前有任务正在执行中，请先等待执行完成")
	}
	// 获取锁
	locked, err := acquireOrderLock(s.C, s.OrderID)
	if err != nil {
		return fmt.Errorf("获取工单锁失败: %v", err)
	}
	if !locked {
		return fmt.Errorf("工单正在执行中，请稍后再试")
	}
	defer releaseOrderLock(s.C, s.OrderID)
	// 更新当前工单进度为执行中
	if tx = global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
		Where("order_id=?", s.OrderID).
		Updates(map[string]any{"progress": "EXECUTING", "executor": s.Username}); tx.Error != nil {
		global.App.Log.Error(tx.Error)
		return tx.Error
	}
	// 获取当前工单所有的任务
	var tasks []ordersModels.InsightOrderTasks
	if tx = global.App.DB.Where("order_id=?", s.OrderID).Find(&tasks); tx.Error != nil {
		return tx.Error
	}
	// 串行执行任务
	for _, task := range tasks {
		// 跳过已完成的任务，避免重复执行
		if task.Progress == "COMPLETED" {
			continue
		}

		// 更新当前任务进度为执行中
		if tx = global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
			Where("task_id=?", task.TaskID).
			Update("progress", "EXECUTING"); tx.Error != nil {
			global.App.Log.Error(tx.Error)
			return tx.Error
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
