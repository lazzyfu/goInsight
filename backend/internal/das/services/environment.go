package services

import (
	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/internal/common/models"

	"github.com/gin-gonic/gin"
)

type GetEnvironmentsService struct {
	C *gin.Context
}

func (s *GetEnvironmentsService) Run() ([]models.InsightDBEnvironments, error) {
	var results []models.InsightDBEnvironments
	global.App.DB.Table("`insight_db_environments` a").
		Select("a.`name`").
		Scan(&results)
	return results, nil
}
