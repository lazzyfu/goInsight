package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/das/dao"
	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/parser"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTableInfoService struct {
	*forms.GetTableInfoForm
	C        *gin.Context
	Username string
}

func (s *GetTableInfoService) parserUUID() (id uuid.UUID, err error) {
	id, err = uuid.Parse(s.InstanceID)
	if err != nil {
		return id, err
	}
	return id, nil
}

// 验证用户是否有指定schema的权限
func (s *GetTableInfoService) validatePerms(uuid uuid.UUID) error {
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

func (s *GetTableInfoService) getInstanceCfg() (instance models.InsightInstances, err error) {
	// 获取DB配置
	r := global.App.DB.Table("`insight_instances` a").
		Select("a.`hostname`, a.`port`, a.`user`, a.`password`, a.`db_type`").
		Where("a.instance_id=?", s.InstanceID).
		Take(&instance)
	// 判断记录是否存在
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return instance, fmt.Errorf("指定DB配置的记录不存在,错误的信息:%s", r.Error.Error())
	}
	return instance, nil
}

func (s *GetTableInfoService) getTableStruc(instanceCfg models.InsightInstances) (data *[]map[string]interface{}, err error) {
	// 解密密码
	plainPassword, err := utils.Decrypt(instanceCfg.Password)
	if err != nil {
		return nil, err
	}
	// 获取DB类型
	dbType := string(instanceCfg.DbType)
	if strings.EqualFold(dbType, "mysql") || strings.EqualFold(dbType, "tidb") {
		db := dao.DB{
			User:     instanceCfg.User,
			Password: plainPassword,
			Host:     instanceCfg.Hostname,
			Port:     instanceCfg.Port,
			Params:   map[string]string{"group_concat_max_len": "1073741824"},
			Ctx:      s.C.Request.Context(),
		}

		_, data, err = db.Query(fmt.Sprintf("show create table `%s`.`%s`", s.Schema, s.Table))
		if err != nil {
			global.App.Log.Error(err.Error())
		}
	}
	if strings.EqualFold(dbType, "clickhouse") {
		db := dao.ClickhouseDB{
			User:     instanceCfg.User,
			Password: plainPassword,
			Host:     instanceCfg.Hostname,
			Port:     instanceCfg.Port,
			Ctx:      s.C.Request.Context(),
		}
		_, data, err = db.Query(fmt.Sprintf("show create table `%s`.`%s`", s.Schema, s.Table))
		if err != nil {
			global.App.Log.Error(err.Error())
		}
	}
	return data, err
}

func (s *GetTableInfoService) getTableBase(instanceCfg models.InsightInstances) (data *[]map[string]interface{}, err error) {
	// 解密密码
	plainPassword, err := utils.Decrypt(instanceCfg.Password)
	if err != nil {
		return nil, err
	}
	// 获取DB类型
	dbType := string(instanceCfg.DbType)

	if strings.EqualFold(dbType, "mysql") || strings.EqualFold(dbType, "tidb") {
		query := fmt.Sprintf(`
					select
						*
					from 
						information_schema.tables 
					where 
						table_schema='%s' and table_name='%s'
				`, s.Schema, s.Table)

		db := dao.DB{
			User:     instanceCfg.User,
			Password: plainPassword,
			Host:     instanceCfg.Hostname,
			Port:     instanceCfg.Port,
			Params:   map[string]string{"group_concat_max_len": "1073741824"},
			Ctx:      s.C.Request.Context(),
		}

		_, data, err = db.Query(query)
		if err != nil {
			global.App.Log.Error(err.Error())
		}
		return data, err
	}
	return data, fmt.Errorf("%s不支持获取表信息", dbType)
}

func (s *GetTableInfoService) Run() (responseData *[]map[string]interface{}, err error) {
	// 获取DB配置
	instanceCfg, err := s.getInstanceCfg()
	if err != nil {
		return responseData, err
	}
	// 解析UUID
	uuid, err := s.parserUUID()
	if err != nil {
		return responseData, err
	}
	// 验证用户是否有指定schema的权限
	if err = s.validatePerms(uuid); err != nil {
		return responseData, err
	}

	if s.Type == "structure" {
		// 获取表结构
		return s.getTableStruc(instanceCfg)
	} else {
		// 获取表基础信息
		return s.getTableBase(instanceCfg)
	}
}
