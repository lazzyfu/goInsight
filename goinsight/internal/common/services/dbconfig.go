package services

import (
	"encoding/json"
	"fmt"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/common/forms"
	"github.com/lazzyfu/goinsight/internal/common/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AdminGetDBConfigServices struct {
	*forms.AdminGetDBConfigForm
	C *gin.Context
}

func (s *AdminGetDBConfigServices) Run() (responseData interface{}, total int64, err error) {
	type DBConfig struct {
		models.InsightDBConfig
		EnvironmentName  string `json:"environment_name"`
		OrganizationName string `json:"organization_name"`
		OrganizationKey  string `json:"organization_key"`
	}
	var dbs []DBConfig
	tx := global.App.DB.Select(`a.id,a.instance_id,a.hostname,a.port,a.use_type,a.db_type,a.inspect_params,a.organization_path,b.id as environment, 
							b.name as environment_name, a.remark, ifnull(
								concat(
									(
										SELECT
											GROUP_CONCAT(
												ia.name
												ORDER BY
													ia.name ASC SEPARATOR '/'
											) AS concatenated_names
										FROM
											insight_organizations ia
										WHERE
											EXISTS (
												SELECT
													1
												FROM
													insight_organizations
												WHERE
													JSON_CONTAINS(c.path, CONCAT('\"', ia.key, '\"'))
											)
									),
									'/',
									c.name
								),
								c.name
							) as organization_name, a.created_at, a.updated_at`).
		Table("insight_db_config a").
		Joins("left join insight_db_environments b on a.environment=b.id").
		Joins("left join insight_organizations c on a.organization_key=c.key")
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`hostname` like ? or `remark` like ? or `instance_id` like ?", "%"+s.Search+"%", "%"+s.Search+"%", "%"+s.Search+"%")
	}
	if s.Environment != "" {
		tx = tx.Where("b.name=?", s.Environment)
	}
	if s.DbType != "" {
		tx = tx.Where("a.db_type=?", s.DbType)
	}
	total = pagination.Pager(&s.PaginationQ, tx, &dbs)
	return &dbs, total, nil
}

type AdminCreateDBConfigService struct {
	*forms.AdminCreateDBConfigForm
	C *gin.Context
}

func (s *AdminCreateDBConfigService) Run() error {
	// 组织KEY
	organizationKey, err := json.Marshal(s.OrganizationKey)
	if err != nil {
		return err
	}
	organizationKeyJson := datatypes.JSON(organizationKey)

	// 审核参数
	jsonInspectParams, err := json.Marshal(s.InspectParams)
	if err != nil {
		return err
	}
	// 新增记录
	db := models.InsightDBConfig{
		Hostname:         s.Hostname,
		Port:             s.Port,
		InspectParams:    datatypes.JSON(jsonInspectParams),
		UseType:          s.UseType,
		DbType:           s.DbType,
		Environment:      s.Environment,
		OrganizationKey:  s.OrganizationKey[len(s.OrganizationKey)-1],
		OrganizationPath: organizationKeyJson,
		Remark:           s.Remark,
		InstanceID:       uuid.New(),
	}
	tx := global.App.DB.Model(&models.InsightDBConfig{})
	result := tx.Create(&db)

	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("使用类型为%s的%s:%d记录已存在", s.UseType, s.Hostname, s.Port)
		}
		return result.Error
	}
	return nil
}

type AdminUpdateDBConfigService struct {
	*forms.AdminUpdateDBConfigForm
	C  *gin.Context
	ID uint64
}

func (s *AdminUpdateDBConfigService) Run() error {
	// 组织KEY
	organizationKey, err := json.Marshal(s.OrganizationKey)
	if err != nil {
		return err
	}
	organizationKeyJson := datatypes.JSON(organizationKey)

	// 审核参数
	jsonInspectParams, err := json.Marshal(s.InspectParams)
	if err != nil {
		return err
	}
	// 更新记录
	result := global.App.DB.Model(&models.InsightDBConfig{}).Where("id=?", s.ID).Updates(map[string]interface{}{
		"hostname":          s.Hostname,
		"port":              s.Port,
		"inspect_params":    datatypes.JSON(jsonInspectParams),
		"use_type":          s.UseType,
		"db_type":           s.DbType,
		"environment":       s.Environment,
		"organization_key":  s.OrganizationKey[len(s.OrganizationKey)-1],
		"organization_path": organizationKeyJson,
		"remark":            s.Remark,
	})
	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("使用类型为%s的%s:%d记录已存在", s.UseType, s.Hostname, s.Port)
		}
		return result.Error
	}
	return nil
}

type AdminDeleteDBConfigService struct {
	C  *gin.Context
	ID uint64
}

func (s *AdminDeleteDBConfigService) Run() error {
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightDBConfig{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
