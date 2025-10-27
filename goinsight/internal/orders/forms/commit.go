package forms

import (
	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/common/models"
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
	Title            string          `form:"title" json:"title" binding:"required,min=3,max=196"`
	Remark           string          `form:"remark" json:"remark" binding:"max=1024"`
	DBType           models.EnumType `form:"db_type" json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	SQLType          models.EnumType `form:"sql_type" json:"sql_type" binding:"required,oneof=DML DDL EXPORT"`
	Environment      int             `form:"environment" json:"environment" binding:"required"`
	InstanceID       string          `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema           string          `form:"schema" json:"schema" binding:"max=1024"`
	CC               []string        `form:"cc" json:"cc"`
	Content          string          `form:"content" json:"content" binding:"required"`
	ExportFileFormat models.EnumType `form:"export_file_format" json:"export_file_format" binding:"required,oneof=XLSX CSV"`
}
