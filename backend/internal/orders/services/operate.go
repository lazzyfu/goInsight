package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lazzyfu/goinsight/internal/global"
	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/models"
	"github.com/lazzyfu/goinsight/pkg/notifier"
	"github.com/lazzyfu/goinsight/pkg/utils"
	"gorm.io/gorm"
)

// 审批
type ApprovalOrderService struct {
	*forms.ApprovalOrderForm
	C        *gin.Context
	Username string
}

func (s *ApprovalOrderService) Run() (err error) {
	// 判断工单是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单`%s`不存在", s.OrderID)
	}
	// 判断当前工单的审批状态
	if !utils.IsContain([]string{"PENDING"}, string(record.Progress)) {
		return fmt.Errorf("非待审批状态，禁止操作")
	}
	// 获取当前审批阶段的审批记录
	var approvalRecords []models.InsightApprovalRecords
	tx = global.App.DB.Table("`insight_approval_records` a").
		Where("a.order_id=? and stage=?", s.OrderID, record.Stage).
		Scan(&approvalRecords)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("审批记录`%s`不存在", s.OrderID)
	}
	// 判断用户是否有当前审批阶段的审批权限及是否已审批
	hasPermission := false
	for _, r := range approvalRecords {
		if r.Approver == s.Username {
			hasPermission = true
			if r.ApprovalStatus != "PENDING" {
				return fmt.Errorf("您已审批过，请勿重复执行")
			}
			break
		}
	}
	if !hasPermission {
		return fmt.Errorf("您没有审批权限或审批阶段未激活")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	txErr := global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新当前用户审批状态
		if err := tx.Model(&models.InsightApprovalRecords{}).
			Where("order_id=? AND stage=? AND approver=?", s.OrderID, record.Stage, s.Username).
			Updates(map[string]any{
				"approval_status": s.Status,
				"msg":             s.Msg,
				"approval_at":     now,
			}).Error; err != nil {
			return err
		}
		// 如果当前审核人驳回，直接驳回整个工单
		if s.Status == "REJECTED" {
			return tx.Model(&models.InsightOrderRecords{}).
				Where("order_id=?", s.OrderID).
				Updates(map[string]any{
					"progress": "REJECTED",
				}).Error
		}
		// 记录操作日志
		action := map[string]string{
			"APPROVED": "通过",
			"REJECTED": "驳回",
		}[s.Status]
		if err := WriteOrderLog(tx, s.OrderID, s.Username, fmt.Sprintf("用户%s%s了工单, 附加消息：%s", s.Username, action, s.Msg)); err != nil {
			global.App.Log.Error("ApprovalOrderService.Run error:", err.Error())
			return err
		}
		// 重新加载当前阶段记录，计算是否通过
		var stageRecords []models.InsightApprovalRecords
		if err := tx.Where("order_id=? AND stage=?", s.OrderID, record.Stage).Find(&stageRecords).Error; err != nil {
			return err
		}
		approvalType := stageRecords[0].ApprovalType
		allApproved := true
		anyApproved := false
		for _, r := range stageRecords {
			switch r.ApprovalStatus {
			case "APPROVED":
				anyApproved = true
			default:
				allApproved = false
			}
		}
		stagePass := false
		if approvalType == "AND" && allApproved {
			stagePass = true
		}
		if approvalType == "OR" && anyApproved {
			stagePass = true
		}
		if stagePass {
			// 当前阶段通过，检查是否还有下一阶段
			var nextStageCount int64
			if err := tx.Model(&models.InsightApprovalRecords{}).
				Where("order_id=? AND stage > ?", s.OrderID, record.Stage).
				Count(&nextStageCount).Error; err != nil {
				return err
			}
			if nextStageCount == 0 {
				// 没有下一阶段，全部审批完成
				return tx.Model(&models.InsightOrderRecords{}).
					Where("order_id=?", s.OrderID).
					Updates(map[string]any{
						"progress": "APPROVED",
					}).Error
			}
			// 有下一阶段，工单阶段 +1
			if err := tx.Model(&models.InsightOrderRecords{}).
				Where("order_id=?", s.OrderID).
				Update("stage", record.Stage+1).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", s.OrderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}

	// 如果阶段通过且有下一阶段，通知下一阶段审批人
	var nextStageRecords []models.InsightApprovalRecords
	if err := global.App.DB.Where("order_id=? AND stage=? AND approval_status='PENDING'", s.OrderID, latestRecord.Stage).Find(&nextStageRecords).Error; err == nil && len(nextStageRecords) > 0 {
		receivers := make([]string, 0, len(nextStageRecords))
		for _, r := range nextStageRecords {
			receivers = append(receivers, r.Approver)
		}
		notifier.SendOrderMessage(receivers, notifier.MsgTypeOrderPendingApproval, notifier.MessageParams{
			Order:     &latestRecord,
			Approvers: receivers,
		})
	}

	// 发送消息给申请人
	receiver := []string{latestRecord.Applicant}
	msgType := notifier.MsgTypeOrderApproved
	if s.Status == "REJECTED" {
		msgType = notifier.MsgTypeOrderRejected
	}
	notifier.SendOrderMessage(receiver, msgType, notifier.MessageParams{
		Order:         &latestRecord,
		Username:      s.Username,
		AdditionalMsg: s.Msg,
	})

	return nil
}

// 认领
type ClaimOrderService struct {
	*forms.ClaimOrderForm
	C        *gin.Context
	Username string
}

func (s *ClaimOrderService) Run() (err error) {
	// 判断工单是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单`%s`不存在", s.OrderID)
	}
	// 判断当前工单的审批状态
	if !utils.IsContain([]string{"APPROVED"}, string(record.Progress)) {
		return fmt.Errorf("当前工单没有审批通过，无法认领")
	}
	// 认领操作
	txErr := global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新工单认领人
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{
				"claimer":  s.Username,
				"progress": "CLAIMED",
			}).Error; err != nil {
			return err
		}
		// 记录操作日志
		if err := WriteOrderLog(tx, s.OrderID, s.Username, fmt.Sprintf("用户%s认领了工单，附加消息：%s", s.Username, s.Msg)); err != nil {
			global.App.Log.Error("ClaimOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", s.OrderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}
	receiver := []string{latestRecord.Applicant}
	notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderClaimed, notifier.MessageParams{
		Order:         &latestRecord,
		Username:      s.Username,
		AdditionalMsg: s.Msg,
	})

	return nil
}

