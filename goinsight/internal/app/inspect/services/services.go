/*
@Time    :   2023/08/31 15:19:46
@Author  :   lazzyfu
*/

package services

import (
	"fmt"
	"goInsight/global"
	"goInsight/internal/app/inspect/forms"
	"goInsight/internal/app/inspect/models"
	"goInsight/internal/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type AdminInspectParamsServices struct {
	*forms.AdminInspectParamsForm
	C *gin.Context
}

func (s *AdminInspectParamsServices) Run() (responseData interface{}, total int64, err error) {
	var params []models.InsightInspectParams
	tx := global.App.DB.Model(&models.InsightInspectParams{})
	// 搜索
	if s.Search != "" {
		tx = tx.Where("`remark` like ?", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &params)
	return &params, total, nil
}

type AdminUpdateInspectParamsService struct {
	*forms.AdminUpdateInspectParamsForm
	C  *gin.Context
	ID uint64
}

func (s *AdminUpdateInspectParamsService) Run() error {
	// 更新记录
	result := global.App.DB.Model(&models.InsightInspectParams{}).Where("id=?", s.ID).Updates(map[string]interface{}{
		"param":  s.Params,
		"remark": s.Remark,
	})
	if result.Error != nil {
		mysqlErr := result.Error.(*mysql.MySQLError)
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("%s记录已存在", s.Remark)
		}
		return result.Error
	}
	return nil
}
