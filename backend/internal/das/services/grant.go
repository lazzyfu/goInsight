package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"
	"github.com/lazzyfu/goinsight/pkg/utils"

	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/das/dao"
	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type GetUserGrantsService struct {
	*forms.UserGrantsForm
	C        *gin.Context
	Username string
}

func (s *GetUserGrantsService) filter(data []map[string]string) (result []map[string]string) {
	var hasAllowedRule bool = false
	for _, i := range data {
		if i["rule"] == "allow" {
			hasAllowedRule = true
			break
		}
	}
	// 如果有allow规则，就移除deny规则
	if hasAllowedRule {
		for _, i := range data {
			if i["rule"] == "allow" {
				result = append(result, i)
			}
		}
	} else {
		// 如果没有allow规则，就保留deny规则
		return data
	}
	return result
}

func (s *GetUserGrantsService) format(schema, tables string) *map[string]interface{} {
	var returnData map[string]interface{} = map[string]interface{}{}
	returnData["schema"] = schema
	if tables == "" {
		returnData["tables"] = "*"
	} else {
		var tmpTables []map[string]string = []map[string]string{}
		for _, i := range strings.Split(tables, ",") {
			val := strings.Split(i, ":")
			tmpTables = append(tmpTables, map[string]string{"table": val[0], "rule": val[1]})
		}
		returnData["tables"] = s.filter(tmpTables)
	}
	return &returnData
}

func (s *GetUserGrantsService) Run() (*map[string]interface{}, error) {
	type result struct {
		Schema string
		Tables string
	}
	var grantResult result
	global.App.DB.Table("insight_das_user_schema_permissions s").
		Select("s.`schema`, group_concat(concat(t.`table`,':',t.rule)) as tables").
		Joins("left join insight_das_user_table_permissions t on s.`schema` = t.`schema` and s.`instance_id` = t.`instance_id`").
		Where("t.username=? and s.instance_id=? and s.`schema`=?", s.Username, s.InstanceID, s.Schema).
		Group("s.`schema`").
		Scan(&grantResult)
	return s.format(grantResult.Schema, grantResult.Tables), nil
}

type AdminGetSchemasGrantService struct {
	*forms.AdminSchemasGrantForm
	C *gin.Context
}

func (s *AdminGetSchemasGrantService) Run() (responseData interface{}, total int64, err error) {
	type result struct {
		models.InsightDASUserSchemaPermissions
		Hostname    string `json:"hostname"`
		Port        int    `json:"port"`
		DBType      string `json:"db_type"`
		Environment string `json:"environment"`
		Remark      string `json:"remark"`
	}
	var schemaPerms []result
	tx := global.App.DB.Select("a.id, a.username,a.schema,a.instance_id, b.hostname, b.port,b.db_type, c.name as environment, b.remark").
		Table("insight_das_user_schema_permissions a").
		Joins("left join insight_db_config b on a.instance_id=b.instance_id").
		Joins("left join insight_db_environments c on b.environment=c.id").
		Order("a.updated_at desc")
	// 搜索
	if s.Search != "" {
		tx = tx.Where("b.`hostname` like ? or b.`remark` like ? or a.`instance_id` like ? or a.username like ? or a.schema like ?", "%"+s.Search+"%", "%"+s.Search+"%", "%"+s.Search+"%", "%"+s.Search+"%", "%"+s.Search+"%")
	}
	if s.Environment != "" {
		tx = tx.Where("c.name = ?", s.Environment)
	}
	total = pagination.Pager(&s.PaginationQ, tx, &schemaPerms)
	return &schemaPerms, total, nil
}

// 获取指定环境的实例
type AdminGetInstancesListService struct {
	*forms.AdminGetInstancesListForm
	C *gin.Context
}

