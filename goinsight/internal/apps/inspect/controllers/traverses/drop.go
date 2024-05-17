/*
@Time    :   2022/06/24 13:12:20
@Author  :   xff
@Desc    :   遍历语法树
*/

package traverses

import (
	"github.com/pingcap/tidb/pkg/parser/ast"
)

// TraverseDropTable
type TraverseDropTable struct {
	Tables         []string // 表名
	IsMatch        int      // 是否匹配当前规则
	IsHasDropTable bool
}

func (c *TraverseDropTable) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.DropTableStmt); ok {
		c.IsMatch++
		for _, table := range stmt.Tables {
			c.Tables = append(c.Tables, table.Name.O)
		}
		c.IsHasDropTable = true
	}
	return in, false
}

func (c *TraverseDropTable) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseTruncateTable
type TraverseTruncateTable struct {
	Table              string // 表名
	IsMatch            int    // 是否匹配当前规则
	IsHasTruncateTable bool
}

func (c *TraverseTruncateTable) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.TruncateTableStmt); ok {
		c.IsMatch++
		c.Table = stmt.Table.Name.O
		c.IsHasTruncateTable = true
	}
	return in, false
}

func (c *TraverseTruncateTable) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
