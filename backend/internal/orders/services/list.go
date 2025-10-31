package services

import (
	"encoding/json"
	"fmt"

	"github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/orders/forms"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/pagination"
	"github.com/lazzyfu/goinsight/pkg/utils"

	ordersModels "github.com/lazzyfu/goinsight/internal/orders/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type GetOrderListServices struct {
	*forms.GetOrderListForm
	C        *gin.Context
	Username string
}

func (s *GetOrderListServices) Run() (responseData interface{}, total int64, err error) {
	type record struct {
		Title            string           `json:"title"`
		Progress         string           `json:"progress"`
		IsRestrictAccess bool             `json:"is_restrict_access"`
		Applicant        string           `json:"applicant"`
		Organization     string           `json:"organization"`
		Environment      string           `json:"environment"`
		SqlType          string           `json:"sql_type"`
		Instance         string           `json:"instance"`
		Schema           string           `json:"schema"`
		Approver         datatypes.JSON   `json:"approver"`
		Reviewer         datatypes.JSON   `json:"reviewer"`
		OrderID          string           `json:"order_id"`
		CreatedAt        models.LocalTime `json:"created_at"`
	}
	var records []record
	tx := global.App.DB.Table("insight_order_records a").
		Select(`
			a.progress, 
			a.title as title, 
			a.applicant, 
			if(length(a.organization)=0, "N/A", a.organization) as organization,
			a.is_restrict_access,
			b.name as environment,
			concat(c.hostname, ':', c.port) as instance, 
			a.schema, 
			a.sql_type, 
			a.approver, 
			a.reviewer, 
			a.order_id, 
			a.created_at
		`).
		Joins("left join insight_db_environments b on a.environment=b.id").
		Joins("left join insight_db_config c on a.instance_id = c.instance_id").
		Order("a.created_at desc")
	// 仅加载我的工单
	if s.OnlyMyOrders {
		tx = tx.Where("a.applicant=?", s.Username)
	}
	// 搜索
	if s.Search != "" {
		tx = tx.Where("a.title like ?", "%"+s.Search+"%")
	}
	if s.Progress != "" {
		tx = tx.Where("a.progress=?", s.Progress)
	}
	if s.Environment > 0 {
		tx = tx.Where("a.environment=?", s.Environment)
	}
	total = pagination.Pager(&s.PaginationQ, tx, &records)
	return &records, total, nil
}

type GetOrderDetailServices struct {
	C        *gin.Context
	OrderID  string
	Username string
}

func (s *GetOrderDetailServices) convertToList(data datatypes.JSON) (users []string) {
	var usersList []map[string]string
	err := json.Unmarshal([]byte(data), &usersList)
	if err != nil {
		global.App.Log.Error("GetDetailServices.convertToList", err.Error())
		return
	}
	for _, entry := range usersList {
		users = append(users, entry["user"])
	}
	return
}

func (s *GetOrderDetailServices) Run() (responseData interface{}, err error) {
	type record struct {
		ordersModels.InsightOrderRecords
		Environment string `json:"environment"`
		Instance    string `json:"instance"`
	}
	var result record
	// 返回记录
	tx := global.App.DB.Table("`insight_order_records` a").
		Select("a.*, b.name as environment, concat(c.hostname, ':', c.port) as instance").
		Joins("left join insight_db_environments b on a.environment=b.id").
		Joins("left join insight_db_config c on a.instance_id = c.instance_id").
		Where("a.order_id=?", s.OrderID).
		Take(&result)
	if tx.RowsAffected == 0 {
		return result, fmt.Errorf("记录`%s`不存在", s.OrderID)
	}

	// 限制访问
	if result.IsRestrictAccess {
		var users []string = []string{result.Applicant}
		users = append(users, s.convertToList(result.Approver)...)
		users = append(users, s.convertToList(result.Reviewer)...)
		users = append(users, s.convertToList(result.CC)...)
		if !utils.IsContain(users, s.Username) {
			result.Content = "您没有权限查看当前工单内容"
		}
	}
	return result, nil
}

type GetOrderApprovalServices struct {
	C        *gin.Context
	OrderID  string
	Username string
}

func (s *GetOrderApprovalServices) Run() (responseData interface{}, err error) {
	var result []ordersModels.InsightApprovalRecords
	// 返回记录
	tx := global.App.DB.Table("`insight_approval_records` a").
		Where("a.order_id=?", s.OrderID).
		Scan(&result)
	if tx.RowsAffected == 0 {
		return result, fmt.Errorf("记录`%s`不存在", s.OrderID)
	}

	return result, nil
}
