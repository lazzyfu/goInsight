package forms

type GetTablesForm struct {
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string `form:"schema" json:"schema" binding:"required"`
}

type GetTableInfoForm struct {
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string `form:"schema" json:"schema" binding:"required"`
	Table      string `form:"table" json:"table" binding:"required"`
	Type       string `form:"type" json:"type" binding:"required,oneof=structure base"`
}
