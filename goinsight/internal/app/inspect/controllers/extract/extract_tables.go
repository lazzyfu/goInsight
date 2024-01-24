/*
@Time    :   2022/09/09 10:35:02
@Author  :   zongfei.fu
@Desc    :   提取表名
*/

package extract

import (
	_ "embed"
	"fmt"
	"sqlSyntaxAudit/common/utils"
	"sqlSyntaxAudit/config"
	"sqlSyntaxAudit/controllers/parser"
	"sqlSyntaxAudit/forms"
	logger "sqlSyntaxAudit/middleware/log"
	"sync"

	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/sirupsen/logrus"
)

// 移除重复的值
func removeDuplicateElement(data []string) []string {
	result := make([]string, 0, len(data))
	temp := map[string]struct{}{}
	for _, item := range data {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func removeElement(data, toBeRemoved []string) []string {
	if toBeRemoved == nil || data == nil {
		return data
	}
	result := make([]string, 0, len(data))
	temp := map[string]struct{}{}
	for _, item := range toBeRemoved {
		temp[item] = struct{}{}
	}
	for _, item := range data {
		if _, ok := temp[item]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// 返回数据
type ReturnData struct {
	Tables []string `json:"tables"` // 表名
	Type   string   `json:"type"`   // 语句类型
	Query  string   `json:"query"`  // 原始SQL
}

// 检查结构体
type Checker struct {
	Form  forms.ExtractTablesForm
	Audit *config.Audit
}

func (c *Checker) Extract(RequestID string) (error, []ReturnData) {
	var returnData []ReturnData
	err := c.Parse()
	if err != nil {
		logger.AppLog.WithFields(logrus.Fields{"request_id": RequestID}).Error(err)
		return err, returnData
	}
	for _, stmt := range c.Audit.TiStmt {
		var data ReturnData = ReturnData{Query: stmt.Text()}
		data.Tables, data.Type = ExtractTablesFromStatement(&stmt)
		returnData = append(returnData, data)
	}
	return nil, returnData
}

// 解析SQL语句
func (c *Checker) Parse() error {
	// 解析SQL
	var warns []error
	var err error
	// 解析
	c.Audit, warns, err = parser.NewParse(c.Form.SqlText, "", "")
	if len(warns) > 0 {
		return fmt.Errorf("Parse Warning: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return fmt.Errorf("sql解析错误：%s", err.Error())
	}
	return nil
}

// 提取表结构体
type ExtractTables struct {
	Tables            []string // 表名
	RecursiveCTENames []string // 递归CTE别名
}

// 迭代select语句
func (e *ExtractTables) checkSelectItem(node ast.ResultSetNode) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.SelectStmt:
		e.checkSubSelectItem(n)
	case *ast.Join:
		e.checkSelectItem(n.Left)
		e.checkSelectItem(n.Right)
	case *ast.TableSource:
		e.checkSelectItem(n.Source)
	case *ast.TableName:
		e.Tables = append(e.Tables, n.Name.L)
	}
}

func (e *ExtractTables) checkCTEItems(node *ast.WithClause) {
	if node.IsRecursive {
		for _, ctE := range node.CTEs {
			e.RecursiveCTENames = append(e.RecursiveCTENames, ctE.Name.L)
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

// TraverseStatement
type TraverseStatement struct {
	Tables            []string // 表名
	RecursiveCTENames []string // 递归CTE别名
	Type              string   // 语句类型
	setTypeOnce       sync.Once
}

func (c *TraverseStatement) setType(typ string) {
	// 只设置一次Type, 避免类INSERT INTO SELECT获取Type错误
	c.setTypeOnce.Do(func() {
		c.Type = typ
	})
}

func (c *TraverseStatement) Enter(in ast.Node) (ast.Node, bool) {
	var e ExtractTables
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		c.setType("SELECT")
		// 处理WITH语句
		if stmt.With != nil {
			e.checkCTEItems(stmt.With)
			c.RecursiveCTENames = append(c.RecursiveCTENames, e.RecursiveCTENames...)
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
		if stmt.IsReplace {
			c.setType("REPLACE")
		} else {
			c.setType("INSERT")
		}
		e.checkSelectItem(stmt.Table.TableRefs)
		e.checkSelectItem(stmt.Select)
		c.Tables = append(c.Tables, e.Tables...)
	case *ast.UpdateStmt:
		c.setType("UPDATE")
		e.checkSelectItem(stmt.TableRefs.TableRefs)
		c.Tables = append(c.Tables, e.Tables...)
	case *ast.DeleteStmt:
		c.setType("DELETE")
		e.checkSelectItem(stmt.TableRefs.TableRefs)
		c.Tables = append(c.Tables, e.Tables...)
	case *ast.CreateTableStmt:
		c.setType("CREATE TABLE")
		c.Tables = append(c.Tables, stmt.Table.Name.L)
	case *ast.CreateViewStmt:
		c.setType("CREATE VIEW")
		c.Tables = append(c.Tables, stmt.ViewName.Name.L)
	case *ast.CreateIndexStmt:
		c.setType("CREATE INDEX")
		c.Tables = append(c.Tables, stmt.Table.Name.L)
	case *ast.AlterTableStmt:
		c.setType("ALTER TABLE")
		c.Tables = append(c.Tables, stmt.Table.Name.L)
	case *ast.DropIndexStmt:
		c.setType("DROP INDEX")
		c.Tables = append(c.Tables, stmt.Table.Name.L)
	case *ast.RenameTableStmt:
		c.setType("RENAME TABLE")
		for _, t := range stmt.TableToTables {
			c.Tables = append(c.Tables, t.OldTable.Name.L)
			c.Tables = append(c.Tables, t.NewTable.Name.L)
		}
	case *ast.DropTableStmt:
		if stmt.IsView {
			c.setType("DROP VIEW")
		} else {
			c.setType("DROP TABLE")
		}
		for _, t := range stmt.Tables {
			c.Tables = append(c.Tables, t.Name.L)
		}
	case *ast.TruncateTableStmt:
		c.setType("TRUNCATE TABLE")
		c.Tables = append(c.Tables, stmt.Table.Name.L)
	}
	return in, false
}

func (c *TraverseStatement) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func ExtractTablesFromStatement(stmt *ast.StmtNode) ([]string, string) {
	v := &TraverseStatement{}
	(*stmt).Accept(v)
	return removeElement(removeDuplicateElement(v.Tables), v.RecursiveCTENames), v.Type
}
