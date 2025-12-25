package parser

import (
	"errors"
	"fmt"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

type TiStmt struct {
	Stmts []ast.StmtNode
}

// 解析一条或多条SQL语句
func NewParse(sqltext, charset, collation string) ([]ast.StmtNode, error) {
	// tidb parser 语法解析
	stmts, warns, err := parser.New().Parse(sqltext, charset, collation)
	if len(warns) > 0 {
		return stmts, fmt.Errorf("Parse Warning: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return stmts, fmt.Errorf("SQL解析错误:%s", err.Error())
	}
	return stmts, nil
}

// 解析一条SQL语句
func NewParseOneStmt(sqltext, charset, collation string) (ast.StmtNode, error) {
	// tidb parser 语法解析
	stmt, err := parser.New().ParseOneStmt(sqltext, charset, collation)
	if err != nil {
		return stmt, fmt.Errorf("SQL解析错误:%s", err.Error())
	}
	return stmt, nil
}

// split
func SplitSQLText(sqltext string) (sqls []string, err error) {
	// 解析SQL
	stmts, err := NewParse(sqltext, "", "")
	if err != nil {
		return nil, err
	}
	for _, stmt := range stmts {
		sqls = append(sqls, stmt.OriginalText())
	}
	return sqls, nil
}

// 检查单次最大允许提交的SQL数量
func CheckMaxAllowedSQLNums(sqltext string) error {
	// 解析SQL
	stmts, err := NewParse(sqltext, "", "")
	if err != nil {
		return err
	}
	if len(stmts) > 2048 {
		return fmt.Errorf("单次最大允许提交2048条SQL，当前SQL语句条数为%d", len(stmts))
	}
	return nil
}

// 检查SQL类型
func CheckSqlType(sqltext, sqltype string) error {
	// 解析SQL
	stmts, err := NewParse(sqltext, "", "")
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		var st string
		switch stmt.(type) {
		case *ast.SelectStmt, *ast.SetOprStmt:
			st = "EXPORT"
		case *ast.DeleteStmt, *ast.InsertStmt, *ast.UpdateStmt:
			st = "DML"
		case *ast.AlterTableStmt, *ast.AlterSequenceStmt, *ast.AlterPlacementPolicyStmt:
			st = "DDL"
		case *ast.CreateDatabaseStmt, *ast.CreateIndexStmt, *ast.CreateTableStmt, *ast.CreateViewStmt, *ast.CreateSequenceStmt, *ast.CreatePlacementPolicyStmt:
			st = "DDL"
		case *ast.DropDatabaseStmt, *ast.DropIndexStmt, *ast.DropTableStmt, *ast.DropSequenceStmt, *ast.DropPlacementPolicyStmt:
			st = "DDL"
		case *ast.RenameTableStmt:
			st = "DDL"
		case *ast.TruncateTableStmt:
			st = "DDL"
		}
		if st != sqltype {
			if sqltype == "DML" {
				return fmt.Errorf("DML模式下，不允许提交%s语句", st)
			}
			if sqltype == "DDL" {
				return fmt.Errorf("DDL模式下，不允许提交%s语句", st)
			}
			if sqltype == "EXPORT" {
				return fmt.Errorf("EXPORT模式下，不允许提交%s语句，仅允许提交SELECT语句", st)
			}
		}
	}
	return nil
}

// 返回语句类型
func GetSqlStatement(sqltext string) (string, error) {
	const (
		CreateDatabase = "CreateDatabase"
		CreateTable    = "CreateTable"
		CreateView     = "CreateView"
		DropTable      = "DropTable"
		DropIndex      = "DropIndex"
		TruncateTable  = "TruncateTable"
		RenameTable    = "RenameTable"
		CreateIndex    = "CreateIndex"
		DropDatabase   = "DropDatabase"
		AlterTable     = "AlterTable"
	)

	stmt, err := NewParseOneStmt(sqltext, "", "")
	if err != nil {
		return "", err
	}
	switch stmt.(type) {
	case *ast.AlterTableStmt:
		return AlterTable, nil
	case *ast.CreateDatabaseStmt:
		return CreateDatabase, nil
	case *ast.CreateIndexStmt:
		return CreateIndex, nil
	case *ast.CreateTableStmt:
		return CreateTable, nil
	case *ast.CreateViewStmt:
		return CreateView, nil
	case *ast.DropIndexStmt:
		return DropIndex, nil
	case *ast.DropTableStmt:
		return DropTable, nil
	case *ast.RenameTableStmt:
		return RenameTable, nil
	case *ast.TruncateTableStmt:
		return TruncateTable, nil
	case *ast.DropDatabaseStmt:
		return DropDatabase, nil
	default:
		return "", errors.New("当前SQL未匹配到规则，执行失败")
	}
}

// 获取表名
func GetTableNameFromAlterStatement(sqltext string) (string, error) {
	stmt, err := NewParseOneStmt(sqltext, "", "")
	if err != nil {
		return "", err
	}
	switch s := stmt.(type) {
	case *ast.AlterTableStmt:
		return s.Table.Name.String(), nil
	}
	return "", errors.New("未提取到表名")
}
