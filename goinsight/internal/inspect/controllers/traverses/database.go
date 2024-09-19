package traverses

import (
	"github.com/pingcap/tidb/pkg/parser/ast"
)

// TraverseCreateViewIsExist
type TraverseCreateDatabaseIsExist struct {
	Name string // 库名
}

func (c *TraverseCreateDatabaseIsExist) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateDatabaseStmt); ok {
		c.Name = stmt.Name.String()
	}
	return in, false
}

func (c *TraverseCreateDatabaseIsExist) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
