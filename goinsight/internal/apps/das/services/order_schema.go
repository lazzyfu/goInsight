/*
@Time    :   2023/08/02 17:47:51
@Author  :   zongfei.fu
@Desc    :
*/

package services

import (
	"goInsight/global"
	"goInsight/internal/apps/das/forms"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderSchemasResult struct {
	InstanceID uuid.UUID `json:"instance_id"`
	Schema     string    `json:"schema"`
	Remark     string    `json:"remark"`
}

type GetOrderSchemasService struct {
	*forms.GetOrderSchemasForm
	C *gin.Context
}

func (s *GetOrderSchemasService) Run() ([]OrderSchemasResult, error) {
	var results []OrderSchemasResult
	global.App.DB.Table("`insight_db_schemas` a").
		Select("a.`instance_id`, a.`schema`, b.`remark`").
		Joins("join `insight_db_config` b on a.instance_id = b.instance_id").
		Joins("join `insight_db_environments` c on b.environment = c.id").
		Where("c.name=?", s.Environment).
		Group("a.`instance_id`, a.`schema`, b.`remark`").
		Scan(&results)
	return results, nil
}
