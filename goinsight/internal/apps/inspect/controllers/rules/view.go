/*
@Time    :   2022/07/06 10:12:42
@Author  :   xff
@Desc    :   None
*/

package rules

import (
	"goInsight/internal/apps/inspect/controllers/extract"
	"goInsight/internal/apps/inspect/controllers/logics"
	"goInsight/internal/apps/inspect/controllers/traverses"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

func CreateViewRules() []Rule {
	return []Rule{
		{
			Hint:      "CreateView#检查视图是否存在",
			CheckFunc: (*Rule).RuleCreateViewIsExist,
		},
	}
}

// RuleCreateViewIsExist
func (r *Rule) RuleCreateViewIsExist(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateViewIsExist{}
	(*tistmt).Accept(v)
	v.Tables, _ = extract.ExtractTablesFromStatement(tistmt)
	logics.LogicCreateViewIsExist(v, r.RuleHint)
}
