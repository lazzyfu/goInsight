package forms

import (
	"goInsight/internal/common/models"
)

type HookOrdersForm struct {
	OrderID  string          `form:"order_id" json:"order_id" binding:"required,uuid"`
	Progress string          `form:"progress" json:"progress" binding:"required,oneof=待审核 已批准"`
	DBType   models.EnumType `form:"db_type" json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	Target   []struct {
		Environment int    `form:"environment" json:"environment" binding:"required"`
		InstanceID  string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
		Schema      string `form:"schema" json:"schema" binding:"max=1024"`
	}
}
