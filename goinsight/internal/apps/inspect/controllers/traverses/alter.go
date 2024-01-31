/*
@Time    :   2022/06/28 10:25:18
@Author  :   zongfei.fu
@Desc    :   遍历语法树,语法参考pingcap文档：https://github.com/pingcap/parser/blob/master/docs/quickstart.md
*/
package traverses

import (
	"goInsight/internal/apps/inspect/controllers/process"
	"goInsight/internal/pkg/utils"

	"github.com/pingcap/tidb/parser/ast"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

// TraverseAlterTableIsExist
type TraverseAlterTableIsExist struct {
	Table string // 表名
}

func (c *TraverseAlterTableIsExist) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
	}
	return in, false
}

func (c *TraverseAlterTableIsExist) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTiDBMerge
type TraverseAlterTiDBMerge struct {
	Table    string // 表名
	SpecsLen int    // 语句长度
}

func (c *TraverseAlterTiDBMerge) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.SpecsLen = len(stmt.Specs)
	}
	return in, false
}

func (c *TraverseAlterTiDBMerge) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableDropCols
type TraverseAlterTableDropColsOrIndexes struct {
	Table   string   // 表名
	IsMatch int      // 是否匹配当前规则
	Cols    []string // 列
	Indexes []string // 索引
}

func (c *TraverseAlterTableDropColsOrIndexes) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableDropColumn:
				c.IsMatch++
				c.Cols = append(c.Cols, spec.OldColumnName.Name.O)
			case ast.AlterTableDropIndex:
				c.IsMatch++
				c.Indexes = append(c.Indexes, spec.Name)
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableDropColsOrIndexes) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableDropTiDBColWithCoveredIndex
type TraverseAlterTableDropTiDBColWithCoveredIndex struct {
	Table   string   // 表名
	IsMatch int      // 是否匹配当前规则
	Cols    []string // 列
}

func (c *TraverseAlterTableDropTiDBColWithCoveredIndex) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableDropColumn:
				c.IsMatch++
				c.Cols = append(c.Cols, spec.OldColumnName.Name.O)
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableDropTiDBColWithCoveredIndex) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableOptions
type TraverseAlterTableOptions struct {
	IsMatch int // 是否匹配当前规则
	process.TableOptions
}

func (c *TraverseAlterTableOptions) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.TableOptions.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableOption:
				// 此处可根据需求实现更多方法
				c.IsMatch++
				for _, node := range spec.Options {
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
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableOptions) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableColCharset
type TraverseAlterTableColCharset struct {
	Table   string // 表名
	IsMatch int    // 是否匹配当前规则
	process.Charset
}

func (c *TraverseAlterTableColCharset) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn:
				c.IsMatch++
				for _, col := range spec.NewColumns {
					colCharset := col.Tp.GetCharset()
					colCollate := col.Tp.GetCollate()
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
		}
	}
	return in, false
}

func (c *TraverseAlterTableColCharset) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableAddColOptions
type TraverseAlterTableAddColOptions struct {
	Table   string // 表名
	IsMatch int    // 是否匹配当前规则
	Cols    []process.ColOptions
}

func (c *TraverseAlterTableAddColOptions) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns:
				c.IsMatch++
				var NotNullFlag, HasDefaultValue, DefaultIsNull, isSkipType, HasComment bool
				var DefaultValue interface{}
				for _, col := range spec.NewColumns {
					for _, opt := range col.Options {
						switch opt.Tp {
						case ast.ColumnOptionComment:
							HasComment = true
						case ast.ColumnOptionDefaultValue:
							HasDefaultValue = true
							if v, ok := opt.Expr.(*driver.ValueExpr); ok {
								// 有默认值，且为NULL，且有NOT NULL约束，如(not null default null)
								DefaultIsNull = v.Datum.IsNull()
								// 获取列的默认值
								DefaultValue = v.GetValue()
							}
							if f, ok := opt.Expr.(*ast.FuncCallExpr); ok {
								DefaultValue = f.FnName.L
							}
						case ast.ColumnOptionPrimaryKey, ast.ColumnOptionAutoIncrement:
							isSkipType = true
						case ast.ColumnOptionNotNull:
							NotNullFlag = true
						}
					}
					if !isSkipType {
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
		}
	}
	return in, false
}

func (c *TraverseAlterTableAddColOptions) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableAddPrimaryKey
type TraverseAlterTableAddPrimaryKey struct {
	Table   string   // 表名
	IsMatch int      // 是否匹配当前规则
	Cols    []string // 主键
}

