package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/notifier"
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
				"claimer":    s.Username,
				"progress":   "CLAIMED",
				"claimed_at": time.Now().Format("2006-01-02 15:04:05"),
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

// 关闭工单
type CloseOrderService struct {
	*forms.CloseOrderForm
	C        *gin.Context
	Username string
}

func (s *CloseOrderService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 只有工单申请人才允许关闭
	if record.Applicant != s.Username {
		return fmt.Errorf("只有工单申请人才能关闭工单")
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
				"closer":    s.Username,
				"progress":  "CLOSED",
				"closed_at": time.Now().Format("2006-01-02 15:04:05"),
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
			Msg:      fmt.Sprintf("用户%s关闭了工单", s.Username),
		}).Error; err != nil {
			global.App.Log.Error("CloseOrderService.Run error:", err.Error())
			return err
		}

		return nil

		// 发送消息，发送给工单申请人
		// receiver := []string{record.Applicant}
		// msg := fmt.Sprintf("您好，用户%s关闭了工单\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		// notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
	})
}

// 反馈
type FeedbackService struct {
	*forms.FeedbackForm
	C        *gin.Context
	Username string
}

func (s *FeedbackService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 已批准->已完成，已批准->执行中，执行中->已完成
	if !utils.IsContain([]string{"已批准", "执行中"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	// 用户点击反馈按钮
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).
			Where("order_id=?", s.OrderID).
			Updates(map[string]interface{}{"progress": s.Progress, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 发送消息，发送给工单申请人
		receiver := []string{record.Applicant}
		msg := fmt.Sprintf("您好，用户%s更新工单状态为：%s\n>工单标题：%s\n>附加消息：%s", s.Username, s.Progress, record.Title, s.Msg)
		notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
	})
}

// 复核
type ReviewService struct {
	*forms.ReviewForm
	C        *gin.Context
	Username string
}

func (s *ReviewService) updateProgress(tx *gorm.DB, progress string) error {
	if err := tx.Model(&models.InsightOrderRecords{}).
		Where("order_id=?", s.OrderID).
		Updates(map[string]interface{}{
			"progress":   progress,
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		}).Error; err != nil {
		return err
	}
	return nil
}

func (s *ReviewService) updateReviewer(tx *gorm.DB, users []map[string]interface{}) error {
	usersJson, err := json.Marshal(users)
	if err != nil {
		return err
	}
	if err := tx.Model(&models.InsightOrderRecords{}).
		Where("order_id=?", s.OrderID).
		Updates(map[string]interface{}{
			"reviewer":   string(usersJson),
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		}).Error; err != nil {
		return err
	}
	return nil
}

func (s *ReviewService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 判断进度
	if !utils.IsContain([]string{"已完成"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	// 获取复核人
	var reviewerList []map[string]interface{}
	err = json.Unmarshal([]byte(record.Reviewer), &reviewerList)
	if err != nil {
		return err
	}
	// 更新复核人信息
	var M int = 0
	var passCount int = 0
	for _, i := range reviewerList {
		if i["user"] == s.Username {
			M += 1
			if i["status"] != "pending" {
				return fmt.Errorf("您已复核过，请不要重复执行")
			}
			i["status"] = "pass"
			i["user"] = s.Username
		}
		if i["status"] == "pass" {
			// 计算为审核通过
			passCount += 1
		}
	}
	if M == 0 {
		return fmt.Errorf("您没有当前工单的复核权限")
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新审批人信息
		if err := s.updateReviewer(tx, reviewerList); err != nil {
			return err
		}
		// 全部都通过了复核，设置为已复核
		if len(reviewerList) == passCount {
			if err := s.updateProgress(tx, "已复核"); err != nil {
				return err
			}
		}
		// 发送消息，发送给工单申请人
		receiver := []string{record.Applicant}
		msg := fmt.Sprintf("您好，用户%s更新工单状态为：已复核\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
	})
}
