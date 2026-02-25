package services

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (s *GetOrderEnvironmentsService) Run() ([]commonModels.InsightInstanceEnvironments, error) {
	var results []commonModels.InsightInstanceEnvironments
	global.App.DB.Table("`insight_instance_environments` a").
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
	var organization usersModels.InsightOrgs
	global.App.DB.Table("`insight_orgs` a").
		Joins("join insight_org_users b on a.key = b.organization_key").
		Joins("join insight_users c on c.uid = b.uid").Where("c.username=?", s.Username).Scan(&organization)

	// 将path json数据转换为数组
	var pathJsonArray []string
	if len(organization.Path) != 0 {
		_ = json.Unmarshal([]byte(organization.Path), &pathJsonArray)
	}
	pathJsonArray = append(pathJsonArray, organization.Key)

	// 获取当前用户当前绑定组织和所有上级组织绑定的实例
	var instances []commonModels.InsightInstances
	tx := global.App.DB.Table("insight_instances").
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
	var roles []commonModels.InsightInstanceSchemas
	tx := global.App.DB.Table("insight_instance_schemas").Where("instance_id=?", s.InstanceID)
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

	var record models.InsightApprovalFlows
	flow := tx.Table("insight_approval_flow_users a").Select(`b.definition`).
		Joins("inner join insight_approval_flow b on a.approval_id = b.approval_id").
		Where("a.username = ?", s.Username).Take(&record)
	if flow.Error != nil || flow.RowsAffected == 0 {
		return fmt.Errorf("未找到您的审批流配置，请联系管理员设置")
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
func (s *CreateOrderService) inspectSQL(instanceCfg commonModels.InsightInstances) ([]checker.ReturnData, error) {
	plainPassword, err := utils.Decrypt(instanceCfg.Password)
	if err != nil {
		return nil, err
	}
	inspect := checker.SyntaxInspectService{
		C:          s.C,
		InstanceID: instanceCfg.InstanceID,
		DbUser:     instanceCfg.User,
		DbPassword: plainPassword,
		DbHost:     instanceCfg.Hostname,
		DbPort:     instanceCfg.Port,
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
							insight_orgs ia
						WHERE
							EXISTS (
								SELECT
									1
								FROM
									insight_orgs
								WHERE
									JSON_CONTAINS(c.path, CONCAT('\"', ia.key, '\"'))
							)
					),
					'/',
					c.name
				),
				c.name
			) as organization`).
		Joins("left join insight_org_users b on a.uid = b.uid").
		Joins("left join insight_orgs c on b.organization_key = c.key").
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
	var config commonModels.InsightInstances
	global.App.DB.Table("`insight_instances`").
		Where("instance_id=?", s.InstanceID).
		First(&config)

	// 检查DDL/DML工单语法检查是否通过
	// 不对EXPORT工单进行语法检查，CheckSqlType已经要求EXPORT工单只能为SELECT语句
	if s.SQLType != "EXPORT" {
		returnData, err := s.inspectSQL(config)
		if err != nil {
			return err
		}

		// 检查语法检查是否通过
		// status: 0表示语法检查通过，1表示语法检查不通过
		status := 0
		for _, row := range returnData {
			for _, sum := range row.Summary {
				if sum.Level != "INFO" {
					status = 1
					break
				}
			}
			if status == 1 {
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
	// 获取工单环境名称
	var env commonModels.InsightInstanceEnvironments
	global.App.DB.Table("`insight_instance_environments` a").
		Select("a.`name`, a.id").
		Where("a.id=?", s.Environment).
		Take(&env)
	// 创建工单
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	title := fmt.Sprintf("%s_[%s] ", s.Title, timeStr)
	record := models.InsightOrderRecords{
		Title:            title,
		OrderID:          orderID,
		Remark:           s.Remark,
		SQLType:          s.SQLType,
		DBType:           s.DBType,
		Environment:      env.Name,
		InstanceID:       instance_id,
		Schema:           s.Schema,
		Applicant:        s.Username,
		Organization:     s.getUserOrg(),
		CC:               cc,
		Content:          s.Content,
		ExportFileFormat: s.ExportFileFormat,
	}
	err = global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 插入工单记录
		if err := tx.Model(&models.InsightOrderRecords{}).Create(&record).Error; err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) {
				switch mysqlErr.Number {
				case 1062:
					return fmt.Errorf("记录已存在，错误:%s", err.Error())
				}
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
		if err := WriteOrderLog(tx, orderID.String(), s.Username, fmt.Sprintf("用户%s提交了工单", s.Username)); err != nil {
			global.App.Log.Error("CreateOrderService.Run error:", err.Error())
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 事务提交成功后，重新查询最新记录后发送消息
	var latestRecord models.InsightOrderRecords
	if err := global.App.DB.Where("order_id = ?", orderID).Take(&latestRecord).Error; err != nil {
		global.App.Log.Error("查询工单记录失败:", err)
		return err
	}

	// 发送消息给工单申请人
	receiver := []string{latestRecord.Applicant}
	receiver = append(receiver, s.CC...)
	notifier.SendOrderMessage(receiver, notifier.MsgTypeOrderSubmitted, notifier.MessageParams{
		Order:    &latestRecord,
		Username: s.Username,
	})

	// 发送消息，发送给第一阶段的审批人
	var approvalRecords []models.InsightApprovalRecords
	global.App.DB.Table("insight_approval_records").
		Where("order_id = ? AND stage = ?", orderID, 1).
		Find(&approvalRecords)
	var approvers []string
	for _, appr := range approvalRecords {
		approvers = append(approvers, appr.Approver)
	}
	notifier.SendOrderMessage(approvers, notifier.MsgTypeOrderPendingApproval, notifier.MessageParams{
		Order:     &latestRecord,
		Approvers: approvers,
	})

	return nil
}