func (c *TraverseAlterTableAddPrimaryKey) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			// add `id` int(10) unsigned primary key AUTO_INCREMENT"
			case ast.AlterTableAddColumns:
				c.IsMatch++
				for _, col := range spec.NewColumns {
					for _, op := range col.Options {
						if op.Tp == ast.ColumnOptionPrimaryKey {
							c.Cols = append(c.Cols, col.Name.Name.O)
						}
					}
				}
			// alter table follow_clue1 add primary key(`name`)
			case ast.AlterTableAddConstraint:
				c.IsMatch++
				if spec.Constraint.Tp == ast.ConstraintPrimaryKey {
					for _, col := range spec.Constraint.Keys {
						// 去重
						if !utils.IsContain(c.Cols, col.Column.Name.O) {
							c.Cols = append(c.Cols, col.Column.Name.O)
						}
					}
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableAddPrimaryKey) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableAddColRepeatDefine
type TraverseAlterTableAddColRepeatDefine struct {
	Table   string   // 表名
	IsMatch int      // 是否匹配当前规则
	Cols    []string // 列名
}

func (c *TraverseAlterTableAddColRepeatDefine) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns:
				c.IsMatch++
				for _, col := range spec.NewColumns {
					c.Cols = append(c.Cols, col.Name.Name.O)
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableAddColRepeatDefine) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableAddIndexPrefix
type TraverseAlterTableAddIndexPrefix struct {
	Table   string              // 表名
	IsMatch int                 // 是否匹配当前规则
	Prefix  process.IndexPrefix // 前缀
}

