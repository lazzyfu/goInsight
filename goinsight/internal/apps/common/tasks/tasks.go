/*
@Time    :   2023/03/17 09:53:21
@Author  :   zongfei.fu
@Desc    :   同步库表
*/

package tasks

import (
	"context"
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/common/models"
	"goInsight/internal/apps/das/dao"
	"goInsight/internal/pkg/utils"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 忽略的库
var ignoredSchemas []string = []string{
	"'PERFORMANCE_SCHEMA'",
	"'INFORMATION_SCHEMA'",
	"'MYSQL'",
	"'SYS'",
	"'SYSTEM'",
	"'PERCONA'",
	"'DM_META'",
	"'DM_HEARTBEAT'",
	"'MYSQL_MONITOR'",
	"'METRICS_SCHEMA'",
	"'TIDB_BINLOG'",
	"'TIDB_LOADER'",
	"'_TEMPORARY_AND_EXTERNAL_TABLES'",
	"'DEFAULT'",
	"'performance_schema'",
	"'information_schema'",
	"'mysql'",
	"'sys'",
	"'system'",
	"'percona'",
	"'dbms_monitor'",
	"'dm_meta'",
	"'dm_heartbeat'",
	"'mysql_monitor'",
	"'metrics_schema'",
	"'tidb_binlog'",
	"'tidb_loader'",
	"'_temporary_and_external_tables'",
	"'default'",
}

// 空库将不会被同步 && 不采集ghost表
var mysqlQuery string = fmt.Sprintf(`
	SELECT 
		SCHEMA_NAME AS TABLE_SCHEMA
	FROM 
		INFORMATION_SCHEMA.SCHEMATA
	WHERE 
		SCHEMA_NAME NOT IN (%s)
	`, strings.Join(ignoredSchemas, ","))

var clickhouseQuery string = fmt.Sprintf(`
	SELECT 
		name AS TABLE_SCHEMA
	FROM 
		system.databases
	WHERE 
		name NOT IN (%s)
`, strings.Join(ignoredSchemas, ","))

// 插入schema数据
func CreateSchemaRecord(instanceID uuid.UUID, row map[string]interface{}) {
	// 如果记录不存在，插入
	var schemas models.InsightDBSchemas
	result := global.App.DB.Table("insight_db_schemas").Where("`instance_id`=? and `schema`=?", instanceID, row["TABLE_SCHEMA"]).First(&schemas)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			schema := models.InsightDBSchemas{InstanceID: instanceID, Schema: row["TABLE_SCHEMA"].(string)}
			global.App.DB.Create(&schema)
		}
	} else {
		// 如果schema删除后又被新建，更新is_deleted状态
		global.App.DB.Model(&models.InsightDBSchemas{}).Where("instance_id = ? and `schema`=?", instanceID, row["TABLE_SCHEMA"]).Update(`IsDeleted`, false)
	}
}

// 将schema记录更新为软删除
func UpdateSchemaRecordAsSoftDel(instanceID uuid.UUID, schema string) {
	// 将指定的schema更新为软删除
	global.App.DB.Model(&models.InsightDBSchemas{}).Where("instance_id=? and `schema`=?", instanceID, schema).Update(`IsDeleted`, true)
}

// 检查源schema是否被删除
func CheckSourceSchemasIsDeleted(instanceID uuid.UUID, data *[]map[string]interface{}) {
	// 获取源schemas
	var sourceSchemas []string
	for _, row := range *data {
		sourceSchemas = append(sourceSchemas, row["TABLE_SCHEMA"].(string))
	}
	// 从库里读取指定cid的schemas
	type Result struct {
		Schema string
	}
	var result []Result
	var localSchemas []string
	global.App.DB.Table("insight_db_schemas").Where("`instance_id`=?", instanceID).Scan(&result)
	for _, i := range result {
		localSchemas = append(localSchemas, i.Schema)
	}
	// 找出源已经删除的schema
	for _, l := range localSchemas {
		if !utils.IsContain(sourceSchemas, l) {
			UpdateSchemaRecordAsSoftDel(instanceID, l)
		}
	}
}

// 同步库表信息
func SyncDBMeta() {
	// 获取配置的数据库
	type Result struct {
		InstanceID uuid.UUID
		Hostname   string
		Port       int
		DbType     string
	}
	var results []Result
	global.App.DB.Table("insight_db_config").Scan(&results)
	// 启动4个并发
	var wg sync.WaitGroup
	ch := make(chan struct{}, 4)
	for _, row := range results {
		ch <- struct{}{}
		wg.Add(1)
		// 获取目标数据库的库表信息
		go func(row Result) {
			defer wg.Done()
			var (
				data *[]map[string]interface{}
				err  error
			)
			if strings.EqualFold(row.DbType, "mysql") || strings.EqualFold(row.DbType, "tidb") {
				db := dao.DB{
					User:     global.App.Config.RemoteDB.UserName,
					Password: global.App.Config.RemoteDB.Password,
					Host:     row.Hostname,
					Port:     row.Port,
					Params:   map[string]string{"group_concat_max_len": "67108864"},
					Ctx:      context.Background(),
				}
				_, data, err = db.Query(mysqlQuery)
				if err != nil {
					global.App.Log.Error(err.Error())
					return
				}
				for _, d := range *data {
					CreateSchemaRecord(row.InstanceID, d)
				}
				// 判断源库表是否被删除
				CheckSourceSchemasIsDeleted(row.InstanceID, data)
			}
			if strings.EqualFold(row.DbType, "clickhouse") {
				db := dao.ClickhouseDB{
					User:     global.App.Config.RemoteDB.UserName,
					Password: global.App.Config.RemoteDB.Password,
					Host:     row.Hostname,
					Port:     row.Port,
					Ctx:      context.Background(),
				}
				_, data, err = db.Query(clickhouseQuery)
				if err != nil {
					global.App.Log.Error(err.Error())
					return
				}
				for _, d := range *data {
					CreateSchemaRecord(row.InstanceID, d)
				}
				// 判断源库表是否被删除
				CheckSourceSchemasIsDeleted(row.InstanceID, data)
			}
			<-ch
		}(row)
	}
	wg.Wait()
}
