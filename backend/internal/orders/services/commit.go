package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/notifier"
	"github.com/lazzyfu/goinsight/pkg/pagination"
	"github.com/lazzyfu/goinsight/pkg/parser"
	"github.com/lazzyfu/goinsight/pkg/utils"

	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	"github.com/lazzyfu/goinsight/internal/inspect/checker"
	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/models"
	usersModels "github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 获取环境
type GetOrderEnvironmentsService struct {
	C *gin.Context
}

func (s *GetOrderEnvironmentsService) Run() ([]commonModels.InsightDBEnvironments, error) {
	var results []commonModels.InsightDBEnvironments
	global.App.DB.Table("`insight_db_environments` a").
		Select("a.`name`, a.id").
		Scan(&results)
	return results, nil
}

// 获取指定环境的实例
type GetOrderInstancesService struct {
	*forms.GetOrderInstancesForm
	C        *gin.Context
	Username string
}

func (s *GetOrderInstancesService) Run() (responseData interface{}, total int64, err error) {
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
type GetOrderSchemasService struct {
	*forms.GetOrderSchemasForm
	C *gin.Context
}

func (s *GetOrderSchemasService) Run() (responseData interface{}, total int64, err error) {
	var roles []commonModels.InsightDBSchemas
	tx := global.App.DB.Table("insight_db_schemas").Where("instance_id=?", s.InstanceID)
	total = pagination.Pager(&s.PaginationQ, tx, &roles)
	return &roles, total, nil
}

// 获取审核/复核/抄送人
type GetOrderUsersService struct {
	*forms.GetOrderUsersForm
	C *gin.Context
}

func (s *GetOrderUsersService) Run() (responseData interface{}, total int64, err error) {
	var roles []usersModels.InsightUsers
	tx := global.App.DB.Table("insight_users").Where("is_active=1")
	total = pagination.Pager(&s.PaginationQ, tx, &roles)
	return &roles, total, nil
}

// 提交工单
type CreateOrderService struct {
	*forms.CreateOrderForm
	C        *gin.Context
	Username string
	Audit    *parser.TiStmt
}

func (s *CreateOrderService) generateApprovalRecords(tx *gorm.DB, orderID uuid.UUID) error {
	type FlowStage struct {
		Type      string   `json:"type"`
		Stage     int      `json:"stage"`
		StageName string   `json:"stage_name"`
		Approvers []string `json:"approvers"`
	}

	var record models.InsightApprovalFlow
	flow := tx.Table("insight_approval_maps a").Select(`b.definition`).
		Joins("inner join insight_approval_flow b on a.approval_id = b.approval_id").
		Where("a.username = ?", s.Username).Take(&record)
	if flow.Error != nil || flow.RowsAffected == 0 {
		return fmt.Errorf("没有查询到审批流配置，请联系系统管理员")
	}
	// [{"type": "OR", "stage": 1, "approvers": ["zhangsan", "lisi"], "stage_name": "安全组审批"}, {"type": "AND", "stage": 1, "approvers": ["admin"], "stage_name": "部门负责人审批"}]
	var stages []FlowStage
	if err := json.Unmarshal(record.Definition, &stages); err != nil {
		return fmt.Errorf("解析审批流JSON失败: %w", err)
	}

	for _, s := range stages {
		for _, approver := range s.Approvers {
			audit := models.InsightApprovalRecords{
				OrderID:        orderID,
				Stage:          s.Stage,
				StageName:      s.StageName,
				Approver:       approver,
				ApprovalStatus: "PENDING",
				ApprovalType:   commonModels.EnumType(s.Type),
			}

			if err := tx.Create(&audit).Error; err != nil {
				return fmt.Errorf("创建审批记录失败: %w", err)
			}
		}
	}
	return nil
}

// 审核SQL
func (s *CreateOrderService) inspectSQL(config commonModels.InsightDBConfig) ([]checker.ReturnData, error) {
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

// 获取用户组织，没有返回空值
func (s *CreateOrderService) getUserOrg() (organization string) {
	type org struct {
		Organization string
	}
	var record org
	tx := global.App.DB.Table("insight_users a").Select(`
			a.username,
			ifnull(
				concat(
					(
						SELECT
							GROUP_CONCAT(
								ia.name
								ORDER BY
									ia.name ASC SEPARATOR '/'
							) AS concatenated_names
						FROM
							insight_organizations ia
						WHERE
							EXISTS (
								SELECT
									1
								FROM
									insight_organizations
								WHERE
									JSON_CONTAINS(c.path, CONCAT('\"', ia.key, '\"'))
							)
					),
					'/',
					c.name
				),
				c.name
			) as organization`).
		Joins("left join insight_organizations_users b on a.uid = b.uid").
		Joins("left join insight_organizations c on b.organization_key = c.key").
		Where("a.username = ?", s.Username).Take(&record)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return
	}
	return record.Organization
}

func (s *CreateOrderService) Run() error {
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
	// 抄送人
	ccData, err := json.Marshal(s.CC)
	if err != nil {
		return err
	}
	cc := datatypes.JSON(ccData)
	// 生成工单ID
	orderID := uuid.New()
	// 创建工单
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	title := fmt.Sprintf("%s_[%s] ", s.Title, timeStr)
	record := models.InsightOrderRecords{
		Title:            title,
		OrderID:          orderID,
		Remark:           s.Remark,
		SQLType:          s.SQLType,
		DBType:           s.DBType,
		Environment:      s.Environment,
		InstanceID:       instance_id,
		Schema:           s.Schema,
		Applicant:        s.Username,
		Organization:     s.getUserOrg(),
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
		// 生成审批流
		err = s.generateApprovalRecords(tx, orderID)
		if err != nil {
			return err
		}
		// 记录操作日志
		if err := tx.Create(&models.InsightOrderLogs{
			OrderID:  orderID,
			Username: s.Username,
			Msg:      fmt.Sprintf("用户%s提交了工单", s.Username),
		}).Error; err != nil {
			global.App.Log.Error("CreateOrderService.Run error:", err.Error())
			return err
		}
		// 获取提交的环境
		var env commonModels.InsightDBEnvironments
		global.App.DB.Table("`insight_db_environments` a").
			Select("a.`name`, a.id").
			Take(&env)
		// 发送消息，发送给工单申请人
		receiver := []string{record.Applicant}
		receiver = append(receiver, s.CC...)

		msg := fmt.Sprintf(
			"您好，用户%s提交了工单\n"+
				">工单标题：%s\n"+
				">备注：%s\n"+
				">抄送：%s\n"+
				">环境：%s\n"+
				">数据库类型：%s\n"+
				">工单类型：%s\n"+
				">库名：%s",
			s.Username, title, s.Remark,
			strings.Join(s.CC, ","),
			env.Name, s.DBType, s.SQLType, s.Schema,
		)

		notifier.SendMessage(title, record.OrderID.String(), receiver, msg)
		return nil
	})
}