func (c *TraverseAlterTableAddIndexPrefix) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.Prefix.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			// 前缀
			if spec.Tp == ast.AlterTableAddConstraint {
				c.IsMatch++
				switch spec.Constraint.Tp {
				case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
					c.Prefix.UniqueKeys = append(c.Prefix.UniqueKeys, spec.Constraint.Name)
				case ast.ConstraintFulltext:
					c.Prefix.SecondaryKeys = append(c.Prefix.FulltextKeys, spec.Constraint.Name)
				case ast.ConstraintIndex, ast.ConstraintKey:
					c.Prefix.SecondaryKeys = append(c.Prefix.SecondaryKeys, spec.Constraint.Name)
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableAddIndexPrefix) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableAddIndexCount
type TraverseAlterTableAddIndexCount struct {
	Table   string // 表名
	IsMatch int    // 是否匹配当前规则
	Number  process.IndexNumber
}

func (c *TraverseAlterTableAddIndexCount) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.Number.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			// 索引数量
			if spec.Tp == ast.AlterTableAddConstraint {
				c.IsMatch++
				switch spec.Constraint.Tp {
				case ast.ConstraintPrimaryKey:
					c.Number.Keys = append(c.Number.Keys, process.IndexLen{Index: "PrimaryKey", Len: len(spec.Constraint.Keys)})
				case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
					c.Number.Number++
					c.Number.Keys = append(c.Number.Keys, process.IndexLen{Index: spec.Constraint.Name, Len: len(spec.Constraint.Keys)})
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableAddIndexCount) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableAddConstraint
type TraverseAlterTableAddConstraint struct {
	Table        string // 表名
	IsMatch      int    // 是否匹配当前规则
	IsForeignKey bool   // 是否定义了外键
}

func (c *TraverseAlterTableAddConstraint) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			if spec.Tp == ast.AlterTableAddConstraint {
				c.IsMatch++
				if spec.Constraint.Tp == ast.ConstraintForeignKey {
					c.IsForeignKey = true
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableAddConstraint) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableAddIndexRepeatDefine
type TraverseAlterTableAddIndexRepeatDefine struct {
	Table   string   // 表名
	IsMatch int      // 是否匹配当前规则
	Indexes []string // 索引
}

func (c *TraverseAlterTableAddIndexRepeatDefine) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			if spec.Tp == ast.AlterTableAddConstraint {
				c.IsMatch++
				switch spec.Constraint.Tp {
				case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintFulltext:
					c.Indexes = append(c.Indexes, spec.Constraint.Name)
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableAddIndexRepeatDefine) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableRedundantIndexes
type TraverseAlterTableRedundantIndexes struct {
	Table     string                 // 表名
	AddCols   []string               // add的列，用于检查alter table xxx add `col1` xxx,add index idx_col1(`col1`)
	DropCols  []string               // drop的列，用于检查alter table xxx drop `col2`,add index idx_col2(`col2`);
	IsMatch   int                    // 是否匹配当前规则
	Redundant process.RedundantIndex // 冗余索引
}

func (c *TraverseAlterTableRedundantIndexes) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.Redundant.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns:
				for _, col := range spec.NewColumns {
					c.AddCols = append(c.AddCols, col.Name.String())
				}
			case ast.AlterTableDropColumn:
				c.DropCols = append(c.DropCols, spec.OldColumnName.Name.O)
			case ast.AlterTableDropIndex:
				c.IsMatch++
				c.Redundant.IndexesCols = append(c.Redundant.IndexesCols, process.IndexColsMap{Index: spec.Name, Tag: "is_drop"})
			case ast.AlterTableAddConstraint:
				c.IsMatch++
				switch spec.Constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintFulltext:
					c.Redundant.Indexes = append(c.Redundant.Indexes, spec.Constraint.Name)
					var idxColsMap []string
					for _, v := range spec.Constraint.Keys {
						idxColsMap = append(idxColsMap, v.Column.Name.L)
					}
					c.Redundant.IndexesCols = append(c.Redundant.IndexesCols, process.IndexColsMap{Index: spec.Constraint.Name, Tag: "is_add", Cols: idxColsMap})
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableRedundantIndexes) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableDisabledIndexes
type TraverseAlterTableDisabledIndexes struct {
	Table           string                  // 表名
	IsMatch         int                     // 是否匹配当前规则
	DisabledIndexes process.DisabledIndexes // 禁止创建索引的列类型
}

func (c *TraverseAlterTableDisabledIndexes) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		c.DisabledIndexes.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			// 索引数量
			if spec.Tp == ast.AlterTableAddConstraint {
				c.IsMatch++
				var idxColsMap []string
				for _, v := range spec.Constraint.Keys {
					idxColsMap = append(idxColsMap, v.Column.Name.L)
				}
				c.DisabledIndexes.IndexesCols = append(c.DisabledIndexes.IndexesCols, process.IndexColsMap{Index: spec.Constraint.Name, Cols: idxColsMap})
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableDisabledIndexes) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableModifyColOptions
type TraverseAlterTableModifyColOptions struct {
	Table   string // 表名
	IsMatch int    // 是否匹配当前规则
	Cols    []process.ColOptions
}

func (c *TraverseAlterTableModifyColOptions) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			if spec.Tp == ast.AlterTableModifyColumn {
				c.IsMatch++
				var NotNullFlag, HasDefaultValue, DefaultIsNull, isSkipType, HasComment bool
				var DefaultValue interface{}
				for _, col := range spec.NewColumns {
					for _, opt := range col.Options {
						switch opt.Tp {
						case ast.ColumnOptionComment:
							HasComment = true
						case ast.ColumnOptionDefaultValue:
							HasDefaultValue = true
							if v, ok := opt.Expr.(*driver.ValueExpr); ok {
								// 有默认值，且为NULL，且有NOT NULL约束，如(not null default null)
								DefaultIsNull = v.Datum.IsNull()
								// 获取列的默认值
								DefaultValue = v.GetValue()
							}
							if f, ok := opt.Expr.(*ast.FuncCallExpr); ok {
								DefaultValue = f.FnName.L
							}
						case ast.ColumnOptionPrimaryKey, ast.ColumnOptionAutoIncrement:
							isSkipType = true
						case ast.ColumnOptionNotNull:
							NotNullFlag = true
						}
					}
					if !isSkipType {
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
		}
	}
	return in, false
}

func (c *TraverseAlterTableModifyColOptions) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableChangeColOptions
type TraverseAlterTableChangeColOptions struct {
	Table   string // 表名
	IsMatch int    // 是否匹配当前规则
	Cols    []process.ColOptions
}

func (c *TraverseAlterTableChangeColOptions) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			if spec.Tp == ast.AlterTableChangeColumn {
				c.IsMatch++
				var NotNullFlag, HasDefaultValue, DefaultIsNull, isSkipType, HasComment bool
				var DefaultValue interface{}
				for _, col := range spec.NewColumns {
					for _, opt := range col.Options {
						switch opt.Tp {
						case ast.ColumnOptionComment:
							HasComment = true
						case ast.ColumnOptionDefaultValue:
							HasDefaultValue = true
							if v, ok := opt.Expr.(*driver.ValueExpr); ok {
								// 有默认值，且为NULL，且有NOT NULL约束，如(not null default null)
								DefaultIsNull = v.Datum.IsNull()
								// 获取列的默认值
								DefaultValue = v.GetValue()
							}
							if f, ok := opt.Expr.(*ast.FuncCallExpr); ok {
								DefaultValue = f.FnName.L
							}
						case ast.ColumnOptionPrimaryKey, ast.ColumnOptionAutoIncrement:
							isSkipType = true
						case ast.ColumnOptionNotNull:
							NotNullFlag = true
						}
					}
					if !isSkipType {
						c.Cols = append(c.Cols,
							process.ColOptions{
								Table:           c.Table,
								OldColumn:       spec.OldColumnName.Name.O,
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
		}
	}
	return in, false
}

func (c *TraverseAlterTableChangeColOptions) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableRenameIndex
type TraverseAlterTableRenameIndex struct {
	Table   string // 表名
	IsMatch int    // 是否匹配当前规则
	Indexes []struct {
		OldIndex string // old_index_name
		NewIndex string // new_index_name
	}
}

func (c *TraverseAlterTableRenameIndex) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			if spec.Tp == ast.AlterTableRenameIndex {
				c.IsMatch++
				c.Indexes = append(c.Indexes, struct {
					OldIndex string
					NewIndex string
				}{OldIndex: spec.FromKey.O, NewIndex: spec.ToKey.O})
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableRenameIndex) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableRenameTblName
type TraverseAlterTableRenameTblName struct {
	Table      string // 表名
	IsMatch    int    // 是否匹配当前规则
	NewTblName string // 新表名
}

func (c *TraverseAlterTableRenameTblName) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			if spec.Tp == ast.AlterTableRenameTable {
				c.IsMatch++
				c.NewTblName = spec.NewTable.Name.O
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableRenameTblName) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableShowCreateTableGetCols
type TraverseAlterTableShowCreateTableGetCols struct {
	Table       string   // 表名
	PrimaryKeys []string // 主键
	Cols        []string // 列
	Indexes     []string // 索引
}

func (c *TraverseAlterTableShowCreateTableGetCols) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.CreateTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		// 获取所有的列
		for _, col := range stmt.Cols {
			c.Cols = append(c.Cols, col.Name.Name.O)
		}
		// 获取主键
		for _, cons := range stmt.Constraints {
			if cons.Tp == ast.ConstraintPrimaryKey {
				for _, col := range cons.Keys {
					c.PrimaryKeys = append(c.PrimaryKeys, col.Column.Name.O)
				}
			}
		}
		// 获取索引
		for _, cons := range stmt.Constraints {
			switch cons.Tp {
			case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintIndex:
				c.Indexes = append(c.Indexes, cons.Name)
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableShowCreateTableGetCols) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableInnodbLargePrefix
type TraverseAlterTableInnodbLargePrefix struct {
	IsMatch     int // 是否匹配当前规则
	LargePrefix process.LargePrefix
}

func (c *TraverseAlterTableInnodbLargePrefix) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.LargePrefix.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddConstraint:
				c.IsMatch++
				switch spec.Constraint.Tp {
				case ast.ConstraintKey, ast.ConstraintIndex, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintFulltext:
					var LargePrefixIndexColsMap process.LargePrefixIndexColsMap
					LargePrefixIndexColsMap.Name = spec.Constraint.Name
					for _, col := range spec.Constraint.Keys {
						LargePrefixIndexColsMap.Keys = append(LargePrefixIndexColsMap.Keys,
							process.LargePrefixIndexPartSpecification{Column: col.Column.Name.L, Ilen: col.Length})
					}
					c.LargePrefix.LargePrefixIndexColsMaps = append(c.LargePrefix.LargePrefixIndexColsMaps, LargePrefixIndexColsMap)
				}
			}
		}
	}
	return in, false
}

func (c *TraverseAlterTableInnodbLargePrefix) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// TraverseAlterTableRowSizeTooLarge
type TraverseAlterTableRowSizeTooLarge struct {
	Table           string // 表名
	IsMatch         int    // 是否匹配当前规则
	RowSizeColsMaps []process.RowSizeTooLargePartSpecification
}

func (c *TraverseAlterTableRowSizeTooLarge) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.AlterTableStmt); ok {
		c.Table = stmt.Table.Name.String()
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn:
				c.IsMatch++
				for _, col := range spec.NewColumns {
					c.RowSizeColsMaps = append(c.RowSizeColsMaps,
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
		}
	}
	return in, false
}

func (c *TraverseAlterTableRowSizeTooLarge) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
