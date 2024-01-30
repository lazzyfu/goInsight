package forms

import (
	"goInsight/internal/app/common/models"
	"goInsight/internal/pkg/pagination"
)

type AdminGetDBConfigForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
	Environment string `form:"environment"`
	DbType      string `form:"db_type"`
}

type AdminCreateDBConfigForm struct {
	Hostname        string          `form:"hostname"  json:"hostname" binding:"required,min=2,max=128"`
	Port            int             `form:"port"  json:"port" binding:"required"`
	InspectParams   string          `form:"inspect_params" json:"inspect_params"`
	UseType         models.EnumType `form:"use_type"  json:"use_type" binding:"required,oneof=查询 工单"`
	DbType          models.EnumType `form:"db_type"  json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	Environment     int             `form:"environment"  json:"environment" binding:"required"`
	OrganizationKey []string        `form:"organization_key"  json:"organization_key" binding:"required"`
	Remark          string          `form:"remark"  json:"remark" binding:"required,min=2,max=256"`
}

type AdminUpdateDBConfigForm struct {
	Hostname        string          `form:"hostname"  json:"hostname" binding:"required,min=2,max=128"`
	Port            int             `form:"port"  json:"port" binding:"required"`
	InspectParams   string          `form:"inspect_params" json:"inspect_params"`
	UseType         models.EnumType `form:"use_type"  json:"use_type" binding:"required,oneof=查询 工单"`
	DbType          models.EnumType `form:"db_type"  json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	Environment     int             `form:"environment"  json:"environment" binding:"required"`
	OrganizationKey []string        `form:"organization_key"  json:"organization_key" binding:"required"`
	Remark          string          `form:"remark"  json:"remark" binding:"required,min=2,max=256"`
}
