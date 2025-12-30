package logics

import (
	"fmt"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/parser"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/process"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicDMLInsertIntoSelect
func LogicDMLInsertIntoSelect(v *traverses.TraverseDMLInsertIntoSelect, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if r.InspectParams.DISABLE_INSERT_INTO_SELECT && v.HasSelectSubQuery {
		r.Warn(fmt.Sprintf("禁止使用 `%s INTO SELECT` 语法", v.DMLType))
		r.IsBreak = true
	}
	if r.InspectParams.DISABLE_ON_DUPLICATE && v.HasOnDuplicate {
		r.Warn(fmt.Sprintf("禁止使用 `%s ... ON DUPLICATE KEY UPDATE` 语法", v.DMLType))
		r.IsBreak = true
	}
}

// LogicDMLNoWhere
func LogicDMLNoWhere(v *traverses.TraverseDMLNoWhere, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if !v.HasWhere && r.InspectParams.DML_MUST_HAVE_WHERE {
		r.Warn(fmt.Sprintf("`%s` 语句必须包含 `WHERE` 条件（避免全表扫描/全表更新）", v.DMLType))
		r.IsBreak = true
	}
}

// LogicDMLInsertWithColumns
func LogicDMLInsertWithColumns(v *traverses.TraverseDMLInsertWithColumns, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if v.DMLType == "REPLACE" && r.InspectParams.DISABLE_REPLACE {
		r.Warn("禁止使用 `REPLACE` 语句（可能导致隐式删除+插入）")
		r.IsBreak = true
		return
	}
	// 获取db表结构
	audit, err := dao.ShowCreateTable(v.Table, r.DB, r.KV)
	if err != nil {
		r.Warn(err.Error())
		return
	}
	// 解析获取的db表结构
	vAudit := &traverses.TraverseAlterTableShowCreateTableGetCols{}
	switch audit := audit.(type) {
	case *parser.Audit:
		(audit.TiStmt[0]).Accept(vAudit)
	}
	// 校验列是否存在（避免拼写错误/脏字段导致执行失败）。
	for _, col := range v.Columns {
		if !utils.IsContain(vAudit.Cols, col) {
			r.Warn(fmt.Sprintf("表`%s`中列`%s`不存在", v.Table, col))
		}
	}
	// INSERT/REPLACE 建议显式指定列名：避免表结构变更后出现列错位。
	if v.ColumnsCount == 0 {
		r.Warn(fmt.Sprintf("`%s` 语句必须显式指定列名（如：`%s INTO t(col1,col2) VALUES ...`）", v.DMLType, v.DMLType))
	} else if !v.ColsValuesIsMatch {
		r.Warn(fmt.Sprintf("`%s` 语句列数量与 VALUES 值数量不匹配", v.DMLType))
	}
	if v.RowsCount > r.InspectParams.MAX_INSERT_ROWS {
		r.Warn(fmt.Sprintf("`%s` 单次写入行数过多：最多允许 %d 行，当前 %d 行；建议拆分为多条执行", v.DMLType, r.InspectParams.MAX_INSERT_ROWS, v.RowsCount))
	}
}

// LogicDMLHasLimit
func LogicDMLHasConstraint(v *traverses.TraverseDMLHasConstraint, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if v.HasLimit && r.InspectParams.DML_DISABLE_LIMIT {
		r.Warn(fmt.Sprintf("`%s` 语句禁止使用 `LIMIT` 子句", v.DMLType))
		r.IsBreak = true
	}
	if v.HasOrderBy && r.InspectParams.DML_DISABLE_ORDERBY {
		r.Warn(fmt.Sprintf("`%s` 语句禁止使用 `ORDER BY` 子句", v.DMLType))
		r.IsBreak = true
	}
	if v.HasSubQuery && r.InspectParams.DML_DISABLE_SUBQUERY {
		r.Warn(fmt.Sprintf("`%s` 语句禁止包含子查询（请改写为 JOIN 或分步执行）", v.DMLType))
		r.IsBreak = true
	}
}

// LogicDMLJoinWithOn
func LogicDMLJoinWithOn(v *traverses.TraverseDMLJoinWithOn, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if v.HasJoin && r.InspectParams.CHECK_DML_JOIN_WITH_ON && !v.IsJoinWithOn {
		r.Warn(fmt.Sprintf("`%s` 的 JOIN 必须包含 `ON` 条件（禁止笛卡尔积）", v.DMLType))
		r.IsBreak = true
	}
}

// LogicDMLMaxUpdateRows
func LogicDMLMaxUpdateRows(v *traverses.TraverseDMLMaxUpdateRows, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	explain := process.Explain{DB: r.DB, SQL: r.Query, KV: r.KV}
	affectedRows, err := explain.Get(r.InspectParams.EXPLAIN_RULE)
	if err != nil {
		r.AffectedRows = 0
		r.Warn(err.Error())
		r.IsBreak = true
		return
	}
	if affectedRows > r.InspectParams.MAX_AFFECTED_ROWS {
		r.AffectedRows = affectedRows
		r.Warn(fmt.Sprintf("`%s` 预计影响/扫描行数 %d，超过上限 %d；建议拆分语句或缩小 WHERE 范围", v.DMLType, affectedRows, r.InspectParams.MAX_AFFECTED_ROWS))
		r.IsBreak = true
		return
	}
	r.IsBreak = true
	r.AffectedRows = affectedRows
}

// LogicDMLMaxInsertRows
func LogicDMLMaxInsertRows(v *traverses.TraverseDMLMaxInsertRows, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	r.AffectedRows = v.RowsCount
	r.IsBreak = true
}
