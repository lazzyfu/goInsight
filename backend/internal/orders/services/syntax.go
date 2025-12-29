package services

import (
	"fmt"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/parser"
	"github.com/lazzyfu/goinsight/pkg/utils"

	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/inspect/checker"
	"github.com/lazzyfu/goinsight/internal/orders/forms"

	"github.com/gin-gonic/gin"
)

type InspectOrderSyntaxService struct {
	*forms.InspectOrderSyntaxForm
	C        *gin.Context
	Username string
}

// 获取实例配置
func (s *InspectOrderSyntaxService) getInstanceConfig() (commonModels.InsightDBConfig, error) {
	// 获取实例配置
	var config commonModels.InsightDBConfig
	tx := global.App.DB.Table("`insight_db_config`").
		Where("instance_id=?", s.InstanceID).
		First(&config)
	if tx.RowsAffected == 0 {
		return config, fmt.Errorf("未找到实例ID为%s的记录", s.InstanceID)
	}
	return config, nil
}

// 审核SQL
func (s *InspectOrderSyntaxService) inspectSQL(instanceCfg commonModels.InsightDBConfig) ([]checker.ReturnData, error) {
	plainPassword, err := utils.Decrypt(instanceCfg.Password)
	if err != nil {
		return nil, err
	}
	inspect := checker.SyntaxInspectService{
		C:          s.C,
		InstanceID: instanceCfg.InstanceID,
		DbUser:     instanceCfg.User,
		DbPassword: plainPassword,
		DbHost:     instanceCfg.Hostname,
		DbPort:     instanceCfg.Port,
		DBSchema:   s.Schema,
		Username:   s.Username,
		SqlText:    s.Content,
	}
	return inspect.Run()
}

func (s *InspectOrderSyntaxService) Run() (any, error) {
	// 判断SQL类型是否匹配，DML工单仅允许提交DML语句，DDL工单仅允许提交DDL语句
	err := parser.CheckSqlType(s.Content, string(s.SQLType))
	if err != nil {
		return nil, err
	}
	if s.SQLType == "EXPORT" {
		// 导出工单仅检查语法是否有效，不审核，CheckSqlType已经判断类型为SELECT了
		return nil, nil
	}
	// clickhouse不审核
	if s.DBType == "ClickHouse" {
		return nil, nil
	}

	// 获取实例配置
	config, err := s.getInstanceConfig()
	if err != nil {
		return nil, err
	}

	// 审核
	returnData, err := s.inspectSQL(config)
	if err != nil {
		return nil, err
	}

	// 检查语法检查是否通过
	// status: 0表示语法检查通过，1表示语法检查不通过
	status := 0
	for _, row := range returnData {
		for _, sum := range row.Summary {
			if sum.Level != "INFO" {
				status = 1
				break
			}
		}
		if status == 1 {
			break
		}
	}

	return map[string]any{"status": status, "data": returnData}, nil
}
