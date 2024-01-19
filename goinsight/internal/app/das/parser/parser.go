/*
@Time    :   2023/03/21 14:53:21
@Author  :   zongfei.fu
@Desc    :   语法解析
go get -v github.com/pingcap/tidb/parser@4084b07
go get -v github.com/pingcap/tidb/types/parser_driver@4084b07
*/

package parser

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
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
