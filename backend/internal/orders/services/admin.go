package services

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lazzyfu/goinsight/internal/global"
	"gorm.io/gorm"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/models"
	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 获取环境
type AdminGetApprovalFlowsService struct {
	C *gin.Context
	*forms.AdminGetApprovalFlowsForm
}

func (s *AdminGetApprovalFlowsService) Run() (responseData any, total int64, err error) {
	var records []models.InsightApprovalFlow
	tx := global.App.DB.Table("`insight_approval_flow` a").
		Select("a.`id`, a.`name`, a.`definition`, a.`created_at`, a.`updated_at`").
		Scan(&records)
	// 搜索
	if s.Search != "" {
		tx = tx.Where("a.title like ?", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &records)
	return &records, total, nil
}

type AdminUpdateApprovalFlowsService struct {
	C *gin.Context
	*forms.AdminUpdateApprovalFlowsForm
	ID uint64
}

func (s *AdminUpdateApprovalFlowsService) Run() (responseData any, total int64, err error) {
	// 更新记录
	result := global.App.DB.Model(&models.InsightApprovalFlow{}).Where("id=?", s.ID).Updates(map[string]any{
		"definition": s.Definition,
		"name":       s.Name,
	})

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return nil, 0, nil
}

type AdminCreateApprovalFlowsService struct {
	C *gin.Context
	*forms.AdminCreateApprovalFlowsForm
}

func (s *AdminCreateApprovalFlowsService) Run() error {
	flow := models.InsightApprovalFlow{
		Definition: s.Definition,
		Name:       s.Name,
		ApprovalID: uuid.New(),
	}

	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightApprovalFlow{}).Create(&flow).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("审批流`%s`已存在", s.Name)
			}
			global.App.Log.Error(err)
			return err
		}
		return nil
	})
}

type AdminDeleteApprovalFlowsService struct {
	C  *gin.Context
	ID uint64
}

func (s *AdminDeleteApprovalFlowsService) Run() error {
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightApprovalFlow{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
