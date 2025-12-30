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

type AdminGetSchemasGrantService struct {
	*forms.AdminSchemasGrantForm
	C *gin.Context
}

func (s *AdminGetSchemasGrantService) Run() (responseData any, total int64, err error) {
	type result struct {
		models.InsightDasSchemaPerms
		Hostname    string `json:"hostname"`
		Port        int    `json:"port"`
		DBType      string `json:"db_type"`
		Environment string `json:"environment"`
		Remark      string `json:"remark"`
	}
	var schemaPerms []result
	tx := global.App.DB.Select("a.id, a.username,a.schema,a.instance_id, b.hostname, b.port,b.db_type, c.name as environment, b.remark").
		Table("insight_das_schema_perms a").
		Joins("left join insight_instances b on a.instance_id=b.instance_id").
		Joins("left join insight_instance_environments c on b.environment=c.id").
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

func (s *AdminGetInstancesListService) Run() (responseData any, total int64, err error) {
	var configs []commonModels.InsightInstances
	tx := global.App.DB.Table("insight_instances").Where("environment=? and db_type=? and use_type='查询'", s.ID, s.DbType)
	total = pagination.Pager(&s.PaginationQ, tx, &configs)
	return &configs, total, nil
}

type AdminGetSchemasListService struct {
	*forms.AdminGetSchemasListForm
	C *gin.Context
}

func (s *AdminGetSchemasListService) Run() (responseData any, total int64, err error) {
	var roles []commonModels.InsightInstanceSchemas
	tx := global.App.DB.Table("insight_instance_schemas").Where("instance_id=?", s.InstanceID)
	total = pagination.Pager(&s.PaginationQ, tx, &roles)
	return &roles, total, nil
}

type AdminGetTablesListService struct {
	*forms.AdminGetTablesListForm
	C *gin.Context
}

// 获取MySQL/TiDB的表信息
func (g *AdminGetTablesListService) getMySQLMetaData(r *InstanceCfg) (data *[]map[string]any, err error) {
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
func (g *AdminGetTablesListService) getClickHouseMetaData(r *InstanceCfg) (data *[]map[string]any, err error) {
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

func (s *AdminGetTablesListService) Run() (responseData *[]map[string]any, err error) {
	var cfg InstanceCfg
	global.App.DB.Table("insight_instances").Where("`instance_id`=?", s.InstanceID).Take(&cfg)

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
	schemaRecord := models.InsightDasSchemaPerms{
		Username:   s.Username,
		Schema:     s.Schema,
		InstanceID: instance_id,
	}
	tableRecords := []models.InsightDasTablePerms{}
	for _, table := range s.Tables {
		tableRecords = append(tableRecords, models.InsightDasTablePerms{
			Username:   s.Username,
			Schema:     s.Schema,
			InstanceID: instance_id,
			Table:      table,
			Rule:       "allow",
		})
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightDasSchemaPerms{}).Create(&schemaRecord).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录已存在，错误:%s", err.Error())
			}
			global.App.Log.Error(err)
			return err
		}
		if err := tx.Model(&models.InsightDasTablePerms{}).CreateInBatches(&tableRecords, len(tableRecords)).Error; err != nil {
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

func (s *AdminGetTablesGrantService) Run() (responseData any, total int64, err error) {
	var tables []models.InsightDasTablePerms
	tx := global.App.DB.Model(&models.InsightDasTablePerms{}).
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
	tableRecords := []models.InsightDasTablePerms{}
	for _, table := range s.Tables {
		tableRecords = append(tableRecords, models.InsightDasTablePerms{
			Username:   s.Username,
			Schema:     s.Schema,
			InstanceID: instance_id,
			Table:      table,
			Rule:       commonModels.EnumType(s.Rule),
		})
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightDasTablePerms{}).CreateInBatches(&tableRecords, len(tableRecords)).Error; err != nil {
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

type AdminDeleteSchemasGrantService struct {
	C  *gin.Context
	ID uint32
}

func (s *AdminDeleteSchemasGrantService) Run() error {
	var schemaPerms models.InsightDasSchemaPerms
	tx := global.App.DB.Model(&models.InsightDasSchemaPerms{}).
		Where("id=?", s.ID).
		Take(&schemaPerms)
	if tx.RowsAffected == 0 {
		return errors.New("删除失败，影响行数为0")
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id=?", s.ID).
			Delete(&models.InsightDasSchemaPerms{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		if err := tx.Where("username=? and instance_id=? AND `schema`=?", schemaPerms.Username, schemaPerms.InstanceID, schemaPerms.Schema).
			Delete(&models.InsightDasTablePerms{}).Error; err != nil {
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
		if err := tx.Where("id=?", s.ID).Delete(&models.InsightDasTablePerms{}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}
