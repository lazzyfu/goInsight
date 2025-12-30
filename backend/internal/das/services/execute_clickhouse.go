package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/query"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/models"
	"github.com/lazzyfu/goinsight/internal/das/parser"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"gorm.io/datatypes"
)

type ExecuteClickHouseQueryService struct {
	*forms.ExecuteClickHouseQueryForm
	C        *gin.Context
	Username string
}

func (s *ExecuteClickHouseQueryService) validateStatementSyntax() (*parser.TiStmt, error) {
	// 检查语法是否有效
	TiStmt, warns, err := parser.NewParse(s.SQLText, "", "")
	if len(warns) > 0 {
		return TiStmt, fmt.Errorf("Parse Warning: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return TiStmt, fmt.Errorf("SQL语法解析错误:%s", err.Error())
	}
	return TiStmt, nil
}

func (s *ExecuteClickHouseQueryService) validateStatementNum(stmts []ast.StmtNode) error {
	// 这里不需要判断len为0，form表单binding:"required"了
	if len(stmts) > 1 {
		return fmt.Errorf("每次仅允许执行一条语句，当前语句数量:%d", len(stmts))
	}
	return nil
}

func (s *ExecuteClickHouseQueryService) validateStatementType(stmt ast.StmtNode) error {
	// 判断语句类型是否被允许
	st := parser.StatementType{}
	statementType := st.Extract(stmt)
	type Result struct {
		Count int
	}
	var result Result
	global.App.DB.Table("insight_das_operations o").
		Select("count(*) as count").
		Where("name=? and is_enable=1", statementType).
		Scan(&result)
	if result.Count == 0 {
		return fmt.Errorf("禁止执行`%s`类型的语句（当前语句不被允许或不支持,请联系DBA）", statementType)
	}
	return nil
}

func (s *ExecuteClickHouseQueryService) Run() (ResponseData, error) {
	// 初始化SQLText为用户输入的SQLText
	var responseData ResponseData = ResponseData{SQLText: s.SQLText}
	var data *[]map[string]interface{}
	// 判断当前用户是否正在对当前实例当前库进行查询，如果有，禁止执行，防止并发
	if err := IsConcurrentRunning(s.Username, s.InstanceID, s.Schema); err != nil {
		return responseData, err
	}
	// 创建数据库记录
	jsonParams, err := json.Marshal(s.Params)
	if err != nil {
		return responseData, err
	}
	record := &models.InsightDASRecords{
		Username:  s.Username,
		Schema:    s.Schema,
		Sqltext:   s.SQLText,
		RequestID: requestid.Get(s.C),
		Params:    datatypes.JSON(jsonParams),
	}
	global.App.DB.Create(&record)
	// 解析UUID
	uuid, err := ParserUUID(s.InstanceID)
	if err != nil {
		return responseData, err
	}
	// 更新InstanceID
	global.App.DB.Model(&models.InsightDASRecords{}).
		Where("request_id=? and username=?", requestid.Get(s.C), s.Username).
		Update("InstanceID", uuid)
	// 根据不同的db类型调用不同的处理逻辑
	DbType, err := GetDbType(s.InstanceID)
	if err != nil {
		return responseData, err
	}
	if !strings.EqualFold(DbType, "clickhouse") {
		return responseData, fmt.Errorf("当前接口仅支持ClickHouse，当前DB类型为%s", DbType)
	}
	// 获取DB配置
	instance, err := GetDBConfig(s.InstanceID)
	if err != nil {
		return responseData, err
	}
	// 检查传入的SQL语法是否正确
	TiStmt, err := s.validateStatementSyntax()
	if err != nil {
		return responseData, err
	}
	// 每次仅允许执行一条SQL
	err = s.validateStatementNum(TiStmt.Stmts)
	if err != nil {
		return responseData, err
	}
	// 提取单条SQL的stmt
	var singleStmt ast.StmtNode = TiStmt.Stmts[0]
	// 生成指纹ID
	fingerId := query.Id(query.Fingerprint(singleStmt.Text()))
	responseData.QueryID = fingerId
	// 更新QueryID
	global.App.DB.Model(&models.InsightDASRecords{}).
		Where("request_id=? and username=?", requestid.Get(s.C), s.Username).
		Update("QueryID", fingerId)
	// 检查语句类型，判断语句类型是否被允许执行
	err = s.validateStatementType(singleStmt)
	if err != nil {
		return responseData, err
	}
	// 提取库表名
	extracter := parser.Extracter{Stmt: singleStmt, Schema: s.Schema}
	extractTables := extracter.Run()
	var tmpExtractTables []string
	for _, tt := range extractTables {
		tmpExtractTables = append(tmpExtractTables, strings.Join([]string{tt.Schema, tt.Table}, "."))
	}
	// 更新表名
	global.App.DB.Model(&models.InsightDASRecords{}).
		Where("request_id=? and username=?", requestid.Get(s.C), s.Username).
		Update("Tables", strings.Join(tmpExtractTables, ";"))
	// 检查库表权限
	checker := CheckUserPerm{
		UserName:   s.Username,
		InstanceID: uuid,
		Tables:     extractTables,
	}
	err = checker.HasSchemaPerms()
	if err != nil {
		return responseData, err
	}
	err = checker.HasTablePerms()
	if err != nil {
		return responseData, err
	}
	// 重写SQL，增加hint和重写limit
	rewrite := parser.Rewrite{Stmt: singleStmt, RequestID: requestid.Get(s.C), DbType: DbType}
	SQLText := rewrite.Run()
	s.SQLText = SQLText
	responseData.SQLText = s.SQLText
	// 更新rewrite sql
	global.App.DB.Model(&models.InsightDASRecords{}).
		Where("request_id=? and username=?", requestid.Get(s.C), s.Username).
		Update("RewriteSqltext", SQLText)
	// 调用clickhouse执行接口
	var executeApi ExecuteApi = ClickHouseExecuteApi{ExecuteClickHouseQueryForm: s.ExecuteClickHouseQueryForm, Ctx: s.C.Request.Context()}
	columns, data, duration, err := CalculateDuration(instance, executeApi.Execute)
	if err != nil {
		return responseData, err
	}
	// 更新耗时和影响行数
	global.App.DB.Model(&models.InsightDASRecords{}).
		Where("request_id=? and username=?", requestid.Get(s.C), s.Username).
		Updates(map[string]interface{}{"ReturnRows": len(*data), "Duration": duration})
	responseData.Duration = fmt.Sprintf("%dms", duration)
	responseData.Data = data
	responseData.Columns = columns
	return responseData, nil
}
