/*
@Time    :   2022/08/25 16:32:17
@Author  :   xff
@Desc    :   None
*/

package rules

import (
	"goInsight/internal/inspect/controllers/logics"
	"goInsight/internal/inspect/controllers/traverses"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

func RenameTableRules() []Rule {
	return []Rule{
		{
			Hint:      "RenameTable#检查",
			CheckFunc: (*Rule).RuleRenameTable,
		},
	}
}

// RuleRenameTable
func (r *Rule) RuleRenameTable(tistmt *ast.StmtNode) {
	v := &traverses.TraverseRenameTable{}
	(*tistmt).Accept(v)
	logics.LogicRenameTable(v, r.RuleHint)
}
