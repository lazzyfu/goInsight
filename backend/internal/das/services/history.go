package services

import (
	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"

	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/models"

	"github.com/gin-gonic/gin"
)

type GetHistoryService struct {
	*forms.GetHistoryForm
	C        *gin.Context
	Username string
}

func (s *GetHistoryService) Run() (responseData *[]models.InsightDASRecords, total int64, err error) {
	var list []models.InsightDASRecords
	tx := global.App.DB.Model(&models.InsightDASRecords{}).Where("username=?", s.Username).Order("created_at desc")
	// 搜索schema
	if s.Search != "" {
		tx = tx.Where("`schema` like ? or `tables` like ? or `sqltext` like ?", "%"+s.Search+"%", "%"+s.Search+"%", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &list)

	return &list, total, nil
}
