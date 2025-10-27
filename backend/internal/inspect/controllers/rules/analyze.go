package rules

import (
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/logics"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

func AnalyzeTableRules() []Rule {
	return []Rule{
		{
			Hint:      "AnalyzeTable#检查",
			CheckFunc: (*Rule).RuleAnalyzeTable,
		},
	}
}

// RuleAnalyzeTable
func (r *Rule) RuleAnalyzeTable(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAnalyzeTable{}
	(*tistmt).Accept(v)
	logics.LogicAnalyzeTable(v, r.RuleHint)
}
