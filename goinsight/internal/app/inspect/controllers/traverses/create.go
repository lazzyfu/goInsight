/*
@Time    :   2022/06/24 13:12:20
@Author  :   zongfei.fu
@Desc    :   遍历语法树,语法参考pingcap文档：https://github.com/pingcap/parser/blob/master/docs/quickstart.md
*/

package traverses

import (
	"goInsight/internal/app/inspect/controllers/process"
	"goInsight/internal/pkg/utils"
	"strings"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/mysql"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

// TraverseCreateTableIsExist
type TraverseCreateTableIsExist struct {
	Table string // 表名
}

func (c *TraverseCreateTableIsExist) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
	}
	return in, false
}

func (c *TraverseCreateTableIsExist) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableAs
type TraverseCreateTableAs struct {
	Table      string // 表名
	IsCreateAs bool   // 是否为create table as语法
}

func (c *TraverseCreateTableAs) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		switch stmt.Select.(type) {
		case *ast.SelectStmt:
			c.IsCreateAs = true
		}
	}
	return in, false
}

func (c *TraverseCreateTableAs) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableLike
type TraverseCreateTableLike struct {
	Table        string // 表名
	IsCreateLike bool   // 是否为create table like语法
}

func (c *TraverseCreateTableLike) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		if stmt.ReferTable != nil {
			c.IsCreateLike = true
		}
	}
	return in, false
}

func (c *TraverseCreateTableLike) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableOptions
type TraverseCreateTableOptions struct {
	process.TableOptions
}

func (c *TraverseCreateTableOptions) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.TableOptions.Table = stmt.Table.Name.String()
		for _, node := range stmt.Options {
			switch node.Tp {
			case ast.TableOptionEngine:
				c.TableOptions.Engine = node.StrValue
			case ast.TableOptionCharset:
				c.TableOptions.Charset = node.StrValue
			case ast.TableOptionCollate:
				c.TableOptions.Collate = node.StrValue
			case ast.TableOptionAutoIncrement:
				c.TableOptions.AutoIncrement = node.UintValue
			case ast.TableOptionRowFormat:
				c.TableOptions.RowFormat = node.StrValue
			case ast.TableOptionComment:
				c.TableOptions.HasComment = true       // 表示有注释，代表不了注释为空
				c.TableOptions.Comment = node.StrValue // 获取注释的值
			}
		}
		if stmt.Partition != nil {
			c.TableOptions.PartitionType = stmt.Partition.PartitionMethod.Tp.String()
			c.TableOptions.IsPartition = true
		}
	}
	return in, false
}

func (c *TraverseCreateTableOptions) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTablePrimaryKey
type TraverseCreateTablePrimaryKey struct {
	Table       string               // 表名
	PrimaryKeys []process.PrimaryKey // 主键
}

