package parser

import (
	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

type TiStmt struct {
	Stmts []ast.StmtNode
}

// NewParse
func NewParse(sqltext, charset, collation string) (*TiStmt, []error, error) {
	t := &TiStmt{}
	// tidb parser 语法解析
	var warns []error
	var err error
	t.Stmts, warns, err = parser.New().Parse(sqltext, charset, collation)
	return t, warns, err
}
