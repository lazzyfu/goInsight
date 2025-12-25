package rules

import (
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/logics"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

func CreateDatabaseRules() []Rule {
	return []Rule{
		{
			Hint:      "CreateDatabase#检查DB是否存在",
			CheckFunc: (*Rule).RuleCreateDatabaseIsExist,
		},
	}
}

// RuleCreateDatabaseIsExist
func (r *Rule) RuleCreateDatabaseIsExist(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateDatabaseIsExist{}
	(*tistmt).Accept(v)
	logics.LogicCreateDatabaseIsExist(v, r.RuleHint)
}
