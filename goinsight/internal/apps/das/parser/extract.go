/*
@Time    :   2022/09/09 10:35:02
@Desc    :   提取SQL语句的库表
*/

package parser

import (
	"strings"

	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

// 返回数据
type Table struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
}
type ReturnData struct {
	Tables []Table `json:"tables"` // 库表名
}

// 提前
type Extracter struct {
	Schema string
	Stmt   ast.StmtNode
}

func (e *Extracter) Run() []Table {
	// 遍历
	v := &TraverseStatement{}
	(e.Stmt).Accept(v)
	// 如果没有解析到schema，赋值Schema为传入的schema
	var newV TraverseStatement = TraverseStatement{CteName: v.CteName}
	for _, item := range v.Tables {
		var tmpSchema string = item.Schema
		if tmpSchema == "" {
			tmpSchema = e.Schema
		}
		newV.Tables = append(newV.Tables, Table{Schema: tmpSchema, Table: item.Table})
	}
	// 移除重复的schema.table & cte name
	return removeElement(removeDuplicateElement(newV.Tables), newV.CteName)
}

// 移除重复的值
func removeDuplicateElement(data []Table) []Table {
	result := make([]Table, 0, len(data))
	temp := map[string]struct{}{}
	for _, item := range data {
		var key string = strings.Join([]string{item.Schema, item.Table}, "")
		if _, ok := temp[key]; !ok {
			temp[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func removeElement(data []Table, toBeRemoved []string) []Table {
	if toBeRemoved == nil || data == nil {
		return data
	}
	result := make([]Table, 0, len(data))
	temp := map[string]struct{}{}
	for _, item := range toBeRemoved {
		temp[item] = struct{}{}
	}
	for _, item := range data {
		if _, ok := temp[item.Table]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// 提取表结构体
type ExtractTables struct {
	Tables  []Table  // 库表名
	CteName []string // RecursiveCTENames
}

// 迭代select语句
func (e *ExtractTables) checkSelectItem(node ast.ResultSetNode) {
	if node == nil {
		return
	}
	// fmt.Println("类型: ", reflect.TypeOf(node))
	switch n := node.(type) {
	case *ast.SelectStmt:
		e.checkSubSelectItem(n)
	case *ast.Join:
		e.checkSelectItem(n.Left)
		e.checkSelectItem(n.Right)
	case *ast.TableSource:
		e.checkSelectItem(n.Source)
	case *ast.TableName:
		e.Tables = append(e.Tables, Table{Schema: n.Schema.L, Table: n.Name.L})
	}
}

func (e *ExtractTables) checkCTEItems(node *ast.WithClause) {
	if node.IsRecursive {
		for _, ctE := range node.CTEs {
			if ctE.IsRecursive {
				e.checkSelectItem(ctE.Query)
			}
			e.CteName = append(e.CteName, ctE.Name.L)
		}
	}
}

// 迭代子查询
func (e *ExtractTables) checkSubSelectItem(node *ast.SelectStmt) {
	if node.From != nil {
		// 迭代from子查询
		e.checkSelectItem(node.From.TableRefs)
	}
	if node.Where != nil {
		e.checkExprItem(node.Where)
	}
	for _, item := range node.Fields.Fields {
		if item.Expr != nil {
			e.checkExprItem(item.Expr)
		}
	}
}

// 迭代表达式
func (e *ExtractTables) checkExprItem(expr ast.ExprNode) {
	switch ex := expr.(type) {
	case *ast.PatternInExpr:
		e.checkExprItem(ex.Sel)
	case *ast.CompareSubqueryExpr:
		e.checkExprItem(ex.R)
	case *ast.BinaryOperationExpr:
		e.checkExprItem(ex.L)
		e.checkExprItem(ex.R)
	case *ast.ExistsSubqueryExpr:
		e.checkExprItem(ex.Sel)
	case *ast.SubqueryExpr:
		e.checkSelectItem(ex.Query)
	}
}

// 遍历语句
type TraverseStatement struct {
	Tables  []Table  // 库表名
	CteName []string // RecursiveCTENames
}

func (c *TraverseStatement) Enter(in ast.Node) (ast.Node, bool) {
	var e ExtractTables
	switch stmt := in.(type) {
	case *ast.ShowStmt:
		if stmt.Tp == ast.ShowCreateTable || stmt.Tp == ast.ShowCreateView || stmt.Tp == ast.ShowCreateDatabase {
			c.Tables = append(c.Tables, Table{Schema: stmt.Table.Schema.L, Table: stmt.Table.Name.L})
		}
	case *ast.SelectStmt:
		// 处理WITH语句
		if stmt.With != nil {
			e.checkCTEItems(stmt.With)
			c.CteName = append(c.CteName, e.CteName...)
			break
		}
		// select 1;
		if stmt.From == nil {
			break
		}
		e.checkSelectItem(stmt.From.TableRefs)
		e.checkExprItem(stmt.Where)
		if stmt.Having != nil {
			e.checkExprItem(stmt.Having.Expr)
		}
		for _, field := range stmt.Fields.Fields {
			e.checkExprItem(field.Expr)
		}
		if stmt.GroupBy != nil {
			for _, gb := range stmt.GroupBy.Items {
				e.checkExprItem(gb.Expr)
			}
		}
		c.Tables = append(c.Tables, e.Tables...)
	case *ast.InsertStmt:
		e.checkSelectItem(stmt.Table.TableRefs)
		e.checkSelectItem(stmt.Select)
		c.Tables = append(c.Tables, e.Tables...)
	case *ast.UpdateStmt:
		e.checkSelectItem(stmt.TableRefs.TableRefs)
		c.Tables = append(c.Tables, e.Tables...)
	case *ast.DeleteStmt:
		e.checkSelectItem(stmt.TableRefs.TableRefs)
		c.Tables = append(c.Tables, e.Tables...)
	case *ast.CreateTableStmt:
		c.Tables = append(c.Tables, Table{Schema: stmt.Table.Schema.L, Table: stmt.Table.Name.L})
	case *ast.CreateViewStmt:
		c.Tables = append(c.Tables, Table{Schema: stmt.ViewName.Schema.L, Table: stmt.ViewName.Name.L})
	case *ast.CreateIndexStmt:
		c.Tables = append(c.Tables, Table{Schema: stmt.Table.Schema.L, Table: stmt.Table.Name.L})
	case *ast.AlterTableStmt:
		c.Tables = append(c.Tables, Table{Schema: stmt.Table.Schema.L, Table: stmt.Table.Name.L})
	case *ast.DropIndexStmt:
		c.Tables = append(c.Tables, Table{Schema: stmt.Table.Schema.L, Table: stmt.Table.Name.L})
	case *ast.RenameTableStmt:
		for _, t := range stmt.TableToTables {
			c.Tables = append(c.Tables, Table{Schema: t.OldTable.Schema.L, Table: t.OldTable.Name.L})
			c.Tables = append(c.Tables, Table{Schema: t.NewTable.Schema.L, Table: t.NewTable.Name.L})
		}
	case *ast.DropTableStmt:
		for _, t := range stmt.Tables {
			c.Tables = append(c.Tables, Table{Schema: t.Schema.L, Table: t.Name.L})
		}
	case *ast.TruncateTableStmt:
		c.Tables = append(c.Tables, Table{Schema: stmt.Table.Schema.L, Table: stmt.Table.Name.L})
	}
	return in, false
}

func (c *TraverseStatement) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