// 转交
type TransferOrderService struct {
	*forms.TransferOrderForm
	C        *gin.Context
	Username string
}

func (s *TransferOrderService) Run() (err error) {
	// 判断工单是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单`%s`不存在", s.OrderID)
	}
	// 判断当前工单的审批状态
	if !utils.IsContain([]string{"CLAIMED"}, string(record.Progress)) {
		return fmt.Errorf("当前工单未被认领，无法转交")
	}
	// 判断当前工单认领人是否等于操作人
	if record.Claimer != s.Username {
		return fmt.Errorf("只有工单认领人才能转交工单")
	}
	// 转交操作
	txErr := global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新工单执行人
		if err != nil {
			return err
		}
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{
				"claimer": s.NewClaimer,
			}).Error; err != nil {
			return err
		}
		// 记录操作日志
		if err := WriteOrderLog(tx, s.OrderID, s.Username, fmt.Sprintf("用户%s转交工单给%s，附加消息：%s", s.Username, s.NewClaimer, s.Msg)); err != nil {
			global.App.Log.Error("TransferOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", s.OrderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}
	receiver := []string{latestRecord.Applicant, s.NewClaimer}
	notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderTransferred, notifier.MessageParams{
		Order:         &latestRecord,
		Username:      s.Username,
		AdditionalMsg: s.Msg,
	})

	return nil
}

// 撤销工单
type RevokeOrderService struct {
	*forms.RevokeOrderForm
	C        *gin.Context
	Username string
}

