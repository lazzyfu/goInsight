/*
@Time    :   2022/06/29 15:30:31
@Author  :   xff
*/

package rules

import (
	"goInsight/internal/apps/inspect/controllers"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

type Rule struct {
	*controllers.RuleHint
	Hint      string                     `json:"hint"` // 规则说明
	CheckFunc func(*Rule, *ast.StmtNode) // 函数名
}
