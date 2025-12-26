package checker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/kv"
	"github.com/lazzyfu/goinsight/pkg/query"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/inspect/config"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/parser"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/process"
	"github.com/lazzyfu/goinsight/internal/inspect/models"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
	"gorm.io/datatypes"
)

// 返回数据
type ReturnData struct {
	Summary      []string `json:"summary"`       // 摘要
	Level        string   `json:"level"`         // 级别,INFO/WARN/ERROR
	AffectedRows int      `json:"affected_rows"` // 影响行数
	Type         string   `json:"type"`          // SQL类型
	FingerId     string   `json:"finger_id"`     // 指纹
	Query        string   `json:"query"`         // 原始SQL
}

// 语法check
type SyntaxInspectService struct {
	C             *gin.Context
	DbUser        string
	DbPassword    string
	DbHost        string
	DbPort        int
	DBSchema      string
	DBParams      datatypes.JSON
	SqlText       string
	Username      string
	Charset       string
	Collation     string
	DB            *dao.DB
	Audit         *parser.Audit
	InspectParams config.InspectParams
}

// 初始化DB
func (s *SyntaxInspectService) initDB() {
	s.DB = &dao.DB{
		User:     s.DbUser,
		Password: s.DbPassword,
		Host:     s.DbHost,
		Port:     s.DbPort,
		Database: s.DBSchema,
	}
}

// 初始化默认审核参数
func (s *SyntaxInspectService) initDefaultInspectParams() error {
	// 读取数据库参数
	var rows []models.InsightInspectParams
	tx := global.App.DB.Model(&models.InsightInspectParams{}).Scan(&rows)
	if tx.RowsAffected == 0 {
		return errors.New("获取审核参数失败，表insight_inspect_params未找到记录")
	}
	// 初始化map存储参数
	jsonParams := make(map[string]json.RawMessage)
	for _, row := range rows {
		err := json.Unmarshal(row.Params, &jsonParams)
		if err != nil {
			return fmt.Errorf("解析JSON参数失败: %v，错误参数：%v", err, row)
		}
	}
	// 序列化参数
	jsonData, err := json.Marshal(jsonParams)
	if err != nil {
		return fmt.Errorf("序列化JSON参数失败: %v", err)
	}
	// 转换为结构体
	var ips config.InspectParams
	err = json.Unmarshal(jsonData, &ips)
	if err != nil {
		return fmt.Errorf("反序列化JSON参数失败: %v", err)
	}
	s.InspectParams = ips
	return nil
}

// 初始化DB审核参数
func (s *SyntaxInspectService) initDBInspectParams() error {
	// 序列化参数
	jsonParams, err := json.Marshal(s.DBParams)
	if err != nil {
		return fmt.Errorf("序列化JSON参数失败: %v", err)
	}
	r := bytes.NewReader(jsonParams)
	decoder := json.NewDecoder(r)
	// 动态参数赋值给默认模板
	// 优先级: post instance params > 内置默认参数
	if err := decoder.Decode(&s.InspectParams); err != nil {
		return err
	}
	return nil
}

