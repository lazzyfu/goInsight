/*
@Time    :   2023/03/17 14:50:56
@Author  :   zongfei.fu
@Desc    :   获取用户有权限的库
*/

package services

import (
	"goInsight/global"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type schemaResult struct {
	InstanceID uuid.UUID `json:"instance_id"`
	Schema     string    `json:"schema"`
	DbType     string    `json:"db_type"`
	Hostname   string    `json:"hostname"`
	Port       int       `json:"port"`
	IsDeleted  bool      `json:"is_deleted"`
	Remark     string    `json:"remark"`
}

type GetSchemasService struct {
	C        *gin.Context
	Username string
}

// 获取用户授权的库
func (s *GetSchemasService) Run() ([]schemaResult, error) {
	var results []schemaResult
	global.App.DB.Table("`insight_das_user_schema_permissions` a").
		Select("a.`instance_id`, b.`db_type`, a.`schema`, b.`hostname`, b.`port`, c.`is_deleted`, b.`remark`").
		Joins("join `insight_db_config` b on a.instance_id = b.instance_id").
		Joins("join `insight_db_schemas` c on b.instance_id = c.instance_id and a.`schema` = c.`schema`").
		Where("a.username=?", s.Username).
		Group("a.`instance_id`, b.`db_type`, a.`schema`, b.`hostname`, b.`port`, c.is_deleted, b.`remark`").
		Scan(&results)
	return results, nil
}
