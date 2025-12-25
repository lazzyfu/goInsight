package services

import (
	"fmt"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/parser"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	ordersModels "github.com/lazzyfu/goinsight/internal/orders/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 生成执行任务
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
