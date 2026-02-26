package models

import (
	"github.com/lazzyfu/goinsight/internal/common/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// 工单记录
/**
# 工单生命周期
"PENDING":   "待审批",
"APPROVED":  "已批准",
"REJECTED":  "已驳回",
"CLAIMED":   "已认领",
"EXECUTING": "执行中",
"COMPLETED": "已完成",
"FAILED":    "已失败",
"REVIEWED":  "已复核",
"REVOKED":   "已撤销",
*/
type InsightOrderRecords struct {
	*models.Model
	Title            string          `gorm:"type:varchar(256);not null;default:'';comment:工单标题;index:idx_title" json:"title"`
	OrderID          uuid.UUID       `gorm:"type:char(36);comment:工单ID;uniqueIndex:uniq_order_id" json:"order_id"`
	HookOrderID      uuid.UUID       `gorm:"type:char(36);comment:HOOK源工单ID;index:idx_hook_order_id" json:"hook_order_id"`
	Remark           string          `gorm:"type:varchar(2048);not null;default:'';comment:工单备注" json:"remark"`
	IsRestrictAccess bool            `gorm:"type:tinyint(1);not null;default:0;comment:是否限制访问" json:"is_restrict_access"`
	DBType           models.EnumType `gorm:"type:ENUM('MySQL', 'TiDB', 'ClickHouse');default:'MySQL';comment:DB类型" json:"db_type"`
	SQLType          models.EnumType `gorm:"type:ENUM('DML', 'DDL', 'EXPORT');default:'DML';comment:SQL类型" json:"sql_type"`
	Environment      string          `gorm:"type:varchar(32);not null;default:'';comment:环境名;index:idx_environment" json:"environment"`
	Applicant        string          `gorm:"type:varchar(32);not null;default:'';comment:申请人;index" json:"applicant"`
	Organization     string          `gorm:"type:varchar(256);not null;default:'';index;comment:组织" json:"organization"`
	Claimer          string          `gorm:"type:varchar(32);not null;default:'';comment:认领人;index" json:"claimer"`
	ClaimUsers       datatypes.JSON  `gorm:"type:json;null;default:null;comment:可认领人员列表" json:"claim_users"`
	Executor         string          `gorm:"type:varchar(32);not null;default:'';comment:工单执行人;index" json:"executor"`
	Approver         datatypes.JSON  `gorm:"type:json;null;default:null;comment:工单审核人" json:"approver"`
	Reviewer         datatypes.JSON  `gorm:"type:json;null;default:null;comment:工单复核人" json:"reviewer"`
	CC               datatypes.JSON  `gorm:"type:json;null;default:null;comment:工单抄送人" json:"cc"`
	InstanceID       uuid.UUID       `gorm:"type:char(36);comment:关联insight_instances的instance_id;index" json:"instance_id"`
	Schema           string          `gorm:"type:varchar(128);not null;default:'';comment:库名" json:"schema"`
	Stage            int             `gorm:"type:tinyint(1);not null;default:1;comment:审批阶段" json:"stage"`
	Progress         models.EnumType `gorm:"type:ENUM('PENDING','APPROVED','REJECTED','CLAIMED','EXECUTING','COMPLETED', 'FAILED','REVIEWED','REVOKED');default:'PENDING';comment:工单进度" json:"progress"`
	FixVersion       string          `gorm:"type:varchar(128);not null;default:'';comment:上线版本;index" json:"fix_version"`
	Content          string          `gorm:"type:text;null;comment:工单内容" json:"content"`
	ExportFileFormat models.EnumType `gorm:"type:ENUM('XLSX', 'CSV');default:'XLSX';comment:导出文件格式" json:"export_file_format"`
}

func (InsightOrderRecords) TableName() string {
	return "insight_order_records"
}

// 审批流
type InsightApprovalFlows struct {
	*models.Model
	ApprovalID uuid.UUID      `gorm:"type:char(36);comment:审批流ID;uniqueIndex:uniq_approval_id" json:"approval_id"`
	Name       string         `gorm:"type:varchar(64);not null;default:'';comment:审批流名称" json:"name"`
	ClaimUsers datatypes.JSON `gorm:"type:json;null;default:null;comment:可认领人员列表" json:"claim_users"`
	Definition datatypes.JSON `json:"definition"` // [{"stage":1, "approvers":["zhangsan","lisi"], "type":"AND", "stage_name": '部门审批'}]
}

func (InsightApprovalFlows) TableName() string {
	return "insight_approval_flow"
}

// 审批流和用户映射表，每个用户只能在一个审批流里面
type InsightApprovalFlowUsers struct {
	*models.Model
	Username   string    `gorm:"type:varchar(32);not null;uniqueIndex:uniq_username;comment:用户名" json:"username"`
	ApprovalID uuid.UUID `gorm:"type:char(36);comment:审批流ID;index:idx_approval_id" json:"approval_id"`
}

func (InsightApprovalFlowUsers) TableName() string {
	return "insight_approval_flow_users"
}

// 审批记录
type InsightApprovalRecords struct {
	*models.Model
	OrderID        uuid.UUID        `gorm:"type:char(36);comment:工单ID;index:index_order_id" json:"order_id"`
	Stage          int              `gorm:"type:tinyint(1);not null;default:1;comment:审批阶段" json:"stage"`
	StageName      string           `gorm:"type:varchar(64);null;default:null;comment:审批阶段名称" json:"stage_name"`
	Approver       string           `gorm:"type:varchar(32);not null;comment:审批人" json:"approver"`
	ApprovalType   models.EnumType  `gorm:"type:ENUM('AND', 'OR');default:'AND';comment:审批类型" json:"approval_type"`
	ApprovalStatus models.EnumType  `gorm:"type:ENUM('PENDING', 'APPROVED', 'REJECTED');default:'PENDING';comment:审批状态" json:"approval_status"`
	ApprovalAt     models.LocalTime `gorm:"comment:审批时间" json:"approval_at"`
	Msg            string           `gorm:"type:varchar(512);not null;default:'';comment:审批意见" json:"msg"`
}

func (InsightApprovalRecords) TableName() string {
	return "insight_approval_records"
}

// 工单操作日志表
type InsightOrderLogs struct {
	*models.Model
	Username string    `gorm:"type:varchar(32);not null;index:idx_username;comment:操作用户" json:"username"`
	OrderID  uuid.UUID `gorm:"type:char(36);comment:工单ID;index:idx_order_id" json:"order_id"`
	Msg      string    `gorm:"type:varchar(1024);null;;comment:操作信息" json:"msg"`
}

func (InsightOrderLogs) TableName() string {
	return "insight_order_logs"
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
	Progress models.EnumType `gorm:"type:ENUM('PENDING', 'EXECUTING', 'COMPLETED', 'FAILED');default:'PENDING';comment:进度" json:"progress"`
	Result   datatypes.JSON  `gorm:"type:json;null;default:null;comment:执行结果" json:"result"`
}

func (InsightOrderTasks) TableName() string {
	return "insight_order_tasks"
}

// 消息推送记录，暂时没有用到
type InsightOrderMessages struct {
	*models.Model
	OrderID  uuid.UUID      `gorm:"type:char(36);comment:关联insight_order_records的order_id;index" json:"order_id"`
	Receiver datatypes.JSON `gorm:"type:json;null;default:null;comment:接收消息的用户" json:"receiver"`
	Response string         `gorm:"type:text;null;comment:第三方返回的响应" json:"response"`
}

func (InsightOrderMessages) TableName() string {
	return "insight_order_messages"
}
