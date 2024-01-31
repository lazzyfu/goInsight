/*
@Time    :   2023/08/31 15:19:46
@Author  :   lazzyfu
*/

package services

import (
	"encoding/json"
	"fmt"
	"goInsight/global"
	"goInsight/internal/app/inspect/forms"
	"goInsight/internal/app/inspect/models"
	"goInsight/internal/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
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
	// 拦截，避免客户端将KEY修改了
	for key := range s.Params {
		var num int64
		global.App.DB.Model(&models.InsightInspectParams{}).Select("json_contains_path(params, 'one', ?)", "$."+key).Where("id=?", s.ID).Scan(&num)
		if num == 0 {
			return fmt.Errorf("修改失败，禁止修改KEY：%s", key)
		}
	}

	// 审核参数
	jsonParams, err := json.Marshal(s.Params)
	if err != nil {
		return err
	}

	// 更新记录
	result := global.App.DB.Model(&models.InsightInspectParams{}).Where("id=?", s.ID).Updates(map[string]interface{}{
		"params": datatypes.JSON(jsonParams),
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
