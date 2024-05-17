/*
@Time    :   2022/06/27 15:49:13
@Author  :   xff
@Desc    :   None
*/

package process

import (
	"fmt"
	"goInsight/internal/pkg/utils"

	"github.com/pingcap/tidb/pkg/parser/mysql"
)

// 字符集
type TableCharset struct {
	Table   string
	Charset string
	Collate string
}
type ColumnCharset struct {
	Table   string
	Column  string
	Tp      byte
	Charset string
	Collate string
}
type Charset struct {
	SupportCharset     []string
	RecommendCollation string
	Table              TableCharset
	Cols               []ColumnCharset
}

// 检查表的字符集
func (c *Charset) CheckTable() error {
	/*
		1.表必须指定字符集
		1.指定的字符集必须是可选的字符集
		2.指定排序规则的同时必须指定字符集
		3.检查字符集排序规则是否符合要求
	*/
	if !utils.IsContain(c.SupportCharset, c.Table.Charset) {
		if c.Table.Charset == "" {
			return fmt.Errorf("表`%s`必须指定字符集，可选字符集为%s【例如:DEFAULT CHARSET=utf8mb4】", c.Table.Table, c.SupportCharset)
		} else {
			return fmt.Errorf("表`%s`指定的字符集`%s`不符合要求，可选字符集为%s", c.Table.Table, c.Table.Charset, c.SupportCharset)
		}
	}

	if len(c.Table.Collate) > 0 && len(c.Table.Charset) > 0 && !utils.HasPrefix(c.Table.Collate, c.Table.Charset+"_", false) {
		// 检查排序规则的前缀，前缀必须为字符集+"_"
		return fmt.Errorf("表`%s`指定的字符集排序规则`%s`不符合要求，应指定前缀为%s的排序规则，推荐的字符集排序规则为%s", c.Table.Table, c.Table.Charset, c.Table.Charset+"_", c.RecommendCollation)
	}
	return nil
}

// 检查列字符集
func (c *Charset) CheckColumn() error {
	// https://dev.mysql.com/doc/refman/5.6/en/charset-column.html
	// 不支持的类型，语法检查就被拦截了；此处JSON类型为BINARY
	for _, col := range c.Cols {
		if col.Tp == mysql.TypeVarchar || col.Tp == mysql.TypeVarString || col.Tp == mysql.TypeString || col.Tp == mysql.TypeEnum || col.Tp == mysql.TypeSet {
			// 检查字符集
			if len(col.Charset) == 0 || len(col.Collate) == 0 {
				return fmt.Errorf("列`%s`必须同时指定字符集和排序规则[表`%s`]", col.Column, col.Table)
			}
			// 检查排序规则
			if !utils.HasPrefix(col.Collate, col.Charset+"_", false) && len(col.Collate) > 0 {
				// 检查排序规则的前缀，前缀必须为字符集+"_"
				return fmt.Errorf("列`%s`指定的字符集排序规则`%s`不符合要求，应指定前缀为%s的排序规则[表`%s`]", col.Column, col.Charset, col.Charset+"_", col.Table)
			}
		}
	}

	return nil
}