func (s *AdminGetInstancesListService) Run() (responseData interface{}, total int64, err error) {
	var configs []commonModels.InsightDBConfig
	tx := global.App.DB.Table("insight_db_config").Where("environment=? and db_type=? and use_type='查询'", s.ID, s.DbType)
	total = pagination.Pager(&s.PaginationQ, tx, &configs)
	return &configs, total, nil
}

type AdminGetSchemasListService struct {
	*forms.AdminGetSchemasListForm
	C *gin.Context
}

func (s *AdminGetSchemasListService) Run() (responseData interface{}, total int64, err error) {
	var roles []commonModels.InsightDBSchemas
	tx := global.App.DB.Table("insight_db_schemas").Where("instance_id=?", s.InstanceID)
	total = pagination.Pager(&s.PaginationQ, tx, &roles)
	return &roles, total, nil
}

type AdminGetTablesListService struct {
	*forms.AdminGetTablesListForm
	C *gin.Context
}

// 获取MySQL/TiDB的表信息
func (g *AdminGetTablesListService) getMySQLMetaData(r *ConfigResult) (data *[]map[string]interface{}, err error) {
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
func (g *AdminGetTablesListService) getClickHouseMetaData(r *ConfigResult) (data *[]map[string]interface{}, err error) {
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

func (s *AdminGetTablesListService) Run() (responseData *[]map[string]interface{}, err error) {
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

type AdminCreateSchemasGrantService struct {
	*forms.AdminCreateSchemasGrantForm
	C *gin.Context
}

func (s *AdminCreateSchemasGrantService) Run() (err error) {
	// 解析UUID
	instance_id, err := utils.ParserUUID(s.InstanceID)
	if err != nil {
		return err
	}
	schemaRecord := models.InsightDASUserSchemaPermissions{
		Username:   s.Username,
		Schema:     s.Schema,
		InstanceID: instance_id,
	}
	tableRecords := []models.InsightDASUserTablePermissions{}
	for _, table := range s.Tables {
		tableRecords = append(tableRecords, models.InsightDASUserTablePermissions{
			Username:   s.Username,
			Schema:     s.Schema,
			InstanceID: instance_id,
			Table:      table,
			Rule:       "allow",
		})
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightDASUserSchemaPermissions{}).Create(&schemaRecord).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录已存在，错误:%s", err.Error())
			}
			global.App.Log.Error(err)
			return err
		}
		if err := tx.Model(&models.InsightDASUserTablePermissions{}).CreateInBatches(&tableRecords, len(tableRecords)).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录已存在，错误:%s", err.Error())
			}
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}

type AdminGetTablesGrantService struct {
	*forms.AdminGetTablesGrantForm
	C *gin.Context
}

func (s *AdminGetTablesGrantService) Run() (responseData interface{}, total int64, err error) {
	var tables []models.InsightDASUserTablePermissions
	tx := global.App.DB.Model(&models.InsightDASUserTablePermissions{}).
		Where("username=? and instance_id=? and `schema`=?", s.Username, s.InstanceID, s.Schema).
		Order("created_at asc").
		Scan(&tables)
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`table` like ? or `rule` like ?", "%"+s.Search+"%", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &tables)
	return &tables, total, nil
}

type AdminCreateTablesGrantService struct {
	*forms.AdminCreateTablesGrantForm
	C *gin.Context
}

func (s *AdminCreateTablesGrantService) Run() (err error) {
	// 解析UUID
	instance_id, err := utils.ParserUUID(s.InstanceID)
	if err != nil {
		return err
	}
	tableRecords := []models.InsightDASUserTablePermissions{}
	for _, table := range s.Tables {
		tableRecords = append(tableRecords, models.InsightDASUserTablePermissions{
			Username:   s.Username,
			Schema:     s.Schema,
			InstanceID: instance_id,
			Table:      table,
			Rule:       commonModels.EnumType(s.Rule),
		})
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightDASUserTablePermissions{}).CreateInBatches(&tableRecords, len(tableRecords)).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录已存在，错误:%s", err.Error())
			}
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}

