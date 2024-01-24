package services

import (
	"context"
	"fmt"
	"goInsight/global"
	"goInsight/internal/app/inspect/controllers/dao"
	"goInsight/internal/app/inspect/controllers/parser"
	"goInsight/internal/app/inspect/forms"
	"goInsight/internal/pkg/kv"
	"goInsight/internal/pkg/query"
	"goInsight/internal/pkg/utils"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

// 返回数据
type ReturnData struct {
	Summary      []string `json:"summary"` // 规则摘要
	Level        string   `json:"level"`   // 提醒级别,INFO/WARN/ERROR
	AffectedRows int      `json:"affected_rows"`
	Type         string   `json:"type"`
	FingerId     string   `json:"finger_id"`
	Query        string   `json:"query"` // 原始SQL
}

// 语法check
type SyntaxInspectService struct {
	Form      forms.SyntaxInspectForm
	C         *gin.Context
	Username  string
	DB        *utils.DB
	Charset   string
	Collation string
	Audit     *parser.Audit
}

// 初始化DB
func (s *SyntaxInspectService) initDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.Form.Timeout)*time.Millisecond)
	defer cancel()
	s.DB = &utils.DB{
		User:     s.Form.DbUser,
		Password: s.Form.DbPassword,
		Host:     s.Form.DbHost,
		Port:     s.Form.DbPort,
		Database: s.Form.DB,
		Ctx:      ctx,
	}
}

func (s *SyntaxInspectService) parser() error {
	// 解析SQL
	var warns []error
	var err error
	// 处理审核参数
	// if err := s.CustomParams(); err != nil {
	// 	return err
	// }
	// 解析
	s.Audit, warns, err = parser.NewParse(s.Form.SqlText, s.Charset, s.Collation)
	if len(warns) > 0 {
		return fmt.Errorf("Parse Warning: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return fmt.Errorf("sql解析错误：%s", err.Error())
	}
	return nil
}

func (s *SyntaxInspectService) Run() (err error, returnData []ReturnData) {
	// 初始化DB
	s.initDB()
	// RequestID
	requestID := requestid.Get(s.C)

	// 存放alter语句中的表名
	var mergeAlters []string

	// 每次请求基于RequestID初始化kv cache
	kv := kv.NewKVCache(requestID)

	// 获取目标数据库变量
	dbVars, err := dao.GetDBVars(s.DB)
	if err != nil {
		errMsg := fmt.Sprintf("获取DB变量失败：%s", err.Error())
		global.App.Log.Error(errMsg)
		return fmt.Errorf(errMsg), returnData
	}
	for k, v := range dbVars {
		kv.Put(k, v)
	}
	s.Charset = dbVars["dbCharset"]

	// 解析SQL
	err = s.parser()
	if err != nil {
		global.App.Log.Error(err)
		return err, returnData
	}

	// 迭代stmt
	for _, stmt := range s.Audit.TiStmt {
		// 移除SQL尾部的分号
		sqlTrim := strings.TrimSuffix(stmt.Text(), ";")

		// 生成指纹ID
		fingerId := query.Id(query.Fingerprint(sqlTrim))

		// 存储指纹ID
		kv.Put(fingerId, true)

		// 迭代
		st := Stmt{s}
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
		default:
			// 不允许的其他语句，有需求可以扩展
			var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "", Level: "WARN"}
			data.Summary = append(data.Summary, "不被允许的审核语句，请联系数据库管理员")
			returnData = append(returnData, data)
		}
	}
}
