/*
@Time    :   2023/09/21 19:49:45
@Author  :   lazzyfu
*/

package models

import (
	"goInsight/internal/app/common/models"

	"gorm.io/datatypes"
)

// 工单记录
type InsightInspectParams struct {
	*models.Model
	Key    string         `gorm:"type:varchar(64);comment:key;uniqueIndex:uniq_key" json:"key"`
	Value  datatypes.JSON `gorm:"type:json;null;default:null;comment:值" json:"value"`
	Remark string         `gorm:"type:varchar(128);not null;default:'';comment:备注" json:"remark"`
}

func (InsightInspectParams) TableName() string {
	return "insight_inspect_params"
}
