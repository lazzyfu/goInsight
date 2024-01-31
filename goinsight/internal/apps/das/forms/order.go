/*
@Time    :   2023/08/02 17:49:29
@Author  :   zongfei.fu
@Desc    :
*/

package forms

type GetOrderSchemasForm struct {
	Environment string `form:"environment" json:"environment" binding:"required"`
}

type GetOrderTablesForm struct {
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string `form:"schema" json:"schema" binding:"required"`
}
