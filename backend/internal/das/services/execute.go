package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"

	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/das/dao"
	"github.com/lazzyfu/goinsight/internal/das/forms"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 获取DB配置
func GetDBConfig(instance_id string) (instanceCfg commonModels.InsightInstances, err error) {
	r := global.App.DB.Table("`insight_instances` a").
		Select("a.`hostname`, a.`port`, a.`User`, a.`Password`").
		Where("a.instance_id=?", instance_id).
		Take(&instanceCfg)
	// 判断记录是否存在
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return instanceCfg, fmt.Errorf("指定DB配置的记录不存在，错误的信息:%s", r.Error.Error())
	}
	return instanceCfg, nil
}

// 获取DB类型
func GetDbType(instance_id string) (string, error) {
	type dbTypeResult struct {
		DbType string `json:"db_type"`
	}
	var result dbTypeResult
	r := global.App.DB.Table("`insight_instances` a").
		Select("a.`db_type`").
		Where("a.instance_id=?", instance_id).
		Take(&result)
	// 判断记录是否存在
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("指定DB配置的记录不存在，错误信息:%s", r.Error.Error())
	}
	return result.DbType, nil
}

// 解析UUID
func ParserUUID(instance_id string) (id uuid.UUID, err error) {
	id, err = uuid.Parse(instance_id)
	if err != nil {
		return id, err
	}
	return id, nil
}

// 计算延时
func CalculateDuration(instanceCfg commonModels.InsightInstances, callback func(instanceCfg commonModels.InsightInstances) (*[]string, *[]map[string]interface{}, error)) (*[]string, *[]map[string]interface{}, int64, error) {
	startTime := time.Now()
	columns, data, err := callback(instanceCfg)
	endTime := time.Now()
	return columns, data, int64(endTime.Sub(startTime) / time.Millisecond), err
}

// 响应data
type ResponseData struct {
	SQLText  string                    `json:"sqltext"`
	Duration string                    `json:"duration"`
	QueryID  string                    `json:"query_id"`
	Data     *[]map[string]interface{} `json:"data"`
	Columns  *[]string                 `json:"columns"`
}

// 执行接口
type ExecuteApi interface {
	Execute(instanceCfg commonModels.InsightInstances) (*[]string, *[]map[string]interface{}, error)
}

type ClickHouseExecuteApi struct {
	*forms.ExecuteClickHouseQueryForm
	Ctx context.Context
}

func (m ClickHouseExecuteApi) Execute(instanceCfg commonModels.InsightInstances) (*[]string, *[]map[string]interface{}, error) {
	// 解密密码
	plainPassword, err := utils.Decrypt(instanceCfg.Password)
	if err != nil {
		return nil, nil, err
	}

	db := dao.ClickhouseDB{
		User:     instanceCfg.User,
		Password: plainPassword,
		Host:     instanceCfg.Hostname,
		Port:     instanceCfg.Port,
		Database: m.Schema,
		Settings: m.Params,
		Ctx:      m.Ctx,
	}
	columns, data, err := db.Query(m.SQLText)
	if err != nil {
		return nil, nil, err
	}
	return columns, data, nil
}

type MySQLExecuteApi struct {
	*forms.ExecuteMySQLQueryForm
	Ctx context.Context
}

func (m MySQLExecuteApi) Execute(instanceCfg commonModels.InsightInstances) (*[]string, *[]map[string]interface{}, error) {
	// 解密密码
	plainPassword, err := utils.Decrypt(instanceCfg.Password)
	if err != nil {
		return nil, nil, err
	}

	db := dao.DB{
		User:     instanceCfg.User,
		Password: plainPassword,
		Host:     instanceCfg.Hostname,
		Port:     instanceCfg.Port,
		Database: m.Schema,
		Params:   m.Params,
		Ctx:      m.Ctx,
	}
	columns, data, err := db.Query(m.SQLText)
	if err != nil {
		return nil, nil, err
	}
	return columns, data, nil
}
