package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/das/dao"
	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/parser"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetDbDictService struct {
	*forms.GetDbDictForm
	C        *gin.Context
	Username string
}

func (s *GetDbDictService) parserUUID() (id uuid.UUID, err error) {
	id, err = uuid.Parse(s.InstanceID)
	if err != nil {
		return id, err
	}
	return id, nil
}

// 验证用户是否有指定schema的权限
func (s *GetDbDictService) validatePerms(uuid uuid.UUID) error {
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

func (s *GetDbDictService) getInstanceCfg() (record models.InsightDBConfig, err error) {
	// 获取DB配置
	r := global.App.DB.Table("`insight_db_config` a").
		Select("a.`hostname`, a.`port`, a.`user`, a.`password`").
		Where("a.instance_id=?", s.InstanceID).
		Take(&record)
	// 判断记录是否存在
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return record, fmt.Errorf("指定DB配置的记录不存在,错误的信息:%s", r.Error.Error())
	}
	return record, nil
}

func (s *GetDbDictService) getDbType() (string, error) {
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

func (s *GetDbDictService) getDbDict(instanceCfg models.InsightDBConfig) (data *[]map[string]interface{}, err error) {
	query := fmt.Sprintf(`
					select
						t.TABLE_NAME,
						t.TABLE_COMMENT,
						t.CREATE_TIME,
						group_concat(
							distinct concat_ws(
								'<b>',
								c.COLUMN_NAME,
								c.COLUMN_TYPE,
								if(c.IS_NULLABLE = 'NO', 'NOT NULL', 'NULL'),
								ifnull(c.COLUMN_DEFAULT, ''),
								ifnull(c.CHARACTER_SET_NAME, ''),
								ifnull(c.COLLATION_NAME, ''),
								ifnull(c.COLUMN_COMMENT, '')
							) separator '<a>'
						) as COLUMNS_INFO,
						group_concat(
							distinct concat_ws(
								'<b>',
								s.INDEX_NAME,
								if(s.NON_UNIQUE = 0, '唯一', '不唯一'),
								s.Cardinality,
								s.INDEX_TYPE,
								s.COLUMN_NAME
							) separator '<a>'
						) as INDEXES_INFO
					from
						COLUMNS c
						join TABLES t on c.TABLE_SCHEMA = t.TABLE_SCHEMA
						and c.TABLE_NAME = t.TABLE_NAME
						left join STATISTICS s on c.TABLE_SCHEMA = s.TABLE_SCHEMA
						and c.TABLE_NAME = s.TABLE_NAME
					where
						t.TABLE_SCHEMA = '%s'
					group by
						t.TABLE_NAME,
						t.TABLE_COMMENT,
						t.CREATE_TIME
				`, s.Schema)

	// 解密密码
	plainPassword, err := utils.Decrypt(instanceCfg.Password)
	if err != nil {
		return nil, err
	}

	db := dao.DB{
		User:     instanceCfg.User,
		Password: plainPassword,
		Host:     instanceCfg.Hostname,
		Port:     instanceCfg.Port,
		Database: "information_schema",
		Params:   map[string]string{"group_concat_max_len": "1073741824"},
		Ctx:      s.C.Request.Context(),
	}

	_, data, err = db.Query(query)
	if err != nil {
		global.App.Log.Error(err.Error())
	}
	return data, err
}

func (s *GetDbDictService) Run() (responseData *[]map[string]interface{}, err error) {
	// 解析UUID
	uuid, err := s.parserUUID()
	if err != nil {
		return responseData, err
	}
	// 验证用户是否有指定schema的权限
	if err = s.validatePerms(uuid); err != nil {
		return responseData, err
	}
	// 获取DB配置
	instanceCfg, err := s.getInstanceCfg()
	if err != nil {
		return responseData, err
	}

	// 获取DB类型
	dbType, err := s.getDbType()
	if err != nil {
		return responseData, err
	}
	if strings.EqualFold(dbType, "mysql") || strings.EqualFold(dbType, "tidb") {
		return s.getDbDict(instanceCfg)
	}
	return responseData, fmt.Errorf("%s不支持获取数据字典", dbType)
}
