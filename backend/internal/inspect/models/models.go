package models

import (
	"github.com/google/uuid"
	"github.com/lazzyfu/goinsight/internal/common/models"

	"gorm.io/datatypes"
)

// 工单记录
type InsightInspectParams struct {
	*models.Model
	Params datatypes.JSON `gorm:"type:json;null;default:null;comment:语法审核参数" json:"params"`
	Remark string         `gorm:"type:varchar(256);null;default:null;uiqueIndex:uniq_remark;comment:备注" json:"remark"`
}

func (InsightInspectParams) TableName() string {
	return "insight_inspect_params"
}

type InsightGlobalInspectParams struct {
	*models.Model
	Title string          `gorm:"type:varchar(64);not null;uniqueIndex:uniq_title;comment:审核参数描述" json:"title"`
	Key   string          `gorm:"type:varchar(64);not null;uniqueIndex:uniq_key;comment:审核参数名称" json:"key"`
	Value string          `gorm:"type:varchar(256);not null;comment:审核参数值" json:"value"`
	Type  models.EnumType `gorm:"type:enum('string','number','boolean','json');not null;default:'string';comment:参数类型" json:"type"`
}

func (InsightGlobalInspectParams) TableName() string {
	return "insight_global_inspect_params"
}

type InsightInstanceInspectParams struct {
	*models.Model
	InstanceID uuid.UUID       `gorm:"type:char(36);uniqueIndex:uniq_instance_key;comment:关联insight_db_config的instance_id" json:"instance_id"`
	Title      string          `gorm:"type:varchar(64);not null;comment:审核参数描述" json:"title"`
	Key        string          `gorm:"type:varchar(64);not null;uniqueIndex:uniq_instance_key;comment:审核参数名称" json:"key"`
	Value      string          `gorm:"type:varchar(256);not null;comment:审核参数值" json:"value"`
	Type       models.EnumType `gorm:"type:enum('string','number','boolean','json');not null;default:'string';comment:参数类型" json:"type"`
}

func (InsightInstanceInspectParams) TableName() string {
	return "insight_instance_inspect_params"
}
