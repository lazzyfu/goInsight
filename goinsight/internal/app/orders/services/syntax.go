/*
@Time    :   2023/10/24 19:14:17
@Author  :   lazzyfu
*/

package services

import (
	"goInsight/global"
	commonModels "goInsight/internal/app/common/models"
	"goInsight/internal/app/orders/api"
	"goInsight/internal/app/orders/forms"
	"goInsight/internal/pkg/parser"

	"github.com/gin-gonic/gin"
)

// hook工单
type SyntaxCheckFormService struct {
	*forms.SyntaxCheckForm
	C        *gin.Context
	Username string
}

func (s *SyntaxCheckFormService) Run() (interface{}, error) {
	// 判断SQL类型是否匹配，DML工单仅允许提交DML语句，DDL工单仅允许提交DDL语句
	err := parser.CheckSqlType(s.Content, string(s.SQLType))
	if err != nil {
		return nil, err
	}
	if s.SQLType == "EXPORT" {
		// 导出工单仅检查语法是否有效，不审核
		return nil, nil
	}
	// clickhouse不审核
	if s.DBType == "ClickHouse" {
		return nil, nil
	}
	// 获取实例配置
	var config commonModels.InsightDBConfig
	global.App.DB.Table("`insight_db_config`").
		Where("instance_id=?", s.InstanceID).
		First(&config)

	// DB参数
	api := api.GAuditApi{
		DbUser:            global.App.Config.RemoteDB.UserName,
		DbPassword:        global.App.Config.RemoteDB.Password,
		DbHost:            config.Hostname,
		DbPort:            config.Port,
		DB:                s.Schema,
		Timeout:           3000,
		CustomAuditParams: map[string]interface{}{},
		SqlText:           s.Content,
	}
	return api.Check()
}
