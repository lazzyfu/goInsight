package forms

import (
	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/common/models"
)

type AdminGetInstancesForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
	Environment string `form:"environment"`
	DbType      string `form:"db_type"`
}

type AdminCreateInstancesForm struct {
	Hostname        string                 `form:"hostname"  json:"hostname" binding:"required,min=2,max=128"`
	Port            int                    `form:"port"  json:"port" binding:"required"`
	User            string                 `form:"user"  json:"user" binding:"required,min=3,max=32"`
	Password        string                 `form:"password"  json:"password" binding:"required,min=8,max=64"`
	InspectParams   map[string]interface{} `form:"inspect_params" json:"inspect_params"`
	UseType         models.EnumType        `form:"use_type"  json:"use_type" binding:"required,oneof=查询 工单"`
	DbType          models.EnumType        `form:"db_type"  json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	Environment     int                    `form:"environment"  json:"environment" binding:"required"`
	OrganizationKey []string               `form:"organization_key"  json:"organization_key" binding:"required"`
	Remark          string                 `form:"remark"  json:"remark" binding:"required,min=2,max=256"`
}

type AdminUpdateInstancesForm struct {
	Hostname        string                 `form:"hostname"  json:"hostname" binding:"required,min=2,max=128"`
	Port            int                    `form:"port"  json:"port" binding:"required"`
	User            string                 `form:"user"  json:"user" binding:"required,min=3,max=32"`
	Password        string                 `form:"password"  json:"password" binding:"required,min=8,max=64"`
	InspectParams   map[string]interface{} `form:"inspect_params" json:"inspect_params"`
	UseType         models.EnumType        `form:"use_type"  json:"use_type" binding:"required,oneof=查询 工单"`
	DbType          models.EnumType        `form:"db_type"  json:"db_type" binding:"required,oneof=MySQL TiDB ClickHouse"`
	Environment     int                    `form:"environment"  json:"environment" binding:"required"`
	OrganizationKey []string               `form:"organization_key"  json:"organization_key" binding:"required"`
	Remark          string                 `form:"remark"  json:"remark" binding:"required,min=2,max=256"`
}

type AdminGetInstanceInspectParamsForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
	InstanceID  string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
}

type AdminCreateInstanceInspectParamsForm struct {
	Title      string `form:"title" json:"title" binding:"required,min=3,max=256"`
	Type       string `form:"type"  json:"type"  binding:"required,oneof=string number boolean"`
	Key        string `form:"key" json:"key" binding:"required,min=3,max=256"`
	Value      string `form:"value" json:"value" binding:"required,min=1,max=256"`
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
}

type AdminUpdateInstanceInspectParamsForm struct {
	Title      string `form:"title" json:"title" binding:"required,min=3,max=256"`
	Type       string `form:"type"  json:"type"  binding:"required,oneof=string number boolean"`
	Key        string `form:"key" json:"key" binding:"required,min=3,max=256"`
	Value      string `form:"value" json:"value" binding:"required,min=1,max=256"`
	InstanceID string `form:"instance_id" json:"instance_id" binding:"required,uuid"`
}
