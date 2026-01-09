package notifier

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"
	ordersModels "github.com/lazzyfu/goinsight/internal/orders/models"
)

// 消息模版类型
type MessageType string

const (
	// 工单提交
	MsgTypeOrderSubmitted MessageType = "ORDER_SUBMITTED"
	// 待审批
	MsgTypeOrderPendingApproval MessageType = "ORDER_PENDING_APPROVAL"
	// 审批通过
	MsgTypeOrderApproved MessageType = "ORDER_APPROVED"
	// 审批驳回
	MsgTypeOrderRejected MessageType = "ORDER_REJECTED"
	// 工单认领
	MsgTypeOrderClaimed MessageType = "ORDER_CLAIMED"
	// 工单转交
	MsgTypeOrderTransferred MessageType = "ORDER_TRANSFERRED"
	// 工单撤销
	MsgTypeOrderRevoked MessageType = "ORDER_REVOKED"
	// 工单完成
	MsgTypeOrderCompleted MessageType = "ORDER_COMPLETED"
	// 工单失败
	MsgTypeOrderFailed MessageType = "ORDER_FAILED"
	// 工单复核
	MsgTypeOrderReviewed MessageType = "ORDER_REVIEWED"
	// 工单执行完成
	MsgTypeOrderExecutionCompleted MessageType = "ORDER_EXECUTION_COMPLETED"
	// 导出文件信息
	MsgTypeExportFileInfo MessageType = "EXPORT_FILE_INFO"
)

// 消息参数
type MessageParams struct {
	// 工单记录
	Order         *ordersModels.InsightOrderRecords
	Task          *ordersModels.InsightOrderTasks
	Username      string   // 操作用户
	AdditionalMsg string   // 附加消息
	Approvers     []string // 待审批消息中需要 @ 的审批人
}

func buildApproverMentions(approvers []string) string {
	if len(approvers) == 0 {
		return ""
	}
	// 去重 + 过滤空值，同时保持尽量稳定的顺序
	seen := make(map[string]struct{}, len(approvers))
	mentions := make([]string, 0, len(approvers))
	for _, a := range approvers {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if _, ok := seen[a]; ok {
			continue
		}
		seen[a] = struct{}{}
		mentions = append(mentions, "@"+a)
	}
	if len(mentions) == 0 {
		return ""
	}
	return strings.Join(mentions, " ")
}

// BuildMessage 根据模版类型和参数构建消息
func BuildMessage(msgType MessageType, params MessageParams) string {
	order := params.Order
	task := params.Task

	switch msgType {
	case MsgTypeOrderSubmitted:
		return fmt.Sprintf(
			"您好，用户%s提交了工单\n"+
				">工单标题：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s",
			params.Username, order.Title,
			order.Environment, order.SQLType, order.DBType, order.Schema,
		)

	case MsgTypeOrderPendingApproval:
		mentions := buildApproverMentions(params.Approvers)
		if mentions != "" {
			mentions = "\n" + mentions
		}
		return fmt.Sprintf(
			"您好，您有新的工单需要审批 %s\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s",
			mentions, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema,
		)

	case MsgTypeOrderApproved:
		return fmt.Sprintf(
			"您好，用户%s通过了工单\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderRejected:
		return fmt.Sprintf(
			"您好，用户%s驳回了工单\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderClaimed:
		return fmt.Sprintf(
			"您好，用户%s认领了工单\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderTransferred:
		return fmt.Sprintf(
			"您好，用户%s将工单转交给了%s\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Claimer, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderRevoked:
		return fmt.Sprintf(
			"您好，用户%s撤销了工单\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderCompleted:
		return fmt.Sprintf(
			"您好，用户%s更新工单状态为：已完成\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderFailed:
		return fmt.Sprintf(
			"您好，用户%s更新工单状态为：已失败\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderReviewed:
		return fmt.Sprintf(
			"您好，用户%s更新工单状态为：已复核\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s\n"+
				">附加消息：%s",
			params.Username, order.Title, order.Applicant, order.Environment, order.SQLType, order.DBType, order.Schema, params.AdditionalMsg,
		)

	case MsgTypeOrderExecutionCompleted:
		return fmt.Sprintf(
			"您好，工单已经执行完成，请知悉\n"+
				">工单标题：%s\n"+
				">申请人：%s\n"+
				">环境：%s\n"+
				">执行人：%s\n"+
				">工单类型：%s\n"+
				">数据库类型：%s\n"+
				">库名：%s",
			order.Title, order.Applicant, order.Environment, order.Executor, order.SQLType, order.DBType, order.Schema,
		)

	case MsgTypeExportFileInfo:
		// 获取下载文件信息
		type exportFile struct {
			FileName      string `json:"file_name"`
			FileSize      int64  `json:"file_size"`
			FilePath      string `json:"file_path"`
			ContentType   string `json:"content_type"`
			EncryptionKey string `json:"encryption_key"`
			ExportRows    int64  `json:"export_rows"`
		}
		var file exportFile
		err := json.Unmarshal([]byte(task.Result), &file)
		if err != nil {
			return "解析导出文件信息异常"
		}
		downloadURL := fmt.Sprintf("%s/orders/tasks/exports/%s", global.App.Config.Notify.NoticeURL, file.FileName)
		return fmt.Sprintf(
			"您好，导出文件信息如下，请查收\n"+
				">工单标题：%s\n"+
				">任务ID：%s\n"+
				">文件名：%s\n"+
				">文件大小：%d字节\n"+
				">数据行数：%d\n"+
				">文件解密密码：%s\n"+
				">文件格式：%s\n"+
				">文件下载路径：%s",
			order.Title, task.TaskID,
			file.FileName,
			file.FileSize,
			file.ExportRows,
			file.EncryptionKey,
			file.ContentType,
			downloadURL,
		)

	default:
		return "未知消息类型"
	}
}

// SendOrderMessage 发送工单消息的便捷方法
func SendOrderMessage(receivers []string, msgType MessageType, params MessageParams) {
	msg := BuildMessage(msgType, params)
	go func(title, orderID string, rec []string, msg string) {
		defer func() {
			if r := recover(); r != nil {
				// 这里可以记录日志
			}
		}()
		SendMessage(title, orderID, rec, msg)
	}(params.Order.Title, params.Order.OrderID.String(), receivers, msg)
}
