/*
@Time    :   2022/06/24 13:12:20
@Author  :   xff
@Desc    :   遍历语法树
*/

package traverses

import (
	"goInsight/internal/inspect/controllers/process"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

// TraverseDisableAuditDMLTables
type TraverseDisableAuditDMLTables struct {
	Tables []string
}

// TraverseDMLInsertIntoSelect
type TraverseDMLInsertIntoSelect struct {
	IsMatch           int // 是否匹配当前规则
	DMLType           string
	HasSelectSubQuery bool // 有 insert/replace into ... select ...
	HasOnDuplicate    bool // 有 OnDuplicate
}

func (c *TraverseDMLInsertIntoSelect) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.InsertStmt:
		c.IsMatch++
		c.DMLType = "INSERT"
		if stmt.IsReplace {
			c.DMLType = "REPLACE"
		}
		if stmt.Select != nil {
			c.HasSelectSubQuery = true
		}
		if !stmt.IsReplace && len(stmt.OnDuplicate) > 0 {
			c.HasOnDuplicate = true
		}
	}
	return in, false
}

func (c *TraverseDMLInsertIntoSelect) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseDMLNoWhere
type TraverseDMLNoWhere struct {
	IsMatch  int // 是否匹配当前规则
	HasWhere bool
	DMLType  string
}

func (c *TraverseDMLNoWhere) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.DeleteStmt:
		c.IsMatch++
		c.DMLType = "DELETE"
		if stmt.Where != nil {
			c.HasWhere = true
		}
	case *ast.UpdateStmt:
		c.IsMatch++
		c.DMLType = "UPDATE"
		if stmt.Where != nil {
			c.IsMatch++
			c.HasWhere = true
		}
	}
	return in, false
}

func (c *TraverseDMLNoWhere) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseDMLInsertWithColumns
type TraverseDMLInsertWithColumns struct {
	Table             string
	IsMatch           int      // 是否匹配当前规则
	Columns           []string // 列名
	ColumnsCount      int      // 指定的列的数量
	ColsValuesIsMatch bool     // 列的数量是否和值的数量匹配
	RowsCount         int      // 一次insert行的数量
	DMLType           string
}

func (c *TraverseDMLInsertWithColumns) CheckSelectItem(node ast.ResultSetNode) {
	// 提取表名
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.Join:
		c.CheckSelectItem(n.Left)
		c.CheckSelectItem(n.Right)
	case *ast.TableSource:
		c.CheckSelectItem(n.Source)
	case *ast.TableName:
		c.Table = n.Name.String()
	}
}

func (c *TraverseDMLInsertWithColumns) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.InsertStmt:
		c.CheckSelectItem(stmt.Table.TableRefs)
		c.IsMatch++
		c.DMLType = "INSERT"
		if stmt.IsReplace {
			c.DMLType = "REPLACE"
		}
		c.ColumnsCount = len(stmt.Columns)
		for _, item := range stmt.Columns {
			c.Columns = append(c.Columns, item.Name.O)
		}
		c.RowsCount = len(stmt.Lists)
		c.ColsValuesIsMatch = true
		for _, row := range stmt.Lists {
			if len(row) != c.ColumnsCount {
				c.ColsValuesIsMatch = false
				break
			}
		}
	}
	return in, false
}

func (c *TraverseDMLInsertWithColumns) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseDMLHasConstraint
type TraverseDMLHasConstraint struct {
	IsMatch     int    // 是否匹配当前规则
	DMLType     string // 语句类型，delete or update
	HasLimit    bool
	HasOrderBy  bool
	HasSubQuery bool
}

func (c *TraverseDMLHasConstraint) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.DeleteStmt:
		c.IsMatch++
		if stmt.Limit != nil {
			c.DMLType = "DELETE"
			c.HasLimit = true
		}
		if stmt.Order != nil {
			c.HasOrderBy = true
		}
		if stmt.Where != nil {
			check := process.CheckItem{}
			check.CheckExprItem(stmt.Where)
			c.HasLimit = check.HasLimit
			c.HasOrderBy = check.HasOrderBy
			c.HasSubQuery = check.HasSubQuery
		}
	case *ast.UpdateStmt:
		c.IsMatch++
		if stmt.Limit != nil {
			c.DMLType = "UPDATE"
			c.HasLimit = true
		}
		if stmt.Order != nil {
			c.HasOrderBy = true
		}
		// eg: update a set a.price = (subquery) where expr;
		for _, item := range stmt.List {
			if item.Expr != nil {
				check := process.CheckItem{}
				check.CheckExprItem(item.Expr)
				c.HasLimit = check.HasLimit
				c.HasOrderBy = check.HasOrderBy
				c.HasSubQuery = check.HasSubQuery
			}
		}
		// eg: update a set a.price=188 where a.id IN (subquery);
		if stmt.Where != nil {
			check := process.CheckItem{}
			check.CheckExprItem(stmt.Where)
			c.HasLimit = check.HasLimit
			c.HasOrderBy = check.HasOrderBy
			c.HasSubQuery = check.HasSubQuery
		}
	}
	return in, false
}

func (c *TraverseDMLHasConstraint) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseDMLJoinWithOn
type TraverseDMLJoinWithOn struct {
	DMLType      string
	IsMatch      int  // 是否匹配当前规则
	HasJoin      bool // 是否为join语句
	IsJoinWithOn bool // join是否有ON语句
}

func (c *TraverseDMLJoinWithOn) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.DeleteStmt:
		c.IsMatch++
		c.DMLType = "DELETE"
		// 检查JOIN操作是否有ON语句
		// eg: delete a from a join b on a.id = b.id where expr;
		if stmt.TableRefs.TableRefs.Tp > 0 {
			c.HasJoin = true
			if stmt.TableRefs.TableRefs.On != nil {
				c.IsJoinWithOn = true
			}
		}
	case *ast.UpdateStmt:
		c.IsMatch++
		c.DMLType = "UPDATE"
		// 检查JOIN操作是否有ON语句
		// eg: update a join b on a.id = b.id set a.is_deleted=1 where b.id is null;
		if stmt.TableRefs.TableRefs.Tp > 0 {
			c.HasJoin = true
			if stmt.TableRefs.TableRefs.On != nil {
				c.IsJoinWithOn = true
			}
		}
	}
	return in, false
}

func (c *TraverseDMLJoinWithOn) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseDMLMaxUpdateRows
type TraverseDMLMaxUpdateRows struct {
	IsMatch int // 是否匹配当前规则
	DMLType string
}

func (c *TraverseDMLMaxUpdateRows) Enter(in ast.Node) (ast.Node, bool) {
	switch in.(type) {
	case *ast.DeleteStmt:
		c.IsMatch++
		c.DMLType = "DELETE"
	case *ast.UpdateStmt:
		c.IsMatch++
		c.DMLType = "UPDATE"
	}
	return in, false
}

func (c *TraverseDMLMaxUpdateRows) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseDMLMaxInsertRows
type TraverseDMLMaxInsertRows struct {
	IsMatch   int // 是否匹配当前规则
	DMLType   string
	RowsCount int // 一次insert行的数量
}

func (c *TraverseDMLMaxInsertRows) Enter(in ast.Node) (ast.Node, bool) {
	switch stmt := in.(type) {
	case *ast.InsertStmt:
		c.IsMatch++
		c.DMLType = "INSERT"
		c.RowsCount = len(stmt.Lists)
	}
	return in, false
}

func (c *TraverseDMLMaxInsertRows) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
