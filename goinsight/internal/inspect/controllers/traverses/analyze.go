/*
@Time    :   2023/04/19 15:11:14
@Author  :   xff
@Desc    :
*/

package traverses

import (
	"github.com/pingcap/tidb/pkg/parser/ast"
)

// TraverseAnalyzeTable
type TraverseAnalyzeTable struct {
	IsMatch    int
	TableNames []string // 表名
}

func (c *TraverseAnalyzeTable) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AnalyzeTableStmt); ok {
		c.IsMatch++
		for _, t := range stmt.TableNames {
			c.TableNames = append(c.TableNames, t.Name.L)
		}
	}
	return in, false
}

func (c *TraverseAnalyzeTable) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
