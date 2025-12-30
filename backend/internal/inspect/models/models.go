package models

import (
	"github.com/google/uuid"
	"github.com/lazzyfu/goinsight/internal/common/models"
)

// 全局审核参数
type InsightInspectGlobalParams struct {
	*models.Model
	Title string          `gorm:"type:varchar(256);not null;uniqueIndex:uniq_title;comment:审核参数描述" json:"title"`
	Key   string          `gorm:"type:varchar(64);not null;uniqueIndex:uniq_key;comment:审核参数名称" json:"key"`
	Value string          `gorm:"type:varchar(256);not null;comment:审核参数值" json:"value"`
	Type  models.EnumType `gorm:"type:enum('string','number','boolean');not null;default:'string';comment:参数类型" json:"type"`
}

func (InsightInspectGlobalParams) TableName() string {
	return "insight_inspect_global_params"
}

// 实例审核参数，优先级>全局审核参数
type InsightInspectInstanceParams struct {
	*models.Model
	InstanceID uuid.UUID       `gorm:"type:char(36);uniqueIndex:uniq_instance_key;comment:关联insight_instances的instance_id" json:"instance_id"`
	Title      string          `gorm:"type:varchar(256);not null;comment:审核参数描述" json:"title"`
	Key        string          `gorm:"type:varchar(64);not null;uniqueIndex:uniq_instance_key;comment:审核参数名称" json:"key"`
	Value      string          `gorm:"type:varchar(256);not null;comment:审核参数值" json:"value"`
	Type       models.EnumType `gorm:"type:enum('string','number','boolean');not null;default:'string';comment:参数类型" json:"type"`
}

func (InsightInspectInstanceParams) TableName() string {
	return "insight_inspect_instance_params"
}
