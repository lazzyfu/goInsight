/*
@Time    :   2023/08/03 16:05:17
@Author  :   zongfei.fu
@Desc    :
*/

package services

import (
	"encoding/json"
	"fmt"
	"goInsight/global"
	commonModels "goInsight/internal/apps/common/models"
	"goInsight/internal/apps/inspect/checker"
	"goInsight/internal/apps/orders/forms"
	"goInsight/internal/apps/orders/models"
	usersModels "goInsight/internal/apps/users/models"
	"goInsight/internal/pkg/notifier"
	"goInsight/internal/pkg/pagination"
	"goInsight/internal/pkg/parser"
	"goInsight/internal/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 获取环境
type GetEnvironmentsService struct {
	C *gin.Context
}

func (s *GetEnvironmentsService) Run() ([]commonModels.InsightDBEnvironments, error) {
	var results []commonModels.InsightDBEnvironments
	global.App.DB.Table("`insight_db_environments` a").
		Select("a.`name`, a.id").
		Scan(&results)
	return results, nil
}

// 获取指定环境的实例
type GetInstancesService struct {
	*forms.GetInstancesForm
	C        *gin.Context
	Username string
}

func (s *GetInstancesService) Run() (responseData interface{}, total int64, err error) {
	// 获取当前用户当前绑定组织和所有上级组织
	var organization usersModels.InsightOrganizations
	global.App.DB.Table("`insight_organizations` a").
		Joins("join insight_organizations_users b on a.key = b.organization_key").
		Joins("join insight_users c on c.uid = b.uid").Where("c.username=?", s.Username).Scan(&organization)

	// 将path json数据转换为数组
	var pathJsonArray []string
	if len(organization.Path) != 0 {
		_ = json.Unmarshal([]byte(organization.Path), &pathJsonArray)
	}
	pathJsonArray = append(pathJsonArray, organization.Key)

	// 获取当前用户当前绑定组织和所有上级组织绑定的实例
	var instances []commonModels.InsightDBConfig
	tx := global.App.DB.Table("insight_db_config").
		Where("environment=? and db_type=? and use_type='工单' and organization_key in ?", s.ID, s.DbType, pathJsonArray)
	total = pagination.Pager(&s.PaginationQ, tx, &instances)
	return &instances, total, nil
}

// 获取指定实例的Schemas
type GetSchemasService struct {
	*forms.GetSchemasForm
	C *gin.Context
}

func (s *GetSchemasService) Run() (responseData interface{}, total int64, err error) {
	var roles []commonModels.InsightDBSchemas
	tx := global.App.DB.Table("insight_db_schemas").Where("instance_id=?", s.InstanceID)
	total = pagination.Pager(&s.PaginationQ, tx, &roles)
	return &roles, total, nil
}

// 获取审核/复核/抄送人
type GetUsersService struct {
	*forms.GetUsersForm
	C *gin.Context
}

func (s *GetUsersService) Run() (responseData interface{}, total int64, err error) {
	var roles []usersModels.InsightUsers
	tx := global.App.DB.Table("insight_users").Where("is_active=1")
	total = pagination.Pager(&s.PaginationQ, tx, &roles)
	return &roles, total, nil
}

// 提交工单
type CreateOrdersService struct {
	*forms.CreateOrderForm
	C        *gin.Context
	Username string
	Audit    *parser.TiStmt
}

// 转json
func (s *CreateOrdersService) toJson(values []string) (datatypes.JSON, error) {
	var tmpData []map[string]interface{}
	for _, u := range values {
		// user：审核人/复核人
		// status : pending, approved, reject
		// time：审核时间
		tmpData = append(tmpData, map[string]interface{}{"user": u, "status": "pending"})
	}
	data, err := json.Marshal(tmpData)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(data), nil
}

// 审核SQL
func (s *CreateOrdersService) inspectSQL(config commonModels.InsightDBConfig) ([]checker.ReturnData, error) {
	inspect := checker.SyntaxInspectService{
		C:          s.C,
		DbUser:     global.App.Config.RemoteDB.UserName,
		DbPassword: global.App.Config.RemoteDB.Password,
		DbHost:     config.Hostname,
		DbPort:     config.Port,
		DBParams:   config.InspectParams,
		DBSchema:   s.Schema,
		Username:   s.Username,
		SqlText:    s.Content,
	}
	return inspect.Run()
}

func (s *CreateOrdersService) Run() error {
	// 判断SQL类型是否匹配，DML工单仅允许提交DML语句，DDL工单仅允许提交DDL语句
	err := parser.CheckSqlType(s.Content, string(s.SQLType))
	if err != nil {
		return err
	}
	// 判断SQL条数
	err = parser.CheckMaxAllowedSQLNums(s.Content)
	if err != nil {
		return err
	}
	// 获取实例配置
	var config commonModels.InsightDBConfig
	global.App.DB.Table("`insight_db_config`").
		Where("instance_id=?", s.InstanceID).
		First(&config)
	// 检查DDL/DML工单语法检查是否通过
	// 不对EXPORT工单进行语法检查，CheckSqlType已经要求EXPORT工单只能为SELECT语句
	if s.SQLType != "EXPORT" {
		returnData, err := s.inspectSQL(config)
		if err != nil {
			return err
		}
		// status: 0表示语法检查通过，1表示语法检查不通过
		status := 0
		for _, row := range returnData {
			if row.Level != "INFO" {
				status = 1
				break
			}
		}
		if status == 1 {
			return fmt.Errorf("SQL语法检查不通过，请先执行【语法检查】")
		}
	}
	// 解析UUID
	instance_id, err := utils.ParserUUID(s.InstanceID)
	if err != nil {
		return err
	}
	// 解析为json格式
	approver, err := s.toJson(s.Approver)
	if err != nil {
		return err
	}
	reviewer, err := s.toJson(s.Reviewer)
	if err != nil {
		return err
	}
	executorData, err := json.Marshal(s.Executor)
	if err != nil {
		return err
	}
	executor := datatypes.JSON(executorData)
	ccData, err := json.Marshal(s.CC)
	if err != nil {
		return err
	}
	cc := datatypes.JSON(ccData)
	// 生成工单ID
	orderID := uuid.New()
	// Title加上时间
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	title := fmt.Sprintf("%s_[%s] ", s.Title, timeStr)
	record := models.InsightOrderRecords{
		Title:            title,
		OrderID:          orderID,
		Remark:           s.Remark,
		IsRestrictAccess: *s.IsRestrictAccess,
		SQLType:          s.SQLType,
		DBType:           s.DBType,
		Environment:      s.Environment,
		InstanceID:       instance_id,
		Schema:           s.Schema,
		Applicant:        s.Username,
		Approver:         approver,
		Reviewer:         reviewer,
		Executor:         executor,
		CC:               cc,
		Content:          s.Content,
		ExportFileFormat: s.ExportFileFormat,
	}
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InsightOrderRecords{}).Create(&record).Error; err != nil {
			mysqlErr := err.(*mysql.MySQLError)
			switch mysqlErr.Number {
			case 1062:
				return fmt.Errorf("记录已存在，错误:%s", err.Error())
			}
			global.App.Log.Error(err)
			return err
		}
		// 操作日志
		log := models.InsightOrderOpLogs{
			Username: s.Username,
			OrderID:  record.OrderID,
			Msg:      fmt.Sprintf("用户%s提交了工单", s.Username),
		}
		if err := tx.Model(&models.InsightOrderOpLogs{}).Create(&log).Error; err != nil {
			global.App.Log.Error(err)
			return err
		}
		// 获取提交的环境
		var env commonModels.InsightDBEnvironments
		global.App.DB.Table("`insight_db_environments` a").
			Select("a.`name`, a.id").
			Take(&env)
		// 发送消息，发送给工单申请人
		receiver := []string{record.Applicant}
		receiver = append(receiver, s.Approver...)
		receiver = append(receiver, s.Reviewer...)
		receiver = append(receiver, s.CC...)

		msg := fmt.Sprintf(
			"您好，用户%s提交了工单，(*￣︶￣)\n"+
				">工单标题：%s\n"+
				">备注：%s\n"+
				">审核人：%s\n"+
				">复核人：%s\n"+
				">执行人：%s\n"+
				">抄送：%s\n"+
				">环境：%s\n"+
				">数据库类型：%s\n"+
				">工单类型：%s\n"+
				">库名：%s",
			s.Username, title, s.Remark,
			strings.Join(s.Approver, ","), strings.Join(s.Reviewer, ","), strings.Join(s.Executor, ","), strings.Join(s.CC, ","),
			env.Name, s.DBType, s.SQLType, s.Schema,
		)

		notifier.SendMessage(title, record.OrderID.String(), receiver, msg)
		return nil
	})
}
