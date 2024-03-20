/*
@Time    :   2023/04/11 15:20:32
@Author  :   zongfei.fu
@Desc    :   提取limit
*/

package parser

import (
	"sync"

	"github.com/pingcap/tidb/pkg/parser/ast"
	driver "github.com/pingcap/tidb/pkg/types/parser_driver"
)

type Limit struct {
	once   sync.Once
	Count  uint64
	Offset uint64
}

func (l *Limit) Enter(in ast.Node) (ast.Node, bool) {
	l.once.Do(func() {
		switch stmt := in.(type) {
		case *ast.SelectStmt:
			// SQL语句没有指定limit
			if stmt.Limit == nil {
				l.Count = 0
				l.Offset = 0
			} else {
				// SQL语句指定了limit
				if stmt.Limit.Count != nil {
					// 获取count
					switch ex := stmt.Limit.Count.(type) {
					case *driver.ValueExpr:
						l.Count = ex.GetUint64()
					}
				}
				if stmt.Limit.Offset != nil {
					// 获取offset
					switch ex := stmt.Limit.Offset.(type) {
					case *driver.ValueExpr:
						l.Offset = ex.GetUint64()
					}
				}
			}
		case *ast.SetOprStmt:
			// SQL语句没有指定limit
			if stmt.Limit == nil {
				l.Count = 0
				l.Offset = 0
			} else {
				// SQL语句指定了limit
				if stmt.Limit.Count != nil {
					// 获取count
					switch ex := stmt.Limit.Count.(type) {
					case *driver.ValueExpr:
						l.Count = ex.GetUint64()
					}
				}
				if stmt.Limit.Offset != nil {
					// 获取offset
					switch ex := stmt.Limit.Offset.(type) {
					case *driver.ValueExpr:
						l.Offset = ex.GetUint64()
					}
				}
			}
		}
	})
	return in, false
}

func (s *Limit) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
