package base

import (
	"fmt"
	"reflect"

	"vitess.io/vitess/go/vt/sqlparser"
)

// Rewrite 用于重写SQL
type Rewrite struct {
	SQL  string
	Stmt sqlparser.Statement
}

// NewRewrite 返回一个*Rewrite对象，如果SQL无法被正常解析，返回nil和错误
func NewRewrite(sql string) (*Rewrite, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}

	return &Rewrite{
		SQL:  sql,
		Stmt: stmt,
	}, err
}

// Rewrite 入口函数
func (rw *Rewrite) Rewrite() error {
	return rw.RewriteDML2Select()
}

// RewriteDML2Select dml2select: DML 转成 SELECT，兼容低版本的 EXPLAIN
func (rw *Rewrite) RewriteDML2Select() error {
	if rw.Stmt == nil {
		return nil
	}
	switch stmt := rw.Stmt.(type) {
	case *sqlparser.Select:
		return nil
	case *sqlparser.Delete: // Multi DELETE not support yet.
		rw.SQL = delete2Select(stmt)
	case *sqlparser.Insert:
		rw.SQL = insert2Select(stmt)
	case *sqlparser.Update: // Multi UPDATE not support yet.
		rw.SQL = update2Select(stmt)
	}
	var err error
	rw.Stmt, err = sqlparser.Parse(rw.SQL)
	return err
}

// delete2Select 将 Delete 语句改写成 Select
func delete2Select(stmt *sqlparser.Delete) string {
	newSQL := &sqlparser.Select{
		SelectExprs: []sqlparser.SelectExpr{
			new(sqlparser.StarExpr),
		},
		From:    stmt.TableExprs,
		Where:   stmt.Where,
		OrderBy: stmt.OrderBy,
		Limit:   stmt.Limit,
	}
	return sqlparser.String(newSQL)
}

// update2Select 将 Update 语句改写成 Select
func update2Select(stmt *sqlparser.Update) string {
	newSQL := &sqlparser.Select{
		SelectExprs: []sqlparser.SelectExpr{
			new(sqlparser.StarExpr),
		},
		From:    stmt.TableExprs,
		Where:   stmt.Where,
		OrderBy: stmt.OrderBy,
		Limit:   stmt.Limit,
	}
	return sqlparser.String(newSQL)
}

// insert语句不需要生成回滚语句，如果insert 的数据不唯一且不包含主键，生成的回滚SQL是有问题的
// insert2Select 将 Insert 语句改写成 Select
func insert2Select(stmt *sqlparser.Insert) string {
	fmt.Println("ccc: ", reflect.TypeOf(stmt.Rows))
	switch row := stmt.Rows.(type) {
	case *sqlparser.Select, *sqlparser.Union:
		return sqlparser.String(row)
	case *sqlparser.Values:
		fmt.Println(row)
	}

	return "select 1 from DUAL"
}
