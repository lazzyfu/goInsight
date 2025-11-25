package services

import (
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/internal/das/forms"

	"github.com/gin-gonic/gin"
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
