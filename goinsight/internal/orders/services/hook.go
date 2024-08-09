/*
@Time    :   2023/10/24 19:14:17
@Author  :   xff
*/

package services

import (
	"encoding/json"
	"fmt"
	"goInsight/global"
	"goInsight/internal/orders/forms"
	"goInsight/internal/orders/models"
	"goInsight/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	commonModels "goInsight/internal/common/models"
)

// hook工单
type HookOrdersService struct {
	*forms.HookOrdersForm
	C        *gin.Context
	Username string
}

// 重置工单审批人/复核人
func (s *HookOrdersService) resetStatus(users datatypes.JSON) (datatypes.JSON, error) {
	var tmpData []map[string]interface{}
	err := json.Unmarshal(users, &tmpData)
	if err != nil {
		return nil, err
	}
	for _, u := range tmpData {
		u["msg"] = ""
		u["status"] = "pending"
		u["time"] = ""
	}
	data, err := json.Marshal(tmpData)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(data), nil
}

func (s *HookOrdersService) Run() error {
	// 判断工单是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}

	// 判断db类型是否一致
	if s.DBType != record.DBType {
		return fmt.Errorf("记录`%s`的db类型(%s)与当前db类型(%s)不一致", s.OrderID, record.DBType, s.DBType)
	}
	// 判断进度
	approver := record.Approver

	// 重置审核状态
	if s.Progress == "待审核" {
		var err error
		approver, err = s.resetStatus(record.Approver)
		if err != nil {
			return err
		}
	}
	// 重置复核状态
	reviewer, err := s.resetStatus(record.Reviewer)
	if err != nil {
		return err
	}

	var hookRecords []models.InsightOrderRecords
	for _, item := range s.Target {
		// 解析UUID
		instance_id, err := utils.ParserUUID(item.InstanceID)
		if err != nil {
			return err
		}
		// 生成新的工单ID
		orderID := uuid.New()
		hookTitle := record.Title
		if !strings.HasPrefix(record.Title, "[Hook]") {
			hookTitle = fmt.Sprintf("[Hook]%s", record.Title)
		}
		hookRecords = append(hookRecords, models.InsightOrderRecords{
			Title:            hookTitle,
			Progress:         commonModels.EnumType(s.Progress),
			OrderID:          orderID,
			HookOrderID:      record.OrderID,
			Remark:           record.Remark,
			IsRestrictAccess: record.IsRestrictAccess,
			DBType:           s.DBType,
			SQLType:          record.SQLType,
			Environment:      item.Environment,
			InstanceID:       instance_id,
			Schema:           item.Schema,
			Applicant:        s.Username, // warn：谁执行的hook，申请人改为谁
			Organization:     record.Organization,
			Approver:         approver,
			Executor:         record.Executor,
			Reviewer:         reviewer,
			CC:               record.CC,
			Content:          record.Content,
		})
	}
	// 批量插入
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).CreateInBatches(&hookRecords, 100).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录已存在，错误:%s", err.Error())
			}
			global.App.Log.Error(err)
			return err
		}
		// 操作日志
		logMsg := fmt.Sprintf("用户%s触发了HOOK操作", s.Username)
		if err := CreateOpLogs(tx, record.OrderID, s.Username, logMsg); err != nil {
			return err
		}
		return nil
	})
}
