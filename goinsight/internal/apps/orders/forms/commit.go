package forms

import (
	"goInsight/internal/apps/common/models"
	"goInsight/internal/pkg/pagination"
)

type GetInstancesForm struct {
	PaginationQ pagination.Pagination
	ID          uint64 `form:"id"  json:"id" binding:"required"`
	DbType      string `form:"db_type" json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
}

type GetSchemasForm struct {
	PaginationQ pagination.Pagination
	InstanceID  string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
}

type GetUsersForm struct {
	PaginationQ pagination.Pagination
}

type CreateOrderForm struct {
	Title            string          `form:"title" json:"title" binding:"required,min=5,max=96"`
	Remark           string          `form:"remark" json:"remark" binding:"max=1024"`
	IsRestrictAccess *bool           `form:"is_restrict_access" json:"is_restrict_access" validate:"boolean" binding:"required"`
	DBType           models.EnumType `form:"db_type" json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	SQLType          models.EnumType `form:"sql_type" json:"sql_type" binding:"required,oneof=DML DDL EXPORT"`
	Environment      int             `form:"environment" json:"environment" binding:"required"`
	InstanceID       string          `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema           string          `form:"schema" json:"schema" binding:"max=1024"`
	Approver         []string        `form:"approver" json:"approver" binding:"required"`
	Executor         []string        `form:"executor" json:"executor" binding:"required"`
	Reviewer         []string        `form:"reviewer" json:"reviewer" binding:"required"`
	CC               []string        `form:"cc" json:"cc"`
	Content          string          `form:"content" json:"content" binding:"required"`
}