func (s *RevokeOrderService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 只有工单申请人才允许撤销工单
	if record.Applicant != s.Username {
		return fmt.Errorf("只有工单申请人才能撤销工单")
	}
	// 判断进度
	if !utils.IsContain([]string{"PENDING", "APPROVED", "CLAIMED"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	// 更新状态为已已撤销
	txErr := global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{
				"progress": "REVOKED",
			}).Error; err != nil {
			return err
		}
		// 记录操作日志
		if err := WriteOrderLog(tx, s.OrderID, s.Username, fmt.Sprintf("用户%s撤销了工单，附加消息：%s", s.Username, s.Msg)); err != nil {
			global.App.Log.Error("RevokeOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", s.OrderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}
	receiver := []string{latestRecord.Applicant}
	notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderRevoked, notifier.MessageParams{
		Order:         &latestRecord,
		Username:      s.Username,
		AdditionalMsg: s.Msg,
	})

	return nil
}

type CompleteOrderService struct {
	*forms.CompleteOrderForm
	C        *gin.Context
	Username string
}

func (s *CompleteOrderService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 判断进度
	if !utils.IsContain([]string{"CLAIMED", "EXECUTING"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	// 判断当前工单认领人是否等于操作人
	if record.Claimer != s.Username {
		return fmt.Errorf("只有工单认领人才能更改工单状态")
	}
	// 用户点击完成按钮
	txErr := global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{"progress": "COMPLETED"}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 记录操作日志
		if err := WriteOrderLog(tx, s.OrderID, s.Username, fmt.Sprintf("用户%s更改工单状态为已完成，附加消息：%s", s.Username, s.Msg)); err != nil {
			global.App.Log.Error("CompleteOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", s.OrderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}
	receiver := []string{latestRecord.Applicant}
	notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderCompleted, notifier.MessageParams{
		Order:         &latestRecord,
		Username:      s.Username,
		AdditionalMsg: s.Msg,
	})

	return nil
}

// 手动更新工单为失败
type FailOrderService struct {
	*forms.FailOrderForm
	C        *gin.Context
	Username string
}

func (s *FailOrderService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 判断进度
	if !utils.IsContain([]string{"CLAIMED", "EXECUTING"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	// 判断当前工单认领人是否等于操作人
	if record.Claimer != s.Username {
		return fmt.Errorf("只有工单认领人才能更改工单状态")
	}
	// 用户点击失败按钮
	txErr := global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{"progress": "FAILED"}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 记录操作日志
		if err := WriteOrderLog(tx, s.OrderID, s.Username, fmt.Sprintf("用户%s更改工单状态为已失败，附加消息：%s", s.Username, s.Msg)); err != nil {
			global.App.Log.Error("FailOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", s.OrderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}
	receiver := []string{latestRecord.Applicant}
	notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderFailed, notifier.MessageParams{
		Order:         &latestRecord,
		Username:      s.Username,
		AdditionalMsg: s.Msg,
	})

	return nil
}

// 复核
type ReviewOrderService struct {
	*forms.ReviewOrderForm
	C        *gin.Context
	Username string
}

func (s *ReviewOrderService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 判断进度
	if !utils.IsContain([]string{"COMPLETED"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	// 只有工单提交人才能复核
	if record.Applicant != s.Username {
		return fmt.Errorf("只有工单提交人才能复核工单")
	}

	txErr := global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{"progress": "REVIEWED"}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 记录操作日志
		if err := WriteOrderLog(tx, s.OrderID, s.Username, fmt.Sprintf("用户%s复核了工单，附加消息：%s", s.Username, s.Msg)); err != nil {
			global.App.Log.Error("ReviewOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", s.OrderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}
	receiver := []string{latestRecord.Applicant}
	notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderReviewed, notifier.MessageParams{
		Order:         &latestRecord,
		Username:      s.Username,
		AdditionalMsg: s.Msg,
	})

	return nil
}
