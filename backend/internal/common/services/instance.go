package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/common/forms"
	"github.com/lazzyfu/goinsight/internal/common/models"
	inspectModel "github.com/lazzyfu/goinsight/internal/inspect/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AdminGetInstancesServices struct {
	*forms.AdminGetInstancesForm
	C *gin.Context
}

func (s *AdminGetInstancesServices) Run() (responseData interface{}, total int64, err error) {
	type DBConfig struct {
		models.InsightInstances
		EnvironmentName  string `json:"environment_name"`
		OrganizationName string `json:"organization_name"`
		OrganizationKey  string `json:"organization_key"`
	}
	var dbs []DBConfig
	tx := global.App.DB.Select(`a.id,a.instance_id,a.hostname,a.port,a.user,a.use_type,a.db_type,a.organization_path,b.id as environment, 
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
											insight_orgs ia
										WHERE
											EXISTS (
												SELECT
													1
												FROM
													insight_orgs
												WHERE
													JSON_CONTAINS(c.path, CONCAT('\"', ia.key, '\"'))
											)
									),
									'/',
									c.name
								),
								c.name
							) as organization_name, a.created_at, a.updated_at`).
		Table("insight_instances a").
		Joins("left join insight_instance_environments b on a.environment=b.id").
		Joins("left join insight_orgs c on a.organization_key=c.key")
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

type AdminCreateInstancesService struct {
	*forms.AdminCreateInstancesForm
	C *gin.Context
}

func (s *AdminCreateInstancesService) Run() error {
	// 组织KEY
	organizationKey, err := json.Marshal(s.OrganizationKey)
	if err != nil {
		return err
	}
	organizationKeyJson := datatypes.JSON(organizationKey)
	// 加密密码
	encryptPassword, err := utils.Encrypt(s.Password)
	if err != nil {
		return err
	}
	// 新增记录
	db := models.InsightInstances{
		Hostname:         s.Hostname,
		Port:             s.Port,
		User:             s.User,
		Password:         encryptPassword,
		UseType:          s.UseType,
		DbType:           s.DbType,
		Environment:      s.Environment,
		OrganizationKey:  s.OrganizationKey[len(s.OrganizationKey)-1],
		OrganizationPath: organizationKeyJson,
		Remark:           s.Remark,
		InstanceID:       uuid.New(),
	}
	tx := global.App.DB.Model(&models.InsightInstances{})
	result := tx.Create(&db)

	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("使用类型为%s的%s:%d记录已存在", s.UseType, s.Hostname, s.Port)
			}
		}
		return result.Error
	}
	return nil
}

type AdminUpdateInstancesService struct {
	*forms.AdminUpdateInstancesForm
	C  *gin.Context
	ID uint64
}

func (s *AdminUpdateInstancesService) Run() error {
	// 组织KEY
	organizationKey, err := json.Marshal(s.OrganizationKey)
	if err != nil {
		return err
	}
	organizationKeyJson := datatypes.JSON(organizationKey)

	updates := map[string]interface{}{
		"hostname":          s.Hostname,
		"port":              s.Port,
		"user":              s.User,
		"use_type":          s.UseType,
		"db_type":           s.DbType,
		"environment":       s.Environment,
		"organization_key":  s.OrganizationKey[len(s.OrganizationKey)-1],
		"organization_path": organizationKeyJson,
		"remark":            s.Remark,
	}

	// 仅当传入非空 password 时才更新密码
	if s.Password != "" {
		encryptPassword, err := utils.Encrypt(s.Password)
		if err != nil {
			return err
		}
		updates["password"] = encryptPassword
	}

	// 更新记录
	result := global.App.DB.Model(&models.InsightInstances{}).Where("id=?", s.ID).Updates(updates)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("使用类型为%s的%s:%d记录已存在", s.UseType, s.Hostname, s.Port)
			}
		}
		return result.Error
	}
	return nil
}

type AdminDeleteInstanceConfigsService struct {
	C  *gin.Context
	ID uint64
}

func (s *AdminDeleteInstanceConfigsService) Run() error {
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightInstances{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

type AdminGetInstanceInspectParamsService struct {
	*forms.AdminGetInstanceInspectParamsForm
	C *gin.Context
}

func (s *AdminGetInstanceInspectParamsService) Run() (responseData any, total int64, err error) {
	var records []inspectModel.InsightInspectInstanceParams
	tx := global.App.DB.Model(&inspectModel.InsightInspectInstanceParams{}).Where("instance_id=?", s.InstanceID)
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`title` like ? ", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &records)
	return &records, total, nil
}

type AdminCreateInstanceInspectParamsService struct {
	*forms.AdminCreateInstanceInspectParamsForm
	C *gin.Context
}

func (s *AdminCreateInstanceInspectParamsService) Run() error {
	instanceID, err := uuid.Parse(s.InstanceID)
	if err != nil {
		return fmt.Errorf("invalid instance_id: %w", err)
	}

	// 新增记录
	db := inspectModel.InsightInspectInstanceParams{
		Title:      s.Title,
		Type:       models.EnumType(s.Type),
		Key:        s.Key,
		Value:      s.Value,
		InstanceID: instanceID,
	}
	tx := global.App.DB.Model(&inspectModel.InsightInspectInstanceParams{})
	result := tx.Create(&db)

	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("实例审核参数`%s`已存在", s.Key)
			}
		}
		return result.Error
	}
	return nil
}

type AdminUpdateInstanceInspectParamsService struct {
	*forms.AdminUpdateInstanceInspectParamsForm
	C  *gin.Context
	ID uint64
}

func (s *AdminUpdateInstanceInspectParamsService) Run() error {
	// 只修改value
	tx := global.App.DB.Model(&inspectModel.InsightInspectInstanceParams{}).Where("id=? and instance_id=? and `key`=?", s.ID, s.InstanceID, s.Key)
	result := tx.Updates(map[string]any{
		"value": s.Value,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type AdminDeleteInstanceInspectParamsService struct {
	C  *gin.Context
	ID uint64
}

func (s *AdminDeleteInstanceInspectParamsService) Run() error {
	tx := global.App.DB.Where("id=?", s.ID).Delete(&inspectModel.InsightInspectInstanceParams{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
