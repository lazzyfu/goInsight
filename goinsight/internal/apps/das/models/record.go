package models

import (
	"goInsight/internal/apps/common/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type InsightDASRecords struct {
	*models.Model
	Username       string         `gorm:"type:varchar(128);not null;comment:执行的用户;index:idx_username" json:"username"`
	InstanceID     uuid.UUID      `gorm:"type:char(36);comment:关联insight_db_config的instance_id;index:idx_instance_id" json:"instance_id"`
	Schema         string         `gorm:"type:varchar(512);not null;default:'';comment:执行的库名;index:idx_schema" json:"schema"`
	Tables         string         `gorm:"type:varchar(4096);not null;default:'';comment:提取的库表名" json:"tables"`
	Sqltext        string         `gorm:"type:text;null;default:null;comment:用户输入的原始SQL" json:"sqltext"`
	QueryID        string         `gorm:"type:char(36);null;default null;comment:查询ID" json:"query_id"`
	Duration       int            `gorm:"type:int;not null;default:0;comment:耗时(单位ms)" json:"duration"`
	RewriteSqltext string         `gorm:"type:text;null;default:null;comment:rewrite后的SQL" json:"rewrite_sqltext"`
	Params         datatypes.JSON `gorm:"type:json;null;default:null;comment:参数" json:"params"` // 以json格式存储
	IsFinish       bool           `gorm:"type:boolean;null;default:False;comment:是否执行完成,0未完成,1已完成" json:"is_finish"`
	IsKill         bool           `gorm:"type:boolean;null;default:False;comment:TiDB连接是否被Kill,0否,1是" json:"is_kill"`
	ErrorMsg       string         `gorm:"type:varchar(8092);not null;default '';comment:错误信息" json:"error_msg"`
	ReturnRows     int            `gorm:"type:int;not null;default:0;comment:返回行数" json:"return_rows"`
	RequestID      string         `gorm:"type:string;not null;default:'';comment:请求ID;index:idx_request_id" json:"request_id"`
}

func (InsightDASRecords) TableName() string {
	return "insight_das_records"
}
