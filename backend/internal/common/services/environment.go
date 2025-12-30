package services

import (
	"fmt"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/common/forms"
	"github.com/lazzyfu/goinsight/internal/common/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type AdminGetEnvironmentServices struct {
	*forms.AdminGetEnvironmentForm
	C *gin.Context
}

func (s *AdminGetEnvironmentServices) Run() (responseData interface{}, total int64, err error) {
	var environments []models.InsightInstanceEnvironments
	tx := global.App.DB.
		Table("insight_instance_environments").
		Order("updated_at desc")
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`name` like ?", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &environments)
	return &environments, total, nil
}

type AdminCreateEnvironmentService struct {
	*forms.AdminCreateEnvironmentForm
	C *gin.Context
}

func (s *AdminCreateEnvironmentService) Run() error {
	tx := global.App.DB.Model(&models.InsightInstanceEnvironments{})
	db := models.InsightInstanceEnvironments{
		Name: s.Name,
	}
	result := tx.Create(&db)
	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("记录`%s`已存在", s.Name)
		}
		return result.Error
	}
	return nil
}

type AdminUpdateEnvironmentService struct {
	*forms.AdminUpdateEnvironmentForm
	C  *gin.Context
	ID uint64
}

func (s *AdminUpdateEnvironmentService) Run() error {
	result := global.App.DB.Model(&models.InsightInstanceEnvironments{}).Where("id=?", s.ID).Updates(map[string]interface{}{
		"name": s.Name,
	})
	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("记录`%s`已存在", s.Name)
		}
		return result.Error
	}
	return nil
}

type AdminDeleteEnvironmentService struct {
	C  *gin.Context
	ID uint64
}

func (s *AdminDeleteEnvironmentService) Run() error {
	tx := global.App.DB.Where("id=?", s.ID).Delete(&models.InsightInstanceEnvironments{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
