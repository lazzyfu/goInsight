/*
@Time    :   2023/08/03 11:27:25
@Author  :   zongfei.fu
@Desc    :
*/

package services

import (
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/das/dao"
	"goInsight/internal/apps/das/forms"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetOrderTablesService struct {
	*forms.GetOrderTablesForm
	C *gin.Context
}

// 获取MySQL/TiDB的表信息
func (g *GetOrderTablesService) getMySQLMetaData(r *ConfigResult) (data *[]map[string]interface{}, err error) {
	var query string = fmt.Sprintf(`
		select 
			table_name as table_name
		from 
			information_schema.tables
		where 
			table_schema='%s' and table_name not regexp '^_(.*)[_ghc|_gho|_del]$'
		`, r.Schema)
	db := dao.DB{
		User:     global.App.Config.RemoteDB.UserName,
		Password: global.App.Config.RemoteDB.Password,
		Host:     r.Hostname,
		Port:     r.Port,
		Params:   map[string]string{"group_concat_max_len": "4194304"},
		Ctx:      g.C.Request.Context(),
	}
	_, data, err = db.Query(query)
	if err != nil {
		global.App.Log.Error(err.Error())
		return data, err
	}
	return data, nil
}

// 获取ClickHouse的表信息
func (g *GetOrderTablesService) getClickHouseMetaData(r *ConfigResult) (data *[]map[string]interface{}, err error) {
	var query string = fmt.Sprintf(`
	select 
		name as table_name
	from 
		system.tables 
	where 
		(database = '%s')
	`, r.Schema)
	db := dao.ClickhouseDB{
		User:     global.App.Config.RemoteDB.UserName,
		Password: global.App.Config.RemoteDB.Password,
		Host:     r.Hostname,
		Port:     r.Port,
		Ctx:      g.C.Request.Context(),
	}
	_, data, err = db.Query(query)
	if err != nil {
		global.App.Log.Error(err.Error())
		return data, err
	}
	return data, nil
}

func (s *GetOrderTablesService) Run() (responseData *[]map[string]interface{}, err error) {
	var config ConfigResult
	global.App.DB.Table("insight_db_config").Where("`instance_id`=?", s.InstanceID).Take(&config)
	config.Schema = s.Schema
	if strings.EqualFold(config.DbType, "mysql") || strings.EqualFold(config.DbType, "tidb") {
		return s.getMySQLMetaData(&config)
	}
	if strings.EqualFold(config.DbType, "clickhouse") {
		return s.getClickHouseMetaData(&config)
	}

	return responseData, nil
}
