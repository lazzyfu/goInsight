/*
@Time    :   2022/07/06 10:12:58
@Author  :   xff
@Desc    :   None
*/

package rules

import (
	"goInsight/internal/inspect/controllers/logics"
	"goInsight/internal/inspect/controllers/traverses"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

func DropTableRules() []Rule {
	return []Rule{
		{
			Hint:      "DropTable#检查",
			CheckFunc: (*Rule).RuleDropTable,
		},
		{
			Hint:      "TruncateTable#检查",
			CheckFunc: (*Rule).RuleTruncateTable,
		},
	}
}

// RuleDropTable
func (r *Rule) RuleDropTable(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDropTable{}
	(*tistmt).Accept(v)
	logics.LogicDropTable(v, r.RuleHint)
}

// RuleTruncateTable
func (r *Rule) RuleTruncateTable(tistmt *ast.StmtNode) {
	v := &traverses.TraverseTruncateTable{}
	(*tistmt).Accept(v)
	logics.LogicTruncateTable(v, r.RuleHint)
}
