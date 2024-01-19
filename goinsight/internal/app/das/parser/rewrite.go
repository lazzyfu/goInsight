/*
@Time    :   2023/03/24 10:06:56
@Author  :   zongfei.fu
@Desc    :   重写sql
*/

package parser

import (
	"fmt"
	"goInsight/global"
	"strings"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	driver "github.com/pingcap/tidb/types/parser_driver"
	utilparser "github.com/pingcap/tidb/util/parser"
)

type Rewrite struct {
	Stmt      ast.StmtNode
	RequestID string
	DbType    string
}

func (r *Rewrite) BindHints(name string, value interface{}) (hints []*ast.TableOptimizerHint) {
	// 绑定hints
	var hint *ast.TableOptimizerHint = &ast.TableOptimizerHint{
		HintName: model.CIStr{O: name, L: name},
		HintData: value,
	}
	hints = append(hints, hint)
	return hints
}

func (r *Rewrite) AddHintsForMaxExecutionTime() {
	// 增加max_execution_time
	switch stmt := r.Stmt.(type) {
	case *ast.SelectStmt:
		var maxTime uint64 = global.App.Config.Das.MaxExecutionTime
		hints := r.BindHints("max_execution_time", maxTime)
		hints = append(hints, stmt.TableHints...)
		stmt.TableHints = hints
		r.Stmt = stmt
	}
}

func (r *Rewrite) RewriteLimitSetCount(stmt ast.StmtNode, value uint64) {
	// 重写count
	switch stmt := r.Stmt.(type) {
	case *ast.SelectStmt:
		if stmt.Limit == nil {
			stmt.Limit = &ast.Limit{}
			stmt.Limit.Count = &driver.ValueExpr{}
		}
		switch ex := stmt.Limit.Count.(type) {
		case *driver.ValueExpr:
			ex.SetValue(value)
		}
	case *ast.SetOprStmt:
		if stmt.Limit == nil {
			stmt.Limit = &ast.Limit{}
			stmt.Limit.Count = &driver.ValueExpr{}
		}
		switch ex := stmt.Limit.Count.(type) {
		case *driver.ValueExpr:
			ex.SetValue(value)
		}
	}
}

// 暂时用不到
// func (r *Rewrite) RewriteLimitSetOffset(stmt *ast.SelectStmt, value uint64) {
// 	// 重写offset
// 	if stmt.Limit == nil {
// 		// 如果SQL没有传递limit，此时Limit会为nil
// 		stmt.Limit = &ast.Limit{}
// 		stmt.Limit.Offset = &driver.ValueExpr{}
// 	}
// 	switch ex := stmt.Limit.Offset.(type) {
// 	case *driver.ValueExpr:
// 		ex.SetValue(value)
// 	}
// }

func (r *Rewrite) RewriteLimit() {
	// 重写limit
	switch stmt := r.Stmt.(type) {
	case *ast.SelectStmt, *ast.SetOprStmt:
		// 遍历
		v := &Limit{}
		(r.Stmt).Accept(v)
		// SQL语句没有limit子句，增加limit N
		if v.Count == 0 {
			r.RewriteLimitSetCount(r.Stmt, global.App.Config.Das.DefaultReturnRows)
		}
		// SQL语句有limit N子句
		if v.Count != 0 {
			// 当N大于定义的MaxReturnRows时，改写为limit MaxReturnRows
			if v.Count > global.App.Config.Das.DefaultReturnRows {
				r.RewriteLimitSetCount(r.Stmt, global.App.Config.Das.MaxReturnRows)
			}
		}
		// SQL语句有LIMIT N, N 或 LIMIT N OFFSET N子句
		if v.Count != 0 && v.Offset != 0 {
			// 当N大于定义的MaxReturnRows时，改写为limit MaxReturnRows
			if v.Count > global.App.Config.Das.DefaultReturnRows {
				r.RewriteLimitSetCount(r.Stmt, global.App.Config.Das.MaxReturnRows)
			}
		}
		r.Stmt = stmt
	}
}

func (r *Rewrite) RewriteExplain() {
	// explain
	switch stmt := r.Stmt.(type) {
	case *ast.ExplainStmt:
		switch stmt.Stmt.(type) {
		case *ast.SelectStmt, *ast.SetOprStmt:
			// mysql没有row格式，仅有traditional，json格式
			if strings.EqualFold(r.DbType, "mysql") && strings.EqualFold(stmt.Format, "row") {
				stmt.Format = "traditional"
			}
		}
		r.Stmt = stmt
	}
}

func (r *Rewrite) ReplaceClickHouseExplain(sql string) string {
	// 如果是clickhouse explain，移除format子句
	switch stmt := r.Stmt.(type) {
	case *ast.ExplainStmt:
		switch stmt.Stmt.(type) {
		case *ast.SelectStmt, *ast.SetOprStmt:
			if strings.EqualFold(r.DbType, "clickhouse") {
				return strings.Replace(sql, "FORMAT = 'row'", "", 1)
			}
		}
	}
	return sql
}

func (r *Rewrite) RestoreSQL() string {
	// 从ast还原SQL
	return utilparser.RestoreWithDefaultDB(r.Stmt, "", "")
}

func (r *Rewrite) AddCommentForRequestID(sql string) string {
	// SQL增加request_id注释
	return strings.Join([]string{sql, fmt.Sprintf("/* %s */", r.RequestID)}, " ")
}

func (r *Rewrite) Run() string {
	r.AddHintsForMaxExecutionTime()
	r.RewriteLimit()
	r.RewriteExplain()
	restoreSQL := r.RestoreSQL()
	replaceRestoreSQL := r.ReplaceClickHouseExplain(restoreSQL)
	return r.AddCommentForRequestID(replaceRestoreSQL)
}
