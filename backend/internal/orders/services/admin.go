package services

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lazzyfu/goinsight/internal/global"
	"gorm.io/gorm"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/models"
	userModels "github.com/lazzyfu/goinsight/internal/users/models"
	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminGetApprovalFlowUnboundService struct {
	C *gin.Context
	*forms.AdminGetApprovalFlowUnboundUsersForm
}

func (s *AdminGetApprovalFlowUnboundService) Run() (responseData any, total int64, err error) {
	var records []userModels.InsightUsers

	tx := global.App.DB.Table("insight_users a").
		Where("NOT EXISTS (SELECT 1 FROM insight_approval_flow_users b WHERE b.username = a.username)")

	// 搜索（精确/模糊看你需要）
	if s.Search != "" {
		like := "%" + s.Search + "%"
		tx = tx.Where("a.username LIKE ? OR a.nick_name LIKE ?", like, like)
	}

	total = pagination.Pager(&s.PaginationQ, tx, &records)
	return &records, total, nil
}

// 获取环境
type AdminGetApprovalFlowsService struct {
	C *gin.Context
	*forms.AdminGetApprovalFlowsForm
}

func (s *AdminGetApprovalFlowsService) Run() (responseData any, total int64, err error) {
	var records []models.InsightApprovalFlows

	// 基础查询（默认不 JOIN）
	base := global.App.DB.Table("insight_approval_flow a")

	// 有搜索时才 JOIN maps
	if s.Search != "" {
		base = base.Joins("LEFT JOIN insight_approval_flow_users b ON a.approval_id = b.approval_id").
			Where("a.name LIKE ? OR b.username = ?", "%"+s.Search+"%", s.Search)
	}

	// 固定字段 + 去重
	base = base.Select("a.id, a.name, a.definition, a.claim_users, a.approval_id, a.created_at, a.updated_at").
		Group("a.id")

	total = pagination.Pager(&s.PaginationQ, base, &records)
	return &records, total, nil
}

type AdminUpdateApprovalFlowsService struct {
	C *gin.Context
	*forms.AdminUpdateApprovalFlowsForm
	ID uint64
}

func (s *AdminUpdateApprovalFlowsService) Run() (responseData any, total int64, err error) {
	claimUsers, err := marshalClaimUsers(s.ClaimUsers)
	if err != nil {
		return nil, 0, err
	}
	// 更新记录
	result := global.App.DB.Model(&models.InsightApprovalFlows{}).Where("id=?", s.ID).Updates(map[string]any{
		"definition":  s.Definition,
		"name":        s.Name,
		"claim_users": claimUsers,
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
	claimUsers, err := marshalClaimUsers(s.ClaimUsers)
	if err != nil {
		return err
	}
	flow := models.InsightApprovalFlows{
		Definition: s.Definition,
		ClaimUsers: claimUsers,
		Name:       s.Name,
		ApprovalID: uuid.New(),
	}

	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightApprovalFlows{}).Create(&flow).Error; err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) {
				switch mysqlErr.Number {
				case 1062:
					return fmt.Errorf("审批流`%s`已存在", s.Name)
				}
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
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightApprovalFlows{})
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
		var record models.InsightApprovalFlowUsers
		tx := global.App.DB.Table("insight_approval_flow_users").Where("username=?", username).First(&record)
		if tx.RowsAffected != 0 {
			return fmt.Errorf("用户%s已经绑定审批流", username)
		}
	}

	// 批量构造写入数据
	records := make([]models.InsightApprovalFlowUsers, 0, len(s.Users))
	for _, u := range s.Users {
		records = append(records, models.InsightApprovalFlowUsers{
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
	var records []models.InsightApprovalFlowUsers
	tx := global.App.DB.Table("insight_approval_flow_users").Where("approval_id=?", s.ApprovalID)

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
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightApprovalFlowUsers{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
