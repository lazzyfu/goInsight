package models

import (
	"github.com/lazzyfu/goinsight/internal/common/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 工单记录
type InsightOrderRecords struct {
	*models.Model
	Title            string          `gorm:"type:varchar(256);not null;default:'';comment:工单标题;index:idx_title" json:"title"`
	OrderID          uuid.UUID       `gorm:"type:char(36);comment:工单ID;uniqueIndex:uniq_order_id" json:"order_id"`
	HookOrderID      uuid.UUID       `gorm:"type:char(36);comment:HOOK源工单ID;index:idx_hook_order_id" json:"hook_order_id"`
	Remark           string          `gorm:"type:varchar(2048);not null;default:'';comment:工单备注" json:"remark"`
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
	FlowID           uint            `json:"flow_id"` // 关联 ApprovalFlow
	CurrentStage     int             `json:"current_stage"`
	WorkflowID       string          `json:"workflow_id"` // goflow 工作流ID
	Status           string          `json:"status"`      // pending, in_approval, approved, rejected, executed
}

func (InsightOrderRecords) TableName() string {
	return "insight_order_records"
}

// 审批流
type ApprovalFlow struct {
	gorm.Model
	ApprovalID uuid.UUID      `gorm:"type:char(36);comment:审批流ID;uniqueIndex:uniq_approval_id" json:"approval_id"`
	Name       string         `json:"name"`
	Definition datatypes.JSON `json:"definition"` // [{"stage":1, "approvers":["zhangsan","lisi"], "type":"AND"}]
}

func (ApprovalFlow) TableName() string {
	return "insight_approval_flow"
}

// 审批流和用户映射表，每个用户只能在一个审批流里面
type ApprovalMaps struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(32);not null;uniqueIndex:uniq_username;comment:用户名" json:"username"`
	ApprovalID uuid.UUID `gorm:"type:char(36);comment:审批流ID;index:idx_approval_id" json:"approval_id"`
}

func (ApprovalMaps) TableName() string {
	return "insight_approval_maps"
}

// 审批记录
type ApprovalRecords struct {
	gorm.Model
	OrderID  uuid.UUID       `json:"order_id" gorm:"char(36)"`
	Stage    int             `json:"stage"` // 审批阶段，1/2/3...
	Approver string          `gorm:"type:varchar(32);not null;comment:审批人" json:"approver"`
	Status   models.EnumType `gorm:"type:ENUM('PENDING', 'APPROVED', 'REJECTED');default:'PENDING';comment:审批状态" json:"status"`
}

func (ApprovalRecords) TableName() string {
	return "insight_approval_records"
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
