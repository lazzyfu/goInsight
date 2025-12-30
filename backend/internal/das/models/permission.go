package models

import (
	"github.com/lazzyfu/goinsight/internal/common/models"

	"github.com/google/uuid"
)

// 用户的库权限
type InsightDasSchemaPerms struct {
	*models.Model
	Username   string    `gorm:"type:varchar(128);not null;comment:用户;uniqueIndex:uniq_schema" json:"username"`
	Schema     string    `gorm:"type:varchar(128);not null;default:'';comment:库名;uniqueIndex:uniq_schema" json:"schema"`
	InstanceID uuid.UUID `gorm:"type:char(36);comment:关联das_config的instance_id;uniqueIndex:uniq_schema;index:idx_instance_id" json:"instance_id"`
}

func (InsightDasSchemaPerms) TableName() string {
	return "insight_das_schema_perms"
}

// 用户的表权限
type InsightDasTablePerms struct {
	*models.Model
	Username   string          `gorm:"type:varchar(128);not null;comment:用户;uniqueIndex:uniq_table" json:"username"`
	Schema     string          `gorm:"type:varchar(128);not null;default:'';comment:库名;uniqueIndex:uniq_table" json:"schema"`
	Table      string          `gorm:"type:varchar(128);not null;default:'';comment:表名;uniqueIndex:uniq_table" json:"table"`
	InstanceID uuid.UUID       `gorm:"type:char(36);comment:关联insight_instances的instance_id;uniqueIndex:uniq_table;index:idx_instance_id" json:"instance_id"`
	Rule       models.EnumType `gorm:"type:ENUM('allow', 'deny');default:allow;comment:规则" json:"rule"`
}

func (InsightDasTablePerms) TableName() string {
	return "insight_das_table_perms"
}

// 允许用户执行的操作
type InsightDASOperations struct {
	*models.Model
	Name     string `gorm:"type:varchar(128);not null;comment:语句类型;uniqueIndex:uniq_name" json:"name"`
	IsEnable bool   `gorm:"type:boolean;null;default:False;comment:是否启用,0未启用,1启用" json:"is_finish"`
	Remark   string `gorm:"type:varchar(1024);not null;default:'';comment:备注" json:"remark"`
}

func (InsightDASOperations) TableName() string {
	return "insight_das_operations"
}
