/*
@Time    :   2023/09/21 19:49:45
@Author  :   xff
*/

package models

import (
	"goInsight/internal/common/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// 工单记录
type InsightOrderRecords struct {
	*models.Model
	Title            string          `gorm:"type:varchar(128);not null;default:'';comment:工单标题;index:idx_title" json:"title"`
	OrderID          uuid.UUID       `gorm:"type:char(36);comment:工单ID;uniqueIndex:uniq_order_id" json:"order_id"`
	HookOrderID      uuid.UUID       `gorm:"type:char(36);comment:HOOK源工单ID;index:idx_hook_order_id" json:"hook_order_id"`
	Remark           string          `gorm:"type:varchar(1024);not null;default:'';comment:工单备注" json:"remark"`
	IsRestrictAccess bool            `gorm:"type:tinyint(1);not null;default:0;comment:是否限制访问" json:"is_restrict_access"`
	DBType           models.EnumType `gorm:"type:ENUM('MySQL', 'TiDB', 'ClickHouse');default:'MySQL';comment:DB类型" json:"db_type"`
	SQLType          models.EnumType `gorm:"type:ENUM('DML', 'DDL', 'EXPORT');default:'DML';comment:SQL类型" json:"sql_type"`
	Environment      int             `gorm:"type:int;null;default:null;comment:环境;index" json:"environment"`
	Applicant        string          `gorm:"type:varchar(32);not null;default:'';comment:申请人;index" json:"applicant"`
	Organization     string          `gorm:"type:varchar(256);not null;default:'';index;comment:组织" json:"organization"`
	Approver         datatypes.JSON  `gorm:"type:json;null;default:null;comment:工单审核人" json:"approver"`
	Executor         datatypes.JSON  `gorm:"type:json;null;default:null;comment:工单执行人" json:"executor"`
	Reviewer         datatypes.JSON  `gorm:"type:json;null;default:null;comment:工单复核人" json:"reviewer"`
	CC               datatypes.JSON  `gorm:"type:json;null;default:null;comment:工单抄送人" json:"cc"`
	InstanceID       uuid.UUID       `gorm:"type:char(36);comment:关联insight_db_config的instance_id;index" json:"instance_id"`
	Schema           string          `gorm:"type:varchar(128);not null;default:'';comment:库名" json:"schema"`
	Progress         models.EnumType `gorm:"type:ENUM('待审核', '已驳回', '已批准', '执行中', '已关闭', '已完成', '已复核');default:'待审核';comment:工单进度" json:"progress"`
	FixVersion       string          `gorm:"type:varchar(128);not null;default:'';comment:上线版本;index" json:"fix_version"`
	Content          string          `gorm:"type:text;null;comment:工单内容" json:"content"`
	ExportFileFormat models.EnumType `gorm:"type:ENUM('XLSX', 'CSV');default:'XLSX';comment:导出文件格式" json:"export_file_format"`
}

func (InsightOrderRecords) TableName() string {
	return "insight_order_records"
}

// 工单操作日志表
type InsightOrderOpLogs struct {
	*models.Model
	Username string    `gorm:"type:varchar(32);not null;index:idx_username;comment:操作用户" json:"username"`
	OrderID  uuid.UUID `gorm:"type:char(36);comment:工单ID;index:idx_order_id" json:"order_id"`
	Msg      string    `gorm:"type:varchar(1024);null;;comment:操作信息" json:"msg"`
}

func (InsightOrderOpLogs) TableName() string {
	return "insight_order_oplogs"
}

// 工单记录生成的执行任务
type InsightOrderTasks struct {
	*models.Model
	OrderID  uuid.UUID       `gorm:"type:char(36);comment:关联insight_order_records的order_id;index" json:"order_id"`
	TaskID   uuid.UUID       `gorm:"type:char(36);comment:任务ID;index" json:"task_id"`
	DBType   models.EnumType `gorm:"type:ENUM('MySQL', 'TiDB', 'ClickHouse');default:'MySQL';comment:DB类型" json:"db_type"`
	SQLType  models.EnumType `gorm:"type:ENUM('DML', 'DDL', 'EXPORT');default:'DML';comment:SQL类型" json:"sql_type"`
	Executor string          `gorm:"type:varchar(128);null;default:null;comment:任务执行人" json:"executor"`
	SQL      string          `gorm:"type:text;null;comment:SQL语句" json:"sql"`
	Progress models.EnumType `gorm:"type:ENUM('未执行', '执行中', '已完成', '已失败', '已暂停');default:'未执行';comment:进度" json:"progress"`
	Result   datatypes.JSON  `gorm:"type:json;null;default:null;comment:执行结果" json:"result"`
}

func (InsightOrderTasks) TableName() string {
	return "insight_order_tasks"
}

// 消息推送记录
type InsightOrderMessages struct {
	*models.Model
	OrderID  uuid.UUID      `gorm:"type:char(36);comment:关联insight_order_records的order_id;index" json:"order_id"`
	Receiver datatypes.JSON `gorm:"type:json;null;default:null;comment:接收消息的用户" json:"receiver"`
	Response string         `gorm:"type:text;null;comment:第三方返回的响应" json:"response"`
}

func (InsightOrderMessages) TableName() string {
	return "insight_order_messages"
}
