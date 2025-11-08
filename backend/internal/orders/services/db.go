package services

import (
	"github.com/lazzyfu/goinsight/internal/orders/models"
	"github.com/lazzyfu/goinsight/pkg/utils"
	"gorm.io/gorm"
)

func WriteOrderLog(tx *gorm.DB, order_id, username, msg string) error {
	// 记录操作日志
	orderID, err := utils.ParserUUID(order_id)
	if err != nil {
		return err
	}
	if err := tx.Create(&models.InsightOrderLogs{
		OrderID:  orderID,
		Username: username,
		Msg:      msg,
	}).Error; err != nil {
		return err
	}
	return nil
}
