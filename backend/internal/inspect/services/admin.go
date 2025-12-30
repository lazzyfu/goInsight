package services

import (
	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/inspect/forms"
	"github.com/lazzyfu/goinsight/internal/inspect/models"

	"github.com/gin-gonic/gin"
)

type AdminGlobalInspectParamsServices struct {
	*forms.AdminGlobalInspectParamsForm
	C *gin.Context
}

func (s *AdminGlobalInspectParamsServices) Run() (responseData any, total int64, err error) {
	var params []models.InsightInspectGlobalParams
	tx := global.App.DB.Model(&models.InsightInspectGlobalParams{})
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`title` like ?", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &params)
	return &params, total, nil
}

type AdminUpdateGlobalInspectParamsService struct {
	*forms.AdminUpdateGlobalInspectParamsForm
	C  *gin.Context
	ID uint64
}

func (s *AdminUpdateGlobalInspectParamsService) Run() error {
	// 只修改value
	tx := global.App.DB.Model(&models.InsightInspectGlobalParams{}).Where("id=? and `key`=?", s.ID, s.Key)
	result := tx.Updates(map[string]any{
		"value": s.Value,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
