/*
@Time    :   2022/07/06 10:12:48
@Author  :   zongfei.fu
@Desc    :   None
*/

package rules

import (
	// "goInsight/internal/app/inspect/controllers/extract"
	"goInsight/internal/app/inspect/controllers/logics"
	"goInsight/internal/app/inspect/controllers/traverses"

	"github.com/pingcap/tidb/parser/ast"
)

func DMLRules() []Rule {
	return []Rule{
		// {
		// 	Hint:      "DML#限制部分表进行语法审核",
		// 	CheckFunc: (*Rule).RuleDisableAuditDMLTables,
		// },
		{
			Hint:      "DML#是否允许INSERT INTO SELECT语法",
			CheckFunc: (*Rule).RuleDMLInsertIntoSelect,
		},
		{
			Hint:      "DML#必须要有WHERE条件",
			CheckFunc: (*Rule).RuleDMLNoWhere,
		},
		{
			Hint:      "DML#INSERT必须指定列名",
			CheckFunc: (*Rule).RuleDMLInsertWithColumns,
		},
		{
			Hint:      "DML#不能有LIMIT/ORDERBY/SubQuery",
			CheckFunc: (*Rule).RuleDMLHasConstraint,
		},
		{
			Hint:      "DML#JOIN操作必须要有ON语句",
			CheckFunc: (*Rule).RuleDMLJoinWithOn,
		},
		{
			Hint:      "DML#更新影响行数",
			CheckFunc: (*Rule).RuleDMLMaxUpdateRows,
		},
		{
			Hint:      "DML#插入影响行数",
			CheckFunc: (*Rule).RuleDMLMaxInsertRows,
		},
	}
}

// RuleDisableAuditDMLTables
// func (r *Rule) RuleDisableAuditDMLTables(tistmt *ast.StmtNode) {
// 	v := &traverses.TraverseDisableAuditDMLTables{}
// 	v.Tables, _ = extract.ExtractTablesFromStatement(tistmt)
// 	logics.LogicDisableAuditDMLTables(v, r.RuleHint)
// }

// RuleDMLInsertIntoSelect
func (r *Rule) RuleDMLInsertIntoSelect(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDMLInsertIntoSelect{}
	(*tistmt).Accept(v)
	logics.LogicDMLInsertIntoSelect(v, r.RuleHint)
}

// RuleDMLNoWhere
func (r *Rule) RuleDMLNoWhere(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDMLNoWhere{}
	(*tistmt).Accept(v)
	logics.LogicDMLNoWhere(v, r.RuleHint)
}

// RuleDMLInsertWithColumns
func (r *Rule) RuleDMLInsertWithColumns(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDMLInsertWithColumns{}
	(*tistmt).Accept(v)
	logics.LogicDMLInsertWithColumns(v, r.RuleHint)
}

// RuleDMLHasConstraint
func (r *Rule) RuleDMLHasConstraint(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDMLHasConstraint{}
	(*tistmt).Accept(v)
	logics.LogicDMLHasConstraint(v, r.RuleHint)
}

// RuleDMLJoinWithOn
func (r *Rule) RuleDMLJoinWithOn(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDMLJoinWithOn{}
	(*tistmt).Accept(v)
	logics.LogicDMLJoinWithOn(v, r.RuleHint)
}

// RuleDMLMaxUpdateRows
func (r *Rule) RuleDMLMaxUpdateRows(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDMLMaxUpdateRows{}
	(*tistmt).Accept(v)
	logics.LogicDMLMaxUpdateRows(v, r.RuleHint)
}

// RuleDMLMaxInsertRows
func (r *Rule) RuleDMLMaxInsertRows(tistmt *ast.StmtNode) {
	v := &traverses.TraverseDMLMaxInsertRows{}
	(*tistmt).Accept(v)
	logics.LogicDMLMaxInsertRows(v, r.RuleHint)
}
