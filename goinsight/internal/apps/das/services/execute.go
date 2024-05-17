/*
@Time    :   2023/04/11 16:04:55
@Author  :   xff
@Desc    :
*/

package services

import (
	"context"
	"errors"
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/das/dao"
	"goInsight/internal/apps/das/forms"
	"goInsight/internal/apps/das/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 获取DB配置
func GetDBConfig(instance_id string) (hostname string, port int, err error) {
	type DASConfigResult struct {
		Hostname string `json:"hostname"`
		Port     int    `json:"port"`
	}
	var result DASConfigResult
	r := global.App.DB.Table("`insight_db_config` a").
		Select("a.`hostname`, a.`port`").
		Where("a.instance_id=?", instance_id).
		Take(&result)
	// 判断记录是否存在
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return hostname, port, fmt.Errorf("指定DB配置的记录不存在，错误的信息:%s", r.Error.Error())
	}
	return result.Hostname, result.Port, nil
}

// 获取DB类型
func GetDbType(instance_id string) (string, error) {
	type dbTypeResult struct {
		DbType string `json:"db_type"`
	}
	var result dbTypeResult
	r := global.App.DB.Table("`insight_db_config` a").
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

// 判断当前用户在max_execution_time内是否在对当前实例当前库进行查询，如果有，禁止执行，防止并发
func IsConcurrentRunning(username, instance_id, schema string) error {
	// 逻辑1，返回响应会更新IsRunning=false
	// 定时任务会定时去目标库按照request_id检索正在运行的SQL，当超过阈值，会自动kill并更新IsRunning=false
	type isRunning struct {
		Count int
	}
	var result isRunning
	global.App.DB.Model(&models.InsightDASRecords{}).
		Select("count(*) as count").
		Where("username=? and instance_id=? and `schema`=? and is_finish=? and created_at>= date_sub(now(), interval ? second)",
			username, instance_id, schema, false, global.App.Config.Das.MaxExecutionTime/1000).
		Take(&result)
	if result.Count > 0 {
		return fmt.Errorf("您有`%s`库的查询正在执行,请等待当前SQL执行完成或%ds后重试", schema, global.App.Config.Das.MaxExecutionTime/1000)
	}
	return nil
}

// 计算延时
func CalculateDuration(host string, port int, callback func(string, int) (*[]string, *[]map[string]interface{}, error)) (*[]string, *[]map[string]interface{}, int64, error) {
	startTime := time.Now()
	columns, data, err := callback(host, port)
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
	Execute(hostname string, port int) (*[]string, *[]map[string]interface{}, error)
}

type ClickHouseExecuteApi struct {
	*forms.ExecuteClickHouseQueryForm
	Ctx context.Context
}

func (m ClickHouseExecuteApi) Execute(hostname string, port int) (*[]string, *[]map[string]interface{}, error) {
	db := dao.ClickhouseDB{
		User:     global.App.Config.RemoteDB.UserName,
		Password: global.App.Config.RemoteDB.Password,
		Host:     hostname,
		Port:     port,
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

func (m MySQLExecuteApi) Execute(hostname string, port int) (*[]string, *[]map[string]interface{}, error) {
	db := dao.DB{
		User:     global.App.Config.RemoteDB.UserName,
		Password: global.App.Config.RemoteDB.Password,
		Host:     hostname,
		Port:     port,
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
