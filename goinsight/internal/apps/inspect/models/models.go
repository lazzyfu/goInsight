/*
@Time    :   2023/09/21 19:49:45
@Author  :   xff
*/

package models

import (
	"goInsight/internal/apps/common/models"

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
