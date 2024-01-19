/*
@Time    :   2023/08/03 16:05:17
@Author  :   zongfei.fu
@Desc    :
*/

package services

import (
	"goInsight/global"
	"goInsight/internal/app/common/models"

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
