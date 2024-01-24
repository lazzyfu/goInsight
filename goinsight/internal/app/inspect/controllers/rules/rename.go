/*
@Time    :   2022/08/25 16:32:17
@Author  :   zongfei.fu
@Desc    :   None
*/

package rules

import (
	"goInsight/internal/app/inspect/controllers/logics"
	"goInsight/internal/app/inspect/controllers/traverses"

	"github.com/pingcap/tidb/parser/ast"
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
