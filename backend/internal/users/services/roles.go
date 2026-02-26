package services

import (
	"errors"
	"fmt"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/users/forms"
	"github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type GetRolesServices struct {
	*forms.GetRolesForm
	C *gin.Context
}

func (s *GetRolesServices) Run() (responseData interface{}, total int64, err error) {
	var roles []models.InsightRoles
	tx := global.App.DB.Table("insight_roles")
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`name` like ? ", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &roles)
	return &roles, total, nil
}

type CreateRolesService struct {
	*forms.CreateRolesForm
	C *gin.Context
}

func (s *CreateRolesService) Run() error {
	tx := global.App.DB.Model(&models.InsightRoles{})
	roles := models.InsightRoles{Name: s.Name}
	result := tx.Create(&roles)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录`%s`已存在", s.Name)
			}
		}
		return result.Error
	}
	return nil
}

type UpdateRolesService struct {
	*forms.UpdateRolesForm
	C  *gin.Context
	ID uint64
}

func (s *UpdateRolesService) Run() error {
	result := global.App.DB.Model(&models.InsightRoles{}).Where("id=?", s.ID).Updates(map[string]interface{}{
		"name": s.Name,
	})
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录`%s`已存在", s.Name)
			}
		}
		return result.Error
	}
	return nil
}

type DeleteRolesService struct {
	C  *gin.Context
	ID uint64
}

func (s *DeleteRolesService) Run() error {
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightRoles{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
