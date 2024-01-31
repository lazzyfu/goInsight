package services

import (
	"goInsight/global"
	"goInsight/internal/apps/orders/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 操作日志
func CreateOpLogs(tx *gorm.DB, order_id uuid.UUID, username, msg string) error {
	log := models.InsightOrderOpLogs{
		Username: username,
		OrderID:  order_id,
		Msg:      msg,
	}
	if err := tx.Model(&models.InsightOrderOpLogs{}).Create(&log).Error; err != nil {
		global.App.Log.Error(err)
		return err
	}
	return nil
}
