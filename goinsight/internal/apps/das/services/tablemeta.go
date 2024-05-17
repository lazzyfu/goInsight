/*
@Time    :   2023/05/10 17:03:25
@Author  :   xff
@Desc    : 	 获取表的元信息
*/

package services

import (
	"errors"
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/das/dao"
	"goInsight/internal/apps/das/forms"
	"goInsight/internal/apps/das/parser"
	"strings"

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

func (s *GetTableInfoService) getConfigFromInstanceID() (hostname string, port int, err error) {
	// 获取DB配置
	type DASConfigResult struct {
		Hostname string `json:"hostname"`
		Port     int    `json:"port"`
	}
	var result DASConfigResult
	r := global.App.DB.Table("`insight_db_config` a").
		Select("a.`hostname`, a.`port`").
		Where("a.instance_id=?", s.InstanceID).
		Take(&result)
	// 判断记录是否存在
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return hostname, port, fmt.Errorf("指定DB配置的记录不存在,错误的信息:%s", r.Error.Error())
	}
	return result.Hostname, result.Port, nil
}

func (s GetTableInfoService) getDbType() (string, error) {
	// 获取DB类型
	type dbTypeResult struct {
		DbType string `json:"db_type"`
	}
	var result dbTypeResult
	r := global.App.DB.Table("`insight_db_config` a").
		Select("a.`db_type`").
		Where("a.instance_id=?", s.InstanceID).
		Take(&result)
	// 判断记录是否存在
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("指定DB配置的记录不存在,错误信息:%s", r.Error.Error())
	}
	return result.DbType, nil
}

func (s *GetTableInfoService) getTableStruc(dbType string, hostname string, port int) (data *[]map[string]interface{}, err error) {
	if strings.EqualFold(dbType, "mysql") || strings.EqualFold(dbType, "tidb") {
		db := dao.DB{
			User:     global.App.Config.RemoteDB.UserName,
			Password: global.App.Config.RemoteDB.Password,
			Host:     hostname,
			Port:     port,
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
			User:     global.App.Config.RemoteDB.UserName,
			Password: global.App.Config.RemoteDB.Password,
			Host:     hostname,
			Port:     port,
			Ctx:      s.C.Request.Context(),
		}
		_, data, err = db.Query(fmt.Sprintf("show create table `%s`.`%s`", s.Schema, s.Table))
		if err != nil {
			global.App.Log.Error(err.Error())
		}
	}
	return data, err
}

func (s *GetTableInfoService) getTableBase(dbType string, hostname string, port int) (data *[]map[string]interface{}, err error) {
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
			User:     global.App.Config.RemoteDB.UserName,
			Password: global.App.Config.RemoteDB.Password,
			Host:     hostname,
			Port:     port,
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
	hostname, port, err := s.getConfigFromInstanceID()
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
	// 获取DB类型
	dbType, err := s.getDbType()
	if err != nil {
		return responseData, err
	}
	if s.Type == "structure" {
		// 获取表结构
		return s.getTableStruc(dbType, hostname, port)
	} else {
		// 获取表基础信息
		return s.getTableBase(dbType, hostname, port)
	}
}