// type UpdateAdminTablesGrantService struct {
// 	*forms.UpdateAdminTablesGrantForm
// 	C *gin.Context
// }

// func (s *UpdateAdminTablesGrantService) BatchDelete(tx *gorm.DB, instance_id uuid.UUID, rule string) error {
// 	// 先删除
// 	tx.Where("username=? and `schema`=? and instance_id=? and rule=?", s.Username, s.Schema, instance_id, rule).Delete(&models.InsightDASUserTablePermissions{})
// 	return nil
// }

// func (s *UpdateAdminTablesGrantService) Run() (err error) {
// 	// 解析UUID
// 	instance_id, err := utils.ParserUUID(s.InstanceID)
// 	if err != nil {
// 		return err
// 	}

// 	allowedTableRecords := []models.InsightDASUserTablePermissions{}
// 	for _, table := range s.AllowedTables {
// 		allowedTableRecords = append(allowedTableRecords, models.InsightDASUserTablePermissions{
// 			Username:   s.Username,
// 			Schema:     s.Schema,
// 			InstanceID: instance_id,
// 			Table:      table,
// 			Rule:       "allow",
// 		})
// 	}

// 	rejectedTableRecords := []models.InsightDASUserTablePermissions{}
// 	for _, table := range s.AllowedTables {
// 		rejectedTableRecords = append(rejectedTableRecords, models.InsightDASUserTablePermissions{
// 			Username:   s.Username,
// 			Schema:     s.Schema,
// 			InstanceID: instance_id,
// 			Table:      table,
// 			Rule:       "deny",
// 		})
// 	}
// 	return global.App.DB.Transaction(func(tx *gorm.DB) error {
// 		if len(s.AllowedTables) > 0 {
// 			fmt.Println("allowedTableRecords....")
// 			_ = s.BatchDelete(tx, instance_id, "allow")
// 			if err := tx.Model(&models.InsightDASUserTablePermissions{}).CreateInBatches(&allowedTableRecords, len(allowedTableRecords)).Error; err != nil {
// 				mysqlErr := err.(*mysql.MySQLError)
// 				switch mysqlErr.Number {
// 				case 1062:
// 					return fmt.Errorf("记录已存在，错误:%s", err.Error())
// 				}
// 				global.App.Log.Error(err)
// 				return err
// 			}
// 		}
// 		if len(s.RejectedTables) > 0 {
// 			fmt.Println("rejectedTableRecords....")
// 			_ = s.BatchDelete(tx, instance_id, "deny")
// 			if err := tx.Model(&models.InsightDASUserTablePermissions{}).CreateInBatches(&rejectedTableRecords, len(rejectedTableRecords)).Error; err != nil {
// 				mysqlErr := err.(*mysql.MySQLError)
// 				switch mysqlErr.Number {
// 				case 1062:
// 					return fmt.Errorf("记录已存在，错误:%s", err.Error())
// 				}
// 				global.App.Log.Error(err)
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

type AdminDeleteSchemasGrantService struct {
	C  *gin.Context
	ID uint32
}

func (s *AdminDeleteSchemasGrantService) Run() error {
	var schemaPerms models.InsightDASUserSchemaPermissions
	tx := global.App.DB.Model(&models.InsightDASUserSchemaPermissions{}).
		Where("id=?", s.ID).
		Take(&schemaPerms)
	if tx.RowsAffected == 0 {
		return errors.New("删除失败，影响行数为0")
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id=?", s.ID).
			Delete(&models.InsightDASUserSchemaPermissions{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		if err := tx.Where("username=? and instance_id=? AND `schema`=?", schemaPerms.Username, schemaPerms.InstanceID, schemaPerms.Schema).
			Delete(&models.InsightDASUserTablePermissions{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}

type AdminDeleteTablesGrantService struct {
	C  *gin.Context
	ID uint32
}

func (s *AdminDeleteTablesGrantService) Run() error {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id=?", s.ID).Delete(&models.InsightDASUserTablePermissions{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}
