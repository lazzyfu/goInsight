package services

import (
	"encoding/json"
	"fmt"
	"goInsight/global"
	"goInsight/internal/app/orders/forms"
	"goInsight/internal/app/orders/models"
	"goInsight/internal/pkg/notifier"
	"goInsight/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 审批
type ApproveService struct {
	*forms.ApproveForm
	C        *gin.Context
	Username string
}

func (s *ApproveService) updateProgress(tx *gorm.DB, progress string) error {
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

func (s *ApproveService) updateApprover(tx *gorm.DB, users []map[string]interface{}) error {
	usersJson, err := json.Marshal(users)
	if err != nil {
		return err
	}
	if err := tx.Model(&models.InsightOrderRecords{}).
		Where("order_id=?", s.OrderID).
		Updates(map[string]interface{}{
			"approver":   string(usersJson),
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		}).Error; err != nil {
		return err
	}
	return nil
}

func (s *ApproveService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 判断审核状态
	if !utils.IsContain([]string{"待审核"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	// 获取允许审核的用户
	var approverList []map[string]interface{}
	err = json.Unmarshal([]byte(record.Approver), &approverList)
	if err != nil {
		return err
	}
	// 更新审核人信息
	var M int = 0
	var passCount int = 0
	for _, i := range approverList {
		if i["user"] == s.Username {
			M += 1
			if i["status"] != "pending" {
				return fmt.Errorf("您已审核过，请不要重复执行")
			}
			i["status"] = s.Status
			i["user"] = s.Username
		}
		if i["status"] == "pass" {
			// 计算为审核通过
			passCount += 1
		}
	}
	if M == 0 {
		return fmt.Errorf("您没有当前工单的审核权限")
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新审批人信息
		if err := s.updateApprover(tx, approverList); err != nil {
			return err
		}
		// 全部都通过了审核，设置为已批准
		if len(approverList) == passCount {
			if err := s.updateProgress(tx, "已批准"); err != nil {
				return err
			}
		}
		// 如果点击驳回，将工单设置为已驳回
		if s.Status == "reject" {
			if err := s.updateProgress(tx, "已驳回"); err != nil {
				return err
			}
		}
		// 操作日志
		logMsg := fmt.Sprintf("用户%s审核通过了工单", s.Username)
		if s.Status == "reject" {
			logMsg = fmt.Sprintf("用户%s驳回了工单", s.Username)
		}
		if err := CreateOpLogs(tx, record.OrderID, s.Username, fmt.Sprintf("%s，附加消息：%s", logMsg, s.Msg)); err != nil {
			return err
		}
		// 发送消息，发送给工单申请人
		receiver := []string{record.Applicant}
		msg := fmt.Sprintf("您好，%s\n>工单标题：%s\n>附加消息：%s", logMsg, record.Title, s.Msg)
		notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
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
		// 操作日志
		logMsg := fmt.Sprintf("用户%s更新工单状态为%s，附加消息：%s", s.Username, s.Progress, s.Msg)
		if err := CreateOpLogs(tx, record.OrderID, s.Username, logMsg); err != nil {
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
		// 操作日志
		logMsg := fmt.Sprintf("用户%s复核了工单，附加消息：%s", s.Username, s.Msg)
		if err := CreateOpLogs(tx, record.OrderID, s.Username, logMsg); err != nil {
			return err
		}
		// 发送消息，发送给工单申请人
		receiver := []string{record.Applicant}
		msg := fmt.Sprintf("您好，用户%s更新工单状态为：已复核\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
	})
}

// 关闭工单
type CloseService struct {
	*forms.CloseForm
	C        *gin.Context
	Username string
}

func (s *CloseService) updateProgress(tx *gorm.DB, progress string) error {
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

func (s *CloseService) Run() (err error) {
	// 判断记录是否存在
	var record models.InsightOrderRecords
	tx := global.App.DB.Table("`insight_order_records`").Where("order_id=?", s.OrderID).Take(&record)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("记录`%s`不存在", s.OrderID)
	}
	// 判断进度
	if !utils.IsContain([]string{"待审核", "已批准", "执行中"}, string(record.Progress)) {
		return fmt.Errorf("非可操作状态，禁止操作")
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新状态为已关闭
		if err := s.updateProgress(tx, "已关闭"); err != nil {
			return err
		}
		// 操作日志
		logMsg := fmt.Sprintf("用户%s关闭了工单，附加消息：%s", s.Username, s.Msg)
		if err := CreateOpLogs(tx, record.OrderID, s.Username, logMsg); err != nil {
			return err
		}
		// 发送消息，发送给工单申请人
		receiver := []string{record.Applicant}
		msg := fmt.Sprintf("您好，用户%s关闭了工单\n>工单标题：%s\n>附加消息：%s", s.Username, record.Title, s.Msg)
		notifier.SendMessage(record.Title, record.OrderID.String(), receiver, msg)
		return nil
	})
}
