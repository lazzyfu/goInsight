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
		Select("a.`id`, a.`name`, a.`definition`, a.`approval_id`, a.`created_at`, a.`updated_at`").
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

type AdminBindUsersToApprovalFlowService struct {
	*forms.AdminBindUsersToApprovalFlowForm
	C *gin.Context
}

func (s *AdminBindUsersToApprovalFlowService) Run() error {
	// 判断当前用户是否绑定审批流
	for _, username := range s.Users {
		var record models.InsightApprovalMaps
		tx := global.App.DB.Table("insight_approval_maps").Where("username=?", username).First(&record)
		if tx.RowsAffected != 0 {
			return fmt.Errorf("用户%s已经绑定审批流", username)
		}
	}

	// 批量构造写入数据
	records := make([]models.InsightApprovalMaps, 0, len(s.Users))
	for _, u := range s.Users {
		records = append(records, models.InsightApprovalMaps{
			ApprovalID: s.ApprovalID,
			Username:   u,
		})
	}

	// 批量写入（事务内）
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(&records, 100).Error
	})

}

type AdminGetApprovalFlowUsersService struct {
	C *gin.Context
	*forms.AdminGetApprovalFlowUsersForm
}

func (s *AdminGetApprovalFlowUsersService) Run() (responseData any, total int64, err error) {
	var records []models.InsightApprovalMaps
	tx := global.App.DB.Table("insight_approval_maps").Where("approval_id=?", s.ApprovalID).Scan(&records)

	// 搜索
	if s.Search != "" {
		tx = tx.Where("username like ?", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &records)
	return &records, total, nil
}

type AdminDeleteUsersFromApprovalFlowService struct {
	C  *gin.Context
	ID uint64
}

func (s *AdminDeleteUsersFromApprovalFlowService) Run() error {
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightApprovalMaps{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
