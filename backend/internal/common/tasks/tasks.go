package tasks

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/das/dao"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 忽略下面的库
var ignoredSchemas []string = []string{
	"'PERFORMANCE_SCHEMA'",
	"'INFORMATION_SCHEMA'",
	"'performance_schema'",
	"'information_schema'",
	"'MYSQL'",
	"'mysql'",
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
	var schemas models.InsightInstanceSchemas
	result := global.App.DB.Table("insight_instance_schemas").Where("`instance_id`=? and `schema`=?", instanceID, row["TABLE_SCHEMA"]).First(&schemas)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			schema := models.InsightInstanceSchemas{InstanceID: instanceID, Schema: row["TABLE_SCHEMA"].(string)}
			global.App.DB.Create(&schema)
		}
	} else {
		// 如果schema删除后又被新建，更新is_deleted状态
		global.App.DB.Model(&models.InsightInstanceSchemas{}).Where("instance_id = ? and `schema`=?", instanceID, row["TABLE_SCHEMA"]).Update(`IsDeleted`, false)
	}
}

// 将schema记录更新为软删除
func UpdateSchemaRecordAsSoftDel(instanceID uuid.UUID, schema string) {
	// 将指定的schema更新为软删除
	global.App.DB.Model(&models.InsightInstanceSchemas{}).Where("instance_id=? and `schema`=?", instanceID, schema).Update(`IsDeleted`, true)
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
	global.App.DB.Table("insight_instance_schemas").Where("`instance_id`=?", instanceID).Scan(&result)
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

// 从用户定义的远程数据库实例同步库信息
func SyncDBMeta() {
	// 获取数据库配置列表
	var records []models.InsightInstances
	global.App.DB.Table("insight_instances").Scan(&records)
	// 启动4个并发
	var wg sync.WaitGroup
	ch := make(chan struct{}, 4)
	for _, row := range records {
		ch <- struct{}{}
		wg.Add(1)
		// 获取目标数据库的库信息
		go func(row models.InsightInstances) {
			defer func() {
				<-ch
				wg.Done()
			}()

			var (
				data *[]map[string]any
				err  error
			)

			// 执行SQL超时
			ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
			defer cancel()

			// 解密密码
			plainPassword, err := utils.Decrypt(row.Password)
			if err != nil {
				global.App.Log.Error(fmt.Sprintf("从主机%s:%d同步元数据失败，密码解密失败，错误信息：%s", row.Hostname, row.Port, err.Error()))
				return
			}

			switch strings.ToLower(string(row.DbType)) {
			case "mysql", "tidb":
				db := dao.DB{
					User:     row.User,
					Password: plainPassword,
					Host:     row.Hostname,
					Port:     row.Port,
					Params:   map[string]string{"group_concat_max_len": "67108864"},
					Ctx:      ctx,
				}
				_, data, err = db.Query(mysqlQuery)
			case "clickhouse":
				db := dao.ClickhouseDB{
					User:     row.User,
					Password: plainPassword,
					Host:     row.Hostname,
					Port:     row.Port,
					Ctx:      ctx,
				}
				_, data, err = db.Query(clickhouseQuery)
			default:
				global.App.Log.Warn(fmt.Sprintf("不支持的数据库类型：%s", row.DbType))
				return
			}
			if err != nil {
				global.App.Log.Error(fmt.Sprintf("从主机%s:%d同步元数据失败，错误信息：%s", row.Hostname, row.Port, err.Error()))
				return
			}
			if len(*data) == 0 {
				global.App.Log.Warn(fmt.Sprintf("从主机%s:%d同步元数据失败，未发现库记录，请检查账号%s是否有SELECT权限", row.Hostname, row.Port, row.User))
			}
			// 创建元数据记录
			for _, d := range *data {
				global.App.Log.Debug(fmt.Sprintf("从主机%s:%d同步元数据成功，主机实例ID：%s 库：%s", row.Hostname, row.Port, row.InstanceID.String(), d["TABLE_SCHEMA"]))
				CreateSchemaRecord(row.InstanceID, d)
			}
			// 判断源库是否被删除
			CheckSourceSchemasIsDeleted(row.InstanceID, data)
		}(row)
	}
	wg.Wait()
}
