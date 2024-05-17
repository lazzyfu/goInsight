/*
@Time    :   2022/06/24 13:12:20
@Author  :   xff
@Desc    :   遍历语法树
*/

package traverses

import (
	"github.com/pingcap/tidb/pkg/parser/ast"
)

// TraverseCreateViewIsExist
type TraverseCreateViewIsExist struct {
	View      string   // 视图名
	OrReplace bool     // 是否为replace语句
	Tables    []string // 表名
}

func (c *TraverseCreateViewIsExist) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateViewStmt); ok {
		c.OrReplace = stmt.OrReplace
		c.View = stmt.ViewName.Name.String()
	}
	return in, false
}

func (c *TraverseCreateViewIsExist) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