func (c *TraverseCreateTablePrimaryKey) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		var keys []string
		// 从I_ID bigint unsigned NOT NULL AUTO_INCREMENT primary key COMMENT '自增ID' 提取主键
		for _, col := range stmt.Cols {
			for _, opt := range col.Options {
				if opt.Tp == ast.ColumnOptionPrimaryKey {
					keys = append(keys, col.Name.Name.O)
				}
			}
		}
		// 从PRIMARY KEY (I_ID) 提取主键
		for _, cons := range stmt.Constraints {
			if cons.Tp == ast.ConstraintPrimaryKey {
				for _, col := range cons.Keys {
					// 去重
					if !utils.IsContain(keys, col.Column.Name.O) {
						keys = append(keys, col.Column.Name.O)
					}
				}
			}
		}
		// 获取主键的属性
		for _, col := range stmt.Cols {
			for _, key := range keys {
				if col.Name.Name.O == key {
					var HasAutoIncrement, HasNotNull bool
					for _, opt := range col.Options {
						switch opt.Tp {
						case ast.ColumnOptionAutoIncrement:
							HasAutoIncrement = true
						case ast.ColumnOptionNotNull:
							HasNotNull = true
						}
					}
					c.PrimaryKeys = append(c.PrimaryKeys, process.PrimaryKey{
						Table:            c.Table,
						Column:           col.Name.Name.O,
						Tp:               col.Tp.GetType(),
						Flag:             col.Tp.GetFlag(),
						HasNotNull:       HasNotNull,
						HasAutoIncrement: HasAutoIncrement,
					})
				}
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTablePrimaryKey) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableConstraint
type TraverseCreateTableConstraint struct {
	Table        string // 表名
	IsForeignKey bool   // 是否定义了外键
}

func (c *TraverseCreateTableConstraint) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, cons := range stmt.Constraints {
			if cons.Tp == ast.ConstraintForeignKey {
				c.IsForeignKey = true
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableConstraint) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableAuditCols
type TraverseCreateTableAuditCols struct {
	Table     string // 表名
	AuditCols []map[string]string
}

func (c *TraverseCreateTableAuditCols) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, node := range stmt.Cols {
			var colOptions []string
			if node.Tp.GetType() == mysql.TypeDatetime || node.Tp.GetType() == mysql.TypeTimestamp {
				for _, o := range node.Options {
					if o.Tp == ast.ColumnOptionDefaultValue {
						// DEFAULT CURRENT_TIMESTAMP
						if funcCall, ok := o.Expr.(*ast.FuncCallExpr); ok {
							if funcCall.FnName.L == ast.CurrentTimestamp {
								colOptions = append(colOptions, "DEFAULT CURRENT_TIMESTAMP")
							}
						}
					}
					if o.Tp == ast.ColumnOptionOnUpdate {
						// ON UPDATE CURRENT_TIMESTAMP
						colOptions = append(colOptions, "ON UPDATE")
						if funcCall, ok := o.Expr.(*ast.FuncCallExpr); ok {
							if funcCall.FnName.L == ast.CurrentTimestamp {
								colOptions = append(colOptions, "CURRENT_TIMESTAMP")
							}
						}
					}
				}
			}
			if len(colOptions) > 0 {
				c.AuditCols = append(c.AuditCols, map[string]string{node.Name.Name.O: strings.Join(colOptions, " ")})
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableAuditCols) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableColsOptions
type TraverseCreateTableColsOptions struct {
	Table       string   // 表名
	Charset     string   // 表字符集
	PrimaryKeys []string // 主键
	Cols        []process.ColOptions
}

func (c *TraverseCreateTableColsOptions) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		// 主键、自增键不检测默认值
		func() {
			// 从I_ID bigint unsigned NOT NULL AUTO_INCREMENT primary key COMMENT '自增ID' 提取主键
			for _, col := range stmt.Cols {
				for _, opt := range col.Options {
					switch opt.Tp {
					case ast.ColumnOptionPrimaryKey, ast.ColumnOptionAutoIncrement:
						c.PrimaryKeys = append(c.PrimaryKeys, col.Name.Name.O)
					}
				}
			}
			// 从PRIMARY KEY (I_ID) 提取主键
			for _, cons := range stmt.Constraints {
				if cons.Tp == ast.ConstraintPrimaryKey {
					for _, col := range cons.Keys {
						// 去重
						if !utils.IsContain(c.PrimaryKeys, col.Column.Name.O) {
							c.PrimaryKeys = append(c.PrimaryKeys, col.Column.Name.O)
						}
					}
				}
			}
		}()

		for _, col := range stmt.Cols {
			var NotNullFlag, HasDefaultValue, DefaultIsNull, HasComment bool
			var DefaultValue interface{}
			for _, op := range col.Options {
				switch op.Tp {
				case ast.ColumnOptionComment:
					HasComment = true
				case ast.ColumnOptionDefaultValue:
					HasDefaultValue = true
					if v, ok := op.Expr.(*driver.ValueExpr); ok {
						// 有默认值，且为NULL，且有NOT NULL约束，如(not null default null)
						DefaultIsNull = v.Datum.IsNull()
						// 获取列的默认值
						DefaultValue = v.GetValue()
					}
					if f, ok := op.Expr.(*ast.FuncCallExpr); ok {
						DefaultValue = f.FnName.L
					}
				case ast.ColumnOptionNotNull:
					NotNullFlag = true
				}
			}
			if !utils.IsContain(c.PrimaryKeys, col.Name.Name.O) {
				c.Cols = append(c.Cols,
					process.ColOptions{
						Table:           c.Table,
						Column:          col.Name.Name.O,
						Tp:              col.Tp.GetType(),
						Flen:            col.Tp.GetFlen(),
						NotNullFlag:     NotNullFlag,
						HasDefaultValue: HasDefaultValue,
						DefaultValue:    DefaultValue,
						DefaultIsNull:   DefaultIsNull,
						HasComment:      HasComment,
					})
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableColsOptions) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableColsRepeatDefine
type TraverseCreateTableColsRepeatDefine struct {
	Table string   // 表名
	Cols  []string // 列名
}

func (c *TraverseCreateTableColsRepeatDefine) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		// 获取所有的列
		for _, col := range stmt.Cols {
			c.Cols = append(c.Cols, col.Name.Name.L)
		}
	}
	return in, false
}

func (c *TraverseCreateTableColsRepeatDefine) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableColsCharset
type TraverseCreateTableColsCharset struct {
	process.Charset
}

func (c *TraverseCreateTableColsCharset) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		// 获取列字符集和排序规则
		for _, col := range stmt.Cols {
			var colCharset, colCollate string
			colCharset = col.Tp.GetCharset()
			colCollate = col.Tp.GetCollate()
			for _, v := range col.Options {
				if v.Tp == ast.ColumnOptionCollate {
					colCollate = v.StrValue
				}
			}
			// 设置了字符集和排序规则
			if colCharset != "" || colCollate != "" {
				c.Charset.Cols = append(
					c.Charset.Cols,
					process.ColumnCharset{
						Table:   stmt.Table.Name.String(),
						Column:  col.Name.Name.O,
						Tp:      col.Tp.GetType(),
						Charset: colCharset,
						Collate: colCollate,
					},
				)
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableColsCharset) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableIndexesPrefix
type TraverseCreateTableIndexesPrefix struct {
	Table  string              // 表名
	Prefix process.IndexPrefix // 前缀
}

func (c *TraverseCreateTableIndexesPrefix) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.Prefix.Table = stmt.Table.Name.String()
		// 前缀
		for _, cons := range stmt.Constraints {
			switch cons.Tp {
			case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
				c.Prefix.UniqueKeys = append(c.Prefix.UniqueKeys, cons.Name)
			case ast.ConstraintIndex, ast.ConstraintKey:
				c.Prefix.SecondaryKeys = append(c.Prefix.SecondaryKeys, cons.Name)
			case ast.ConstraintFulltext:
				c.Prefix.SecondaryKeys = append(c.Prefix.FulltextKeys, cons.Name)
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableIndexesPrefix) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableIndexesCount
type TraverseCreateTableIndexesCount struct {
	Table  string              // 表名
	Number process.IndexNumber // 索引数量
}

func (c *TraverseCreateTableIndexesCount) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.Number.Table = stmt.Table.Name.String()
		// 索引数量
		for _, cons := range stmt.Constraints {
			switch cons.Tp {
			case ast.ConstraintPrimaryKey:
				c.Number.Keys = append(c.Number.Keys, process.IndexLen{Index: "PrimaryKey", Len: len(cons.Keys)})
			case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
				c.Number.Number++
				c.Number.Keys = append(c.Number.Keys, process.IndexLen{Index: cons.Name, Len: len(cons.Keys)})
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableIndexesCount) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableIndexesRepeatDefine
type TraverseCreateTableIndexesRepeatDefine struct {
	Table   string   // 表名
	Indexes []string // 索引名
}

func (c *TraverseCreateTableIndexesRepeatDefine) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		// 获取所有的列
		for _, cons := range stmt.Constraints {
			switch cons.Tp {
			case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintFulltext:
				c.Indexes = append(c.Indexes, cons.Name)
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableIndexesRepeatDefine) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableRedundantIndexes
type TraverseCreateTableRedundantIndexes struct {
	Table     string                 // 表名
	Redundant process.RedundantIndex // 冗余索引
}

func (c *TraverseCreateTableRedundantIndexes) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.Redundant.Table = stmt.Table.Name.String()
		// 冗余索引
		for _, col := range stmt.Cols {
			// 获取所有的列
			c.Redundant.Cols = append(c.Redundant.Cols, col.Name.Name.L)
		}
		for _, cons := range stmt.Constraints {
			if cons.Tp == ast.ConstraintKey || cons.Tp == ast.ConstraintIndex ||
				cons.Tp == ast.ConstraintUniq || cons.Tp == ast.ConstraintUniqKey ||
				cons.Tp == ast.ConstraintUniqIndex || cons.Tp == ast.ConstraintFulltext {
				c.Redundant.Indexes = append(c.Redundant.Indexes, cons.Name)
				var idxColsMap []string
				for _, v := range cons.Keys {
					idxColsMap = append(idxColsMap, v.Column.Name.L)
				}
				c.Redundant.IndexesCols = append(c.Redundant.IndexesCols, process.IndexColsMap{Index: cons.Name, Cols: idxColsMap, Tag: "is_meta"})
			}
		}
	}
	return in, false
}

func (c *TraverseCreateTableRedundantIndexes) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableDisabledIndexes
type TraverseCreateTableDisabledIndexes struct {
	Table           string                  // 表名
	DisabledIndexes process.DisabledIndexes // 禁止创建索引的列类型
}

func (c *TraverseCreateTableDisabledIndexes) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.DisabledIndexes.Table = stmt.Table.Name.String()
		// 禁止创建索引的列类型
		for _, col := range stmt.Cols {
			if col.Tp.GetType() == mysql.TypeBlob ||
				col.Tp.GetType() == mysql.TypeTinyBlob ||
				col.Tp.GetType() == mysql.TypeMediumBlob ||
				col.Tp.GetType() == mysql.TypeLongBlob {
				c.DisabledIndexes.Cols = append(c.DisabledIndexes.Cols, col.Name.Name.O)
			}
		}
		for _, cons := range stmt.Constraints {
			var idxColsMap []string
			for _, v := range cons.Keys {
				idxColsMap = append(idxColsMap, v.Column.Name.L)
			}
			c.DisabledIndexes.IndexesCols = append(c.DisabledIndexes.IndexesCols, process.IndexColsMap{Index: cons.Name, Cols: idxColsMap})
		}
	}
	return in, false
}

func (c *TraverseCreateTableDisabledIndexes) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableInnodbLargePrefix
type TraverseCreateTableInnodbLargePrefix struct {
	LargePrefix process.LargePrefix
}

func (c *TraverseCreateTableInnodbLargePrefix) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.LargePrefix.Table = stmt.Table.Name.String()
		for _, node := range stmt.Options {
			switch node.Tp {
			case ast.TableOptionCharset:
				c.LargePrefix.Charset = node.StrValue
			}
		}
		var LargePrefixIndexColsMaps []process.LargePrefixIndexColsMap
		for _, cons := range stmt.Constraints {
			switch cons.Tp {
			case ast.ConstraintKey, ast.ConstraintIndex, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintFulltext:
				var LargePrefixIndexColsMap process.LargePrefixIndexColsMap
				LargePrefixIndexColsMap.Name = cons.Name
				for _, col := range cons.Keys {
					LargePrefixIndexColsMap.Keys = append(LargePrefixIndexColsMap.Keys,
						process.LargePrefixIndexPartSpecification{Column: col.Column.Name.L, Ilen: col.Length})
				}
				LargePrefixIndexColsMaps = append(LargePrefixIndexColsMaps, LargePrefixIndexColsMap)
			}
		}
		for _, i := range LargePrefixIndexColsMaps {
			var tmpLargePrefixIndexColsMap process.LargePrefixIndexColsMap
			for _, j := range i.Keys {
				for _, col := range stmt.Cols {
					if j.Column == col.Name.Name.L {
						tmpLargePrefixIndexColsMap.Name = i.Name
						tmpLargePrefixIndexColsMap.Keys = append(tmpLargePrefixIndexColsMap.Keys,
							process.LargePrefixIndexPartSpecification{
								Column:  j.Column,
								Ilen:    j.Ilen,
								Tp:      col.Tp.GetType(),
								Flen:    col.Tp.GetFlen(),
								Decimal: col.Tp.GetDecimal(),
								Charset: col.Tp.GetCharset(),
								Elems:   col.Tp.GetElems(),
							})
					}
				}
			}
			c.LargePrefix.LargePrefixIndexColsMaps = append(c.LargePrefix.LargePrefixIndexColsMaps, tmpLargePrefixIndexColsMap)
		}
	}
	return in, false
}

