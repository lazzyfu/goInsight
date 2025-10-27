package forms

import "github.com/lazzyfu/goinsight/pkg/pagination"

type UserGrantsForm struct {
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string `form:"schema" json:"schema" binding:"required"`
}

type AdminSchemasGrantForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
	Environment string `form:"environment"`
}

type AdminGetInstancesListForm struct {
	PaginationQ pagination.Pagination
	ID          uint64 `form:"id"  json:"id" binding:"required"`
	DbType      string `form:"db_type" json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
}

type AdminGetSchemasListForm struct {
	PaginationQ pagination.Pagination
	InstanceID  string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
}

type AdminGetTablesListForm struct {
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string `form:"schema" json:"schema" binding:"required"`
}

type AdminCreateSchemasGrantForm struct {
	Username   string   `form:"username" json:"username" binding:"required"`
	InstanceID string   `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string   `form:"schema" json:"schema" binding:"required"`
	Tables     []string `form:"tables" json:"tables"`
}

type AdminGetTablesGrantForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
	Username    string `form:"username" json:"username" binding:"required"`
	InstanceID  string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema      string `form:"schema" json:"schema" binding:"required"`
}

type AdminCreateTablesGrantForm struct {
	Username   string   `form:"username" json:"username" binding:"required"`
	InstanceID string   `form:"instance_id" json:"instance_id" binding:"required,uuid"`
	Schema     string   `form:"schema" json:"schema" binding:"required"`
	Tables     []string `form:"tables" json:"tables"`
	Rule       string   `form:"rule" json:"rule" binding:"required"`
}
