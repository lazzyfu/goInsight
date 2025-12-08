package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/notifier"
	"github.com/lazzyfu/goinsight/pkg/pagination"
	"github.com/lazzyfu/goinsight/pkg/parser"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
	"github.com/lazzyfu/goinsight/internal/orders/api/execute"
	"github.com/lazzyfu/goinsight/internal/orders/forms"
	ordersModels "github.com/lazzyfu/goinsight/internal/orders/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenOrderTasksService struct {
	*forms.GenOrderTasksForm
	C        *gin.Context
	Username string
}

func (s *GenOrderTasksService) subTasksExist() bool {
	// 如果tasks记录存在，跳过
	var count int64
	global.App.DB.
		Model(&ordersModels.InsightOrderTasks{}).
		Where("order_id = ?", s.OrderID).
		Count(&count)
	return count == 0
}

func (s *GenOrderTasksService) Run() (err error) {
	// 工单是否存在
	var record ordersModels.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 检查是否有执行权限
	if s.Username != record.Claimer {
		return fmt.Errorf("您不是工单认领人，没有执行工单权限")
	}
	// 判断审核状态
	// 'PENDING','APPROVED','REJECTED','CLAIMED','EXECUTING','COMPLETED', 'FAILED','REVIEWED','REVOKED'
	if !utils.IsContain([]string{"CLAIMED", "EXECUTING", "COMPLETED", "FAILED", "REVIEWED"}, string(record.Progress)) {
		return fmt.Errorf("当前工单状态，禁止操作")
	}
	// 如果tasks记录存在，跳过
	if !s.subTasksExist() {
		return nil
	}
	// tasks记录不存在，生成记录
	sqls, err := parser.SplitSQLText(record.Content)
	if err != nil {
		return err
	}
	// 批量插入
	var tasks []map[string]any
	for _, sql := range sqls {
		tasks = append(tasks, map[string]any{
			"OrderID":    s.OrderID,
			"TaskID":     uuid.New(),
			"DBType":     record.DBType,
			"SQLType":    record.SQLType,
			"SQL":        sql,
			"created_at": time.Now().Format("2006-01-02 15:04:05"),
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ordersModels.InsightOrderTasks{}).Create(tasks).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}

// 获取任务列表
type GetTasksServices struct {
	*forms.GetTasksForm
	C        *gin.Context
	OrderID  string `json:"order_id"`
	Username string
}

func (s *GetTasksServices) Run() (responseData any, total int64, err error) {
	var records []ordersModels.InsightOrderTasks
	tx := global.App.DB.Table("`insight_order_tasks`").Where("order_id=?", s.OrderID).Scan(&records)
	if tx.RowsAffected == 0 {
		return records, total, fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 搜索
	if s.Search != "" {
		tx = tx.Where("sql like ?", "%"+s.Search+"%")
	}
	if s.Progress != "" {
		tx = tx.Where("progress=?", s.Progress)
	}
	total = pagination.Pager(&s.PaginationQ, tx, &records)
	return &records, total, nil
}

type PreviewTasksServices struct {
	*forms.PreviewTasksForm
	C *gin.Context
}

func (s *PreviewTasksServices) Run() (responseData any, err error) {
	type record struct {
		Total                        int `json:"total"`
		Unexecuted                   int `json:"unexecuted"`
		Processing                   int `json:"processing"`
		Completed                    int `json:"completed"`
		CompletedWithRollbackFailure int `json:"completed_with_rollback_failure"`
		Failed                       int `json:"failed"`
		Paused                       int `json:"paused"`
	}
	var records record
	global.App.DB.Table("`insight_order_tasks`").
		Select("COUNT(*) as total, SUM(if(progress='未执行',1,0)) as unexecuted, SUM(if(progress='执行中',1,0)) as processing, SUM(if(progress='已完成',1,0)) as completed, SUM(if(progress='已失败',1,0)) as failed,SUM(if(progress='已暂停',1,0)) as paused").
		Where("order_id=?", s.OrderID).
		Take(&records)

	return records, nil
}

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

// 检查当前工单的所有任务中是否有执行中的任务
func checkTasksProgressIsDoing(order_id string) bool {
	var records []ordersModels.InsightOrderTasks
	global.App.DB.Table("`insight_order_tasks`").Where("order_id=?", order_id).Scan(&records)
	for _, record := range records {
		if record.Progress == "执行中" {
			return false
		}
	}
	return true
}

// 检查当前工单的所有任务中是否有已暂停的任务
func checkTasksProgressIsPause(order_id string) bool {
	var records []ordersModels.InsightOrderTasks
	global.App.DB.Table("`insight_order_tasks`").Where("order_id=?", order_id).Scan(&records)
	for _, record := range records {
		if record.Progress == "已暂停" {
			return false
		}
	}
	return true
}

// 检查工单状态
func checkOrderStatus(order_id string, username string) error {
	var record ordersModels.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", order_id).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单记录`%s`不存在", order_id)
	}
	// 检查是否有执行权限
	var executorList []string
	err := json.Unmarshal([]byte(record.Executor), &executorList)
	if err != nil {
		return err
	}
	if !utils.IsContain(executorList, username) {
		return fmt.Errorf("您没有执行工单权限")
	}
	// 当工单的状态不为已批准或执行中的时候，禁止执行
	if utils.IsContain([]string{"已批准", "执行中"}, string(record.Progress)) {
		return nil
	}
	return fmt.Errorf("执行失败，当前工单状态为：%s", string(record.Progress))
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

// 执行单个任务
type ExecuteSingleTaskService struct {
	*forms.ExecuteSingleTaskForm
	C        *gin.Context
	Username string
}

func (s *ExecuteSingleTaskService) Run() (err error) {
	// 当工单的状态不为已批准或执行中的时候，禁止执行任务
	if err = checkOrderStatus(s.OrderID, s.Username); err != nil {
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
	if !checkTasksProgressIsDoing(s.OrderID) {
		return errors.New("当前有任务正在执行中，请先等待执行完成")
	}
	// 更新当前任务进度为执行中，工单状态为执行中
	if err := func() error {
		return global.App.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&ordersModels.InsightOrderTasks{}).
				Where("id=? and order_id=?", s.ID, s.OrderID).
				Update("progress", "执行中").Error; err != nil {
				global.App.Log.Error(err)
				return err
			}
			if err := tx.Model(&ordersModels.InsightOrderRecords{}).
				Where("order_id=?", s.OrderID).
				Update("progress", "执行中").Error; err != nil {
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
			taskProgress = "已失败"
		case base.RollbackSQLError:
			taskProgress = "已完成"
		default:
			taskProgress = "已失败"
		}
		global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
			Where("id=? and order_id=?", s.ID, s.OrderID).
			Updates(map[string]any{"progress": taskProgress, "result": data})
		return err
	}

	// 没有错误返回，更新任务状态为已完成
	global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
		Where("id=? and order_id=?", s.ID, s.OrderID).
		Updates(map[string]any{"progress": "已完成", "result": data})

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
	// 当工单的状态不为已批准或执行中的时候，禁止执行任务
	if err = checkOrderStatus(s.OrderID, s.Username); err != nil {
		return err
	}
	// 判断当前工单的所有任务中是否存在执行中的任务，如果存在，不执行
	if !checkTasksProgressIsDoing(s.OrderID) {
		return errors.New("当前有任务正在执行中，请先等待执行完成")
	}
	//  判断当前工单的所有任务中是否存在已暂停的任务，如果存在，不执行；可手动执行单个任务
	if !checkTasksProgressIsPause(s.OrderID) {
		return errors.New("当前有任务正在执行中，请先等待执行完成")
	}
	// 更新当前工单进度为执行中
	if err := global.App.DB.Model(&ordersModels.InsightOrderRecords{}).
		Where("order_id=?", s.OrderID).
		Update("progress", "执行中").Error; err != nil {
		global.App.Log.Error(err)
		return err
	}
	// 获取工单所有的任务
	var tasks []ordersModels.InsightOrderTasks
	tx := global.App.DB.Table("`insight_order_tasks`").Where("order_id=?", s.OrderID).Scan(&tasks)
	if tx.RowsAffected == 0 {
		return errors.New("任务记录不存在")
	}
	// 执行任务
	for _, task := range tasks {
		// 跳过已完成的任务
		if task.Progress == "已完成" {
			continue
		}

		// 更新当前任务进度为执行中
		if err := tx.Model(&ordersModels.InsightOrderTasks{}).
			Where("task_id=?", task.TaskID).
			Update("progress", "执行中").Error; err != nil {
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
				taskProgress = "已失败"
			case base.RollbackSQLError:
				taskProgress = "已完成"
			default:
				taskProgress = "已失败"
			}
			global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
				Where("task_id=?", task.TaskID).
				Updates(map[string]any{"progress": taskProgress, "result": data})
		} else {
			global.App.DB.Model(&ordersModels.InsightOrderTasks{}).
				Where("task_id=?", task.TaskID).
				Updates(map[string]any{"progress": "已完成", "result": data})

			// 导出工单需要发送导出文件信息给申请人、抄送人
			go sendExportFileInfoToApplicant(task.TaskID)
		}
	}
	// 更新工单状态为已完成
	updateOrderStatusToFinish(s.OrderID)
	return nil
}