func (c *TraverseCreateTableInnodbLargePrefix) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableRowSizeTooLarge
type TraverseCreateTableRowSizeTooLarge struct {
	RowSizeTooLarge process.RowSizeTooLarge
}

func (c *TraverseCreateTableRowSizeTooLarge) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.RowSizeTooLarge.Table = stmt.Table.Name.String()
		for _, node := range stmt.Options {
			switch node.Tp {
			case ast.TableOptionCharset:
				c.RowSizeTooLarge.Charset = node.StrValue
			}
		}
		for _, col := range stmt.Cols {
			c.RowSizeTooLarge.RowSizeTooLargeColsMaps = append(c.RowSizeTooLarge.RowSizeTooLargeColsMaps,
				process.RowSizeTooLargePartSpecification{
					Column:  col.Name.Name.L,
					Tp:      col.Tp.GetType(),
					Flen:    col.Tp.GetFlen(),
					Decimal: col.Tp.GetDecimal(),
					Charset: col.Tp.GetCharset(),
					Elems:   col.Tp.GetElems(),
				})
		}
	}
	return in, false
}

func (c *TraverseCreateTableRowSizeTooLarge) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseCreateTableColsTp
type TraverseCreateTableColsTp struct {
	Table   string // 表名
	Charset string // 表字符集
	Cols    []process.LargePrefixIndexPartSpecification
}

func (c *TraverseCreateTableColsTp) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, node := range stmt.Options {
			switch node.Tp {
			case ast.TableOptionCharset:
				c.Charset = node.StrValue
			}
		}

		for _, col := range stmt.Cols {
			c.Cols = append(c.Cols,
				process.LargePrefixIndexPartSpecification{
					Column:  col.Name.Name.O,
					Tp:      col.Tp.GetType(),
					Flen:    col.Tp.GetFlen(),
					Decimal: col.Tp.GetDecimal(),
					Charset: col.Tp.GetCharset(),
					Elems:   col.Tp.GetElems(),
				})
		}
	}
	return in, false
}

func (c *TraverseCreateTableColsTp) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
