package services

import (
	"fmt"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/models"

	"github.com/gin-gonic/gin"
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
	var orderRecord models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&orderRecord)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单`%s`不存在", s.OrderID)
	}
	// 判断当前工单的审批状态
	if !utils.IsContain([]string{"PENDING"}, string(orderRecord.Progress)) {
		return fmt.Errorf("非待审批状态，禁止操作")
	}
	// 获取当前审批阶段的审批记录
	var approvalRecords []models.InsightApprovalRecords
	tx = global.App.DB.Table("`insight_approval_records` a").
		Where("a.order_id=? and stage=?", s.OrderID, orderRecord.Stage).
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
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新当前用户审批状态
		if err := tx.Model(&models.InsightApprovalRecords{}).
			Where("order_id=? AND stage=? AND approver=?", s.OrderID, orderRecord.Stage, s.Username).
			Updates(map[string]any{
				"approval_status": s.Status,
				"msg":             s.Msg,
				"approval_at":     now,
			}).Error; err != nil {
			return err
		}

		// 记录操作日志
		orderID, err := utils.ParserUUID(s.OrderID)
		if err != nil {
			return err
		}
		var action string
		switch s.Status {
		case "APPROVED":
			action = "通过"
		case "REJECTED":
			action = "驳回"
		}
		if err := tx.Create(&models.InsightOrderLogs{
			OrderID:  orderID,
			Username: s.Username,
			Msg:      fmt.Sprintf("用户%s%s了工单", s.Username, action),
		}).Error; err != nil {
			global.App.Log.Error("ApprovalOrderService.Run error:", err.Error())
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

		// 重新加载当前阶段记录，计算是否通过
		var stageRecords []models.InsightApprovalRecords
		if err := tx.Where("order_id=? AND stage=?", s.OrderID, orderRecord.Stage).Find(&stageRecords).Error; err != nil {
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
				Where("order_id=? AND stage > ?", s.OrderID, orderRecord.Stage).
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
				Update("stage", orderRecord.Stage+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 认领
type ClaimOrderService struct {
	*forms.ClaimOrderForm
	C        *gin.Context
	Username string
}

func (s *ClaimOrderService) Run() (err error) {
	// 判断工单是否存在
	var orderRecord models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&orderRecord)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单`%s`不存在", s.OrderID)
	}
	// 判断当前工单的审批状态
	if !utils.IsContain([]string{"APPROVED"}, string(orderRecord.Progress)) {
		return fmt.Errorf("当前工单没有审批通过，无法认领")
	}
	// 认领操作
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
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
		orderID, err := utils.ParserUUID(s.OrderID)
		if err != nil {
			return err
		}
		if err := tx.Create(&models.InsightOrderLogs{
			OrderID:  orderID,
			Username: s.Username,
			Msg:      fmt.Sprintf("用户%s认领了工单", s.Username),
		}).Error; err != nil {
			global.App.Log.Error("ClaimOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})
}

// 转交
type TransferOrderService struct {
	*forms.TransferOrderForm
	C        *gin.Context
	Username string
}

func (s *TransferOrderService) Run() (err error) {
	// 判断工单是否存在
	var orderRecord models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&orderRecord)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("工单`%s`不存在", s.OrderID)
	}
	// 判断当前工单的审批状态
	if !utils.IsContain([]string{"CLAIMED"}, string(orderRecord.Progress)) {
		return fmt.Errorf("当前工单未被认领，无法转交")
	}
	// 判断当前工单认领人是否等于操作人
	if orderRecord.Claimer != s.Username {
		return fmt.Errorf("只有工单认领人才能转交工单")
	}
	// 转交操作
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新工单执行人
		if err != nil {
			return err
		}
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{
				"executor": s.NewExecutor,
			}).Error; err != nil {
			return err
		}
		// 记录操作日志
		orderID, err := utils.ParserUUID(s.OrderID)
		if err != nil {
			return err
		}
		if err := tx.Create(&models.InsightOrderLogs{
			OrderID:  orderID,
			Username: s.Username,
			Msg:      fmt.Sprintf("用户%s转交工单给%s", s.Username, s.NewExecutor),
		}).Error; err != nil {
			global.App.Log.Error("TransferOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})
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
	// 更新状态为已关闭
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{
				"progress": "REVOKED",
			}).Error; err != nil {
			return err
		}

		// 记录操作日志
		orderID, err := utils.ParserUUID(s.OrderID)
		if err != nil {
			return err
		}
		if err := tx.Create(&models.InsightOrderLogs{
			OrderID:  orderID,
			Username: s.Username,
			Msg:      fmt.Sprintf("用户%s撤销了工单", s.Username),
		}).Error; err != nil {
			global.App.Log.Error("RevokeOrderService.Run error:", err.Error())
			return err
		}

		return nil

		// 发送消息，发送给工单申请人
		// receiver := []string{record.Applicant}
		// msg := fmt.Sprintf("您好，用户%s撤销了工单\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		// notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
	})
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
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{"progress": "COMPLETED"}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 发送消息，发送给工单申请人
		// receiver := []string{record.Applicant}
		// msg := fmt.Sprintf("您好，用户%s完成了工单\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		// notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
	})
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
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{"progress": "FAILED"}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 发送消息，发送给工单申请人
		// receiver := []string{record.Applicant}
		// msg := fmt.Sprintf("您好，用户%s更新工单为失败状态\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		// notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
	})
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

	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]any{"progress": "REVIEWED"}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 发送消息，发送给工单申请人
		// receiver := []string{record.Applicant}
		// msg := fmt.Sprintf("您好，用户%s更新工单状态为：已复核\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		// notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
	})
}