func (s *SyntaxInspectService) parser() error {
	// 解析SQL
	var warns []error
	var err error
	// 解析
	s.Audit, warns, err = parser.NewParse(s.SqlText, s.Charset, s.Collation)
	if len(warns) > 0 {
		return fmt.Errorf("Parse Warning: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return fmt.Errorf("sql解析错误：%w", err)
	}
	return nil
}

// 判断多条alter语句是否需要合并
func (s *SyntaxInspectService) mergeAlters(kv *kv.KVCache, mergeAlters []string) ReturnData {
	data := ReturnData{FingerId: utils.GenerateSimpleRandomString(16), Level: LevelInfo}
	dbVersionIns := process.DbVersion{Version: kv.Get("dbVersion").(string)}
	if s.InspectParams.ENABLE_MYSQL_MERGE_ALTER_TABLE && !dbVersionIns.IsTiDB() {
		if ok, val := utils.IsRepeat(mergeAlters); ok {
			for _, v := range val {
				data.Summary = append(data.Summary, fmt.Sprintf("[MySQL数据库]表`%s`的多条ALTER操作，请合并为一条ALTER语句", v))
			}
		}
	}
	if len(data.Summary) > 0 {
		data.Level = LevelWarn
	}
	return data
}

func (s *SyntaxInspectService) Run() (returnData []ReturnData, err error) {
	// 初始化系统默认审核参数
	err = s.initDefaultInspectParams()
	if err != nil {
		return nil, err
	}
	// 获取DB定义的审核参数，优先级>系统默认审核参数
	err = s.initDBInspectParams()
	if err != nil {
		return nil, err
	}
	// 初始化DB
	s.initDB()
	// RequestID
	requestID := requestid.Get(s.C)
	// 存放alter语句中的表名
	var mergeAlters []string
	// 每次请求基于RequestID初始化kv cache
	kv := kv.NewKVCache(requestID)
	defer kv.Clear(requestID)
	// 获取目标数据库变量
	dbVars, err := dao.GetDBVars(s.DB)
	if err != nil {
		return returnData, fmt.Errorf("获取DB变量失败：%s", err.Error())
	}
	for k, v := range dbVars {
		kv.Put(k, v)
	}
	s.Charset = dbVars["dbCharset"]
	// 解析SQL
	err = s.parser()
	if err != nil {
		return returnData, err
	}
	// 迭代stmt
	st := Stmt{s}
	for _, stmt := range s.Audit.TiStmt {
		// 移除SQL尾部的分号
		sqlTrim := strings.TrimSuffix(stmt.Text(), ";")
		// 生成指纹ID
		fingerId := query.Id(query.Fingerprint(sqlTrim))
		// 存储指纹ID
		kv.Put(fingerId, true)
		switch stmt.(type) {
		case *ast.SelectStmt:
			// select语句不允许审核
			var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DML", Level: "WARN"}
			data.Summary = append(data.Summary, "发现SELECT语句，请删除SELECT语句后重新审核")
			returnData = append(returnData, data)
		case *ast.CreateTableStmt:
			returnData = append(returnData, st.CreateTableStmt(stmt, kv, fingerId))
		case *ast.CreateViewStmt:
			returnData = append(returnData, st.CreateViewStmt(stmt, kv, fingerId))
		case *ast.AlterTableStmt:
			data, mergeAlter := st.AlterTableStmt(stmt, kv, fingerId)
			mergeAlters = append(mergeAlters, mergeAlter)
			returnData = append(returnData, data)
		case *ast.DropTableStmt, *ast.TruncateTableStmt:
			returnData = append(returnData, st.DropTableStmt(stmt, kv, fingerId))
		case *ast.DeleteStmt, *ast.InsertStmt, *ast.UpdateStmt:
			returnData = append(returnData, st.DMLStmt(stmt, kv, fingerId))
		case *ast.RenameTableStmt:
			returnData = append(returnData, st.RenameTableStmt(stmt, kv, fingerId))
		case *ast.AnalyzeTableStmt:
			returnData = append(returnData, st.AnalyzeTableStmt(stmt, kv, fingerId))
		case *ast.CreateDatabaseStmt:
			returnData = append(returnData, st.CreateDatabaseStmt(stmt, kv, fingerId))
		default:
			// 不允许的其他语句，有需求可以扩展
			var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "", Level: "WARN"}
			data.Summary = append(data.Summary, "未识别或禁止的审核语句，请联系数据库管理员")
			returnData = append(returnData, data)
		}
	}
	// 判断多条alter语句是否需要合并
	if len(mergeAlters) > 1 {
		mergeData := s.mergeAlters(kv, mergeAlters)
		if len(mergeData.Summary) > 0 {
			returnData = append(returnData, mergeData)
		}
	}
	// 比如只传递了注释,如:#
	if len(s.Audit.TiStmt) == 0 {
		return
	}
	return
}
