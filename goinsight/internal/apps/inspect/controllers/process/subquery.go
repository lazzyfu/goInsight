/*
@Time    :   2022/07/15 14:16:21
@Author  :   zongfei.fu
@Desc    :   迭代子查询
*/

package process

import (
	"github.com/pingcap/tidb/parser/ast"
)

// 迭代子查询
type CheckItem struct {
	HasSubQuery bool
	HasOrderBy  bool
	HasLimit    bool
}

func (c *CheckItem) CheckExprItem(expr ast.ExprNode) {
	switch e := expr.(type) {
	case *ast.PatternInExpr:
		c.CheckExprItem(e.Sel)
	case *ast.CompareSubqueryExpr:
		c.CheckExprItem(e.R)
	case *ast.BinaryOperationExpr:
		c.CheckExprItem(e.L)
		c.CheckExprItem(e.R)
	case *ast.ExistsSubqueryExpr:
		c.CheckExprItem(e.Sel)
	case *ast.SubqueryExpr:
		if !c.HasSubQuery {
			c.HasSubQuery = true
		}
		c.CheckSelectItem(e.Query)
	}
}

func (c *CheckItem) CheckSelectItem(node ast.ResultSetNode) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.SelectStmt:
		c.CheckSubSelectItem(n)
	case *ast.Join:
		c.CheckSelectItem(n.Left)
		c.CheckSelectItem(n.Right)
	case *ast.TableSource:
		c.CheckSelectItem(n.Source)
	}
}

func (c *CheckItem) CheckSubSelectItem(node *ast.SelectStmt) {
	if node.From != nil {
		// 迭代from子查询
		c.CheckSelectItem(node.From.TableRefs)
	}
	if node.Where != nil {
		c.CheckExprItem(node.Where)
	}
	if node.OrderBy != nil {
		if !c.HasOrderBy {
			c.HasOrderBy = true
		}
	}
	if node.Limit != nil {
		if !c.HasLimit {
			c.HasLimit = true
		}
	}
	for _, item := range node.Fields.Fields {
		if item.Expr != nil {
			c.CheckExprItem(item.Expr)
		}
	}
}
