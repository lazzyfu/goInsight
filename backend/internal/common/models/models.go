package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 工单环境
type InsightDBEnvironments struct {
	*Model
	Name string `gorm:"type:varchar(32);not null;default:'';comment:环境名;uniqueIndex:uniq_name" json:"name"`
}

func (InsightDBEnvironments) TableName() string {
	return "insight_db_environments"
}

// 工单DB配置
type InsightDBConfig struct {
	*Model
	InstanceID       uuid.UUID      `gorm:"type:char(36);uniqueIndex:uniq_instance_id" json:"instance_id"`
	Hostname         string         `gorm:"type:varchar(128);not null;default:'';uniqueIndex:uniq_hostname;comment:主机名" json:"hostname"`
	Port             int            `gorm:"type:int;not null;default 3306;uniqueIndex:uniq_hostname;comment:端口" json:"port"`
	User             string         `gorm:"type:varchar(32);not null;default:'';comment:用户名" json:"user"`
	Password         string         `gorm:"type:varchar(128);not null;default:'';comment:密码" json:"password"`
	UseType          EnumType       `gorm:"type:ENUM('查询', '工单');default:工单;uniqueIndex:uniq_hostname;comment:用途" json:"use_type"`
	DbType           EnumType       `gorm:"type:ENUM('MySQL', 'TiDB', 'ClickHouse');default:MySQL;comment:数据库类型" json:"db_type"`
	Environment      int            `gorm:"type:int;null;default:null;comment:环境;index" json:"environment"`
	OrganizationKey  string         `gorm:"type:varchar(256);not null;index:organization_key;comment:搜索路径" json:"organization_key"`
	OrganizationPath datatypes.JSON `gorm:"type:json;null;default:null;comment:绝对路径" json:"organization_path"`
	Remark           string         `gorm:"type:varchar(256);not null;default:'';comment:备注" json:"remark"`
}

func (InsightDBConfig) TableName() string {
	return "insight_db_config"
}

func (u *InsightDBConfig) BeforeCreate(tx *gorm.DB) (err error) {
	u.InstanceID, _ = uuid.NewUUID()
	return
}

// 自动采集和存储InsightDBConfig配置实例的库
type InsightDBSchemas struct {
	*Model
	InstanceID uuid.UUID `gorm:"type:char(36);comment:关联insight_db_config的instance_id;uniqueIndex:uniq_schema" json:"instance_id"`
	Schema     string    `gorm:"type:varchar(128);not null;default:'';comment:库名;uniqueIndex:uniq_schema" json:"schema"`
	IsDeleted  bool      `gorm:"type:boolean;not null;default:false;comment:是否删除;uniqueIndex:uniq_schema" json:"is_deleted"`
}

func (InsightDBSchemas) TableName() string {
	return "insight_db_schemas"
}
