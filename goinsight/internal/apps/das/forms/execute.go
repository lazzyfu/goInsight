/*
@Time    :   2023/03/21 14:20:29
@Author  :   xff
@Desc    :
*/

package forms

type ExecuteMySQLQueryForm struct {
	InstanceID string            `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string            `form:"schema" json:"schema" binding:"required"`
	Params     map[string]string `form:"params" json:"params"`
	SQLText    string            `form:"sqltext" json:"sqltext" binding:"required"`
}

type ExecuteClickHouseQueryForm struct {
	InstanceID string                 `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string                 `form:"schema" json:"schema" binding:"required"`
	Params     map[string]interface{} `form:"params" json:"params"`
	SQLText    string                 `form:"sqltext" json:"sqltext" binding:"required"`
}
