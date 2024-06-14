/*
@Time    :   2022/07/06 10:12:33
@Author  :   xff
@Desc    :   None
*/

package parser

import (
	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

type Audit struct {
	Query  string
	TiStmt []ast.StmtNode // 通过TiDB解析出的抽象语法树
}

// NewParse
func NewParse(sqltext, charset, collation string) (*Audit, []error, error) {
	q := &Audit{Query: sqltext}

	// tidb parser 语法解析
	var warns []error
	var err error
	q.TiStmt, warns, err = parser.New().Parse(sqltext, charset, collation)
	return q, warns, err
}
