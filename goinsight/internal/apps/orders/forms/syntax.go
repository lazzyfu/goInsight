package forms

import (
	"goInsight/internal/apps/common/models"
)

type SyntaxInspectForm struct {
	DBType     models.EnumType `form:"db_type" json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	SQLType    models.EnumType `form:"sql_type" json:"sql_type" binding:"required,oneof=DML DDL EXPORT"`
	InstanceID string          `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string          `form:"schema" json:"schema" binding:"max=1024"`
	Content    string          `form:"content" json:"content" binding:"required"`
}
