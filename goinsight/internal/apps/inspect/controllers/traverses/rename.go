/*
@Time    :   2022/08/25 16:41:19
@Author  :   xff
@Desc    :   None
*/

package traverses

import (
	"github.com/pingcap/tidb/pkg/parser/ast"
)

type RenameTable struct {
	OldTable string // 表名
	NewTable string // 是否匹配当前规则
}

// TraverseRenameTable
type TraverseRenameTable struct {
	IsMatch int
	Tables  []RenameTable
}

func (c *TraverseRenameTable) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.RenameTableStmt); ok {
		c.IsMatch++
		for _, t := range stmt.TableToTables {
			c.Tables = append(c.Tables, RenameTable{
				OldTable: t.OldTable.Name.String(),
				NewTable: t.NewTable.Name.String(),
			})
		}
	}
	return in, false
}

func (c *TraverseRenameTable) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
