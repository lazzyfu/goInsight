package services

import (
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/das/dao"
	"github.com/lazzyfu/goinsight/internal/das/forms"

	"github.com/gin-gonic/gin"
)

type GetOrderTablesService struct {
	*forms.GetOrderTablesForm
	C *gin.Context
}

// 获取MySQL/TiDB的表信息
func (g *GetOrderTablesService) getMySQLMetaData(r *InstanceCfg) (data *[]map[string]interface{}, err error) {
	var query string = fmt.Sprintf(`
		select 
			table_name as table_name
		from 
			information_schema.tables
		where 
			table_schema='%s' and table_name not regexp '^_(.*)[_ghc|_gho|_del]$'
		`, r.Schema)
	db := dao.DB{
		User:     r.User,
		Password: r.PlainPassword,
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
func (g *GetOrderTablesService) getClickHouseMetaData(r *InstanceCfg) (data *[]map[string]interface{}, err error) {
	var query string = fmt.Sprintf(`
	select 
		name as table_name
	from 
		system.tables 
	where 
		(database = '%s')
	`, r.Schema)

	db := dao.ClickhouseDB{
		User:     r.User,
		Password: r.PlainPassword,
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
	var cfg InstanceCfg
	global.App.DB.Table("insight_db_config").Where("`instance_id`=?", s.InstanceID).Take(&cfg)
	plainPassword, err := utils.Decrypt(cfg.Password)
	if err != nil {
		return
	}
	cfg.PlainPassword = plainPassword
	cfg.Schema = s.Schema
	if strings.EqualFold(cfg.DbType, "mysql") || strings.EqualFold(cfg.DbType, "tidb") {
		return s.getMySQLMetaData(&cfg)
	}
	if strings.EqualFold(cfg.DbType, "clickhouse") {
		return s.getClickHouseMetaData(&cfg)
	}

	return responseData, nil
}
