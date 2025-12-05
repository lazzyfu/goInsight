package services

import (
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/das/dao"
	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/parser"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstanceCfg struct {
	Hostname      string
	Port          int
	User          string
	Password      string
	PlainPassword string
	Schema        string
	DbType        string
}

type GetTablesService struct {
	*forms.GetTablesForm
	C        *gin.Context
	Username string
}

func (s *GetTablesService) parserUUID() (id uuid.UUID, err error) {
	id, err = uuid.Parse(s.InstanceID)
	if err != nil {
		return id, err
	}
	return id, nil
}

// 验证用户是否有指定schema的权限
func (s *GetTablesService) validatePerms(uuid uuid.UUID) error {
	// 检查库表权限
	var tables []parser.Table
	tables = append(tables, parser.Table{Schema: s.Schema})
	checker := CheckUserPerm{
		UserName:   s.Username,
		InstanceID: uuid,
		Tables:     tables,
	}
	if err := checker.HasSchemaPerms(); err != nil {
		return err
	}
	return nil
}

// 获取MySQL/TiDB的表&字段信息
func (s *GetTablesService) getMySQLMetaData(instanceCfg *InstanceCfg) (data *[]map[string]interface{}, err error) {
	/*
		table_name: test
		table_schema: pmm
		columns: id$$int(10) unsigned$$@@created_at$$datetime(3)$$创建时间@@updated_at$$datetime(3)$$更新时间
	*/
	var query string = fmt.Sprintf(`
		SELECT 
			table_schema as table_schema,
			table_name as table_name, 
			group_concat(concat(column_name, '$$', column_type, "$$", column_comment) SEPARATOR '@@') as columns
		from 
			information_schema.columns
		where 
			table_schema='%s' and table_name not regexp '^_(.*)[_ghc|_gho|_del]$'
		group by 
			table_schema, table_name order by table_name
		`, instanceCfg.Schema)
	db := dao.DB{
		User:     instanceCfg.User,
		Password: instanceCfg.PlainPassword,
		Host:     instanceCfg.Hostname,
		Port:     instanceCfg.Port,
		Params:   map[string]string{"group_concat_max_len": "4194304"},
		Ctx:      s.C.Request.Context(),
	}
	_, data, err = db.Query(query)
	if err != nil {
		global.App.Log.Error(err.Error())
		return data, err
	}
	return data, nil
}

// 获取ClickHouse的表&字段信息
func (s *GetTablesService) getClickHouseMetaData(instanceCfg *InstanceCfg) (data *[]map[string]interface{}, err error) {
	/*
		table_name: test
		table_schema: pmm
		columns: I_ID#bigint(18),USER_ID#bigint(20),D_CREATED_AT#datetime
	*/
	var query string = fmt.Sprintf(`
		SELECT
			database as table_schema,
			table as table_name,
			columns
		FROM
		(
			SELECT
				database,
				table,
				groupArray(concat_col) AS col_array,
				arrayStringConcat(col_array, '@@') AS columns
			FROM
				(
					SELECT
						database,
						table,
						concat(name, '$$', type, '$$', comment) AS concat_col,
						name AS column
					FROM
						system.columns
					WHERE
						(database = '%s')
				)
			GROUP BY
				database,
				table
		)
		ORDER BY table ASC
	`, instanceCfg.Schema)

	db := dao.ClickhouseDB{
		User:     instanceCfg.User,
		Password: instanceCfg.PlainPassword,
		Host:     instanceCfg.Hostname,
		Port:     instanceCfg.Port,
		Ctx:      s.C.Request.Context(),
	}
	_, data, err = db.Query(query)
	if err != nil {
		global.App.Log.Error(err.Error())
		return data, err
	}
	return data, nil
}

func (s *GetTablesService) Run() (responseData *[]map[string]interface{}, err error) {
	// 解析UUID
	uuid, err := s.parserUUID()
	if err != nil {
		return responseData, err
	}
	// 验证用户是否有指定schema的权限
	if err = s.validatePerms(uuid); err != nil {
		return responseData, err
	}
	var cfg InstanceCfg
	global.App.DB.Table("insight_db_config").Where("`instance_id`=?", s.InstanceID).Take(&cfg)
	// 解密密码
	plainPassword, err := utils.Decrypt(cfg.Password)
	if err != nil {
		return nil, err
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
