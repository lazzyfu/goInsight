package forms

type GetDbDictForm struct {
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string `form:"schema" json:"schema" binding:"required"`
}
