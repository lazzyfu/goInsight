/*
@Time    :   2022/07/06 10:12:14
@Author  :   zongfei.fu
@Desc    :   None
*/

package logics

import (
	"fmt"
	"sqlSyntaxAudit/common/utils"
	"sqlSyntaxAudit/controllers/process"
	"strings"
)

// LogicCreateTableIsExist
func LogicCreateTableIsExist(v *TraverseCreateTableIsExist, r *Rule) {
	// 检查表是否存在,如果表存在,skip下面的检查
	if err, msg := DescTable(v.Table, r.DB); err == nil {
		r.Summary = append(r.Summary, msg)
		r.IsSkipNextStep = true
	}
}

// LogicCreateTableAs
func LogicCreateTableAs(v *TraverseCreateTableAs, r *Rule) {
	if v.IsCreateAs {
		// 不深入检查AS后面的语法
		if !r.AuditConfig.ENABLE_CREATE_TABLE_AS {
			r.Summary = append(r.Summary, fmt.Sprintf("不允许使用create table as语法[表`%s`]", v.Table))
			r.IsSkipNextStep = true
		}
	}
}

// LogicCreateTableLike
func LogicCreateTableLike(v *TraverseCreateTableLike, r *Rule) {
	if v.IsCreateLike {
		if !r.AuditConfig.ENABLE_CREATE_TABLE_LIKE {
			r.Summary = append(r.Summary, fmt.Sprintf("不允许使用create table like语法[表`%s`]", v.Table))
			r.IsSkipNextStep = true
		}
	}
}

// LogicCreateTableOptions
func LogicCreateTableOptions(v *TraverseCreateTableOptions, r *Rule) {
	v.Type = "create"
	v.TableOptions.AuditConfig = r.AuditConfig
	fns := []func() error{
		v.CheckTableLength,
		v.CheckTableIdentifer,
		v.CheckTableIdentiferKeyword,
		v.CheckTableEngine,
		v.CheckTablePartition,
		v.CheckTableComment,
		v.CheckTableCharset,
		v.CheckTableAutoIncrementInitValue,
	}
	for _, fn := range fns {
		if err := fn(); err != nil {
			r.Summary = append(r.Summary, err.Error())
		}
	}
}

// LogicCreateTablePrimaryKey
func LogicCreateTablePrimaryKey(v *TraverseCreateTablePrimaryKey, r *Rule) {
	// 必须定义主键
	if r.AuditConfig.CHECK_TABLE_PRIMARY_KEY {
		if len(v.PrimaryKeys) == 0 {
			r.Summary = append(r.Summary, fmt.Sprintf("表`%s`必须定义主键", v.Table))
		}
		if len(v.PrimaryKeys) > 1 {
			r.Summary = append(r.Summary, fmt.Sprintf("表`%s`有且只能定义一个主键", v.Table))
		}
	}
	// 检查主键是否为bigint类型
	for _, item := range v.PrimaryKeys {
		var p process.PrimaryKey = item
		p.AuditConfig = r.AuditConfig
		fns := []func() error{
			p.CheckBigint,
			p.CheckUnsigned,
			p.CheckAutoIncrement,
			p.CheckNotNull,
		}
		for _, fn := range fns {
			if err := fn(); err != nil {
				r.Summary = append(r.Summary, err.Error())
			}
		}
	}
}

// LogicCreateTableConstraint
func LogicCreateTableConstraint(v *TraverseCreateTableConstraint, r *Rule) {
	if !r.AuditConfig.ENABLE_FOREIGN_KEY && v.IsForeignKey {
		// 禁止使用外键
		r.Summary = append(r.Summary, fmt.Sprintf("表`%s`禁止定义外键", v.Table))
	}
}

// LogicCreateTableAuditCols
func LogicCreateTableAuditCols(v *TraverseCreateTableAuditCols, r *Rule) {
	if r.AuditConfig.CHECK_TABLE_AUDIT_TYPE_COLUMNS {
		// 启用审计类型的字段, 必须定义2个审计字段, 字段名和注释名不做要求, 如:
		// `UPDATED_AT` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
		// `CREATED_AT` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
		var colsOptionsArray []string
		for _, item := range v.AuditCols {
			for _, value := range item {
				colsOptionsArray = append(colsOptionsArray, value)
			}
		}
		if !utils.IsContain(colsOptionsArray, "DEFAULT CURRENT_TIMESTAMP") {
			r.Summary = append(r.Summary, fmt.Sprintf("表`%s`未定义字段类型为%s的审计字段【例如：CREATED_AT datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'】", v.Table, "DEFAULT CURRENT_TIMESTAMP"))
		}
		if !utils.IsContain(colsOptionsArray, "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP") {
			r.Summary = append(r.Summary, fmt.Sprintf("表`%s`未定义字段类型为%s的审计字段【例如：UPDATED_AT datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'】", v.Table, "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"))
		}
	}
}

// LogicCreateTableColsOptions
func LogicCreateTableColsOptions(v *TraverseCreateTableColsOptions, r *Rule) {
	for _, col := range v.Cols {
		col.AuditConfig = r.AuditConfig
		fns := []func() error{
			col.CheckColumnLength,
			col.CheckColumnIdentifer,
			col.CheckColumnIdentiferKeyword,
			col.CheckColumnComment,
			col.CheckColumnCharToVarchar,
			col.CheckColumnMaxVarcharLength,
			col.CheckColumnFloatDouble,
			col.CheckColumnNotAllowedType,
			col.CheckColumnNotNull,
			col.CheckColumnDefaultValue,
		}
		for _, fn := range fns {
			if err := fn(); err != nil {
				r.Summary = append(r.Summary, err.Error())
			}
		}
	}
}

// LogicCreateTableColsRepeatDefine
func LogicCreateTableColsRepeatDefine(v *TraverseCreateTableColsRepeatDefine, r *Rule) {
	// 查找重复的列名
	if ok, data := utils.IsRepeat(v.Cols); ok {
		r.Summary = append(r.Summary, fmt.Sprintf("发现重复的列名`%s`[表`%s`]", strings.Join(data, ","), v.Table))
	}
}

// LogicCreateTableColsCharset
func LogicCreateTableColsCharset(v *TraverseCreateTableColsCharset, r *Rule) {
	// 列字符集检查
	if r.AuditConfig.CHECK_COLUMN_CHARSET {
		if len(v.Cols) > 0 {
			if err := v.CheckColumn(); err != nil {
				r.Summary = append(r.Summary, err.Error())
			}
		}
	}
}

// LogicCreateTableIndexesPrefix
func LogicCreateTableIndexesPrefix(v *TraverseCreateTableIndexesPrefix, r *Rule) {
	// 检查唯一索引前缀、如唯一索引必须以uniq_为前缀
	var indexPrefixCheck process.IndexPrefix = v.Prefix
	indexPrefixCheck.AuditConfig = r.AuditConfig
	if r.AuditConfig.CHECK_UNIQ_INDEX_PREFIX {
		if err := indexPrefixCheck.CheckUniquePrefix(); err != nil {
			r.Summary = append(r.Summary, err.Error())
		}
	}
	// 检查二级索引前缀、如二级索引必须以idx_为前缀
	if r.AuditConfig.CHECK_SECONDARY_INDEX_PREFIX {
		if err := indexPrefixCheck.CheckSecondaryPrefix(); err != nil {
			r.Summary = append(r.Summary, err.Error())
		}
	}
	// 检查全文索引前缀、如全文索引必须以full_为前缀
	if r.AuditConfig.CHECK_FULLTEXT_INDEX_PREFIX {
		if err := indexPrefixCheck.CheckFulltextPrefix(); err != nil {
			r.Summary = append(r.Summary, err.Error())
		}
	}
}

// LogicCreateTableIndexesCount
func LogicCreateTableIndexesCount(v *TraverseCreateTableIndexesCount, r *Rule) {
	// 检查二级索引数量
	var indexNumberCheck process.IndexNumber = v.Number
	indexNumberCheck.AuditConfig = r.AuditConfig
	if err := indexNumberCheck.CheckSecondaryIndexesNum(); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
	if err := indexNumberCheck.CheckPrimaryKeyColsNum(); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
}

// LogicCreateTableIndexesRepeatDefine
func LogicCreateTableIndexesRepeatDefine(v *TraverseCreateTableIndexesRepeatDefine, r *Rule) {
	// 查找重复的索引
	if ok, data := utils.IsRepeat(v.Indexes); ok {
		r.Summary = append(r.Summary, fmt.Sprintf("发现重复的索引`%s`[表`%s`]", strings.Join(data, ","), v.Table))
	}
}

// LogicCreateTableRedundantIndexes
func LogicCreateTableRedundantIndexes(v *TraverseCreateTableRedundantIndexes, r *Rule) {
	// 检查索引,建索引时,指定的列必须存在、索引中的列,不能重复、索引名不能重复
	// 不能有重复的索引,包括(索引名不同,字段相同；冗余索引,如(a),(a,b))
	var redundantIndexCheck process.RedundantIndex = v.Redundant
	if err := redundantIndexCheck.CheckRepeatCols(); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
	if err := redundantIndexCheck.CheckRepeatColsWithDiffIndexes(); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
	if err := redundantIndexCheck.CheckRedundantColsWithDiffIndexes(); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
}

// LogicCreateTableDisabledIndexes
func LogicCreateTableDisabledIndexes(v *TraverseCreateTableDisabledIndexes, r *Rule) {
	// BLOB/TEXT类型不能设置为索引
	var indexTypesCheck process.DisabledIndexes = v.DisabledIndexes
	if err := indexTypesCheck.Check(); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
}

// LogicCreateTableInnodbLargePrefix
func LogicCreateTableInnodbLargePrefix(v *TraverseCreateTableInnodbLargePrefix, r *Rule) {
	var largePrefix process.LargePrefix = v.LargePrefix
	if err := largePrefix.Check(r.KV); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
}

// LogicCreateTableRowSizeTooLarge
func LogicCreateTableRowSizeTooLarge(v *TraverseCreateTableRowSizeTooLarge, r *Rule) {
	var rowSizeTooLarge process.RowSizeTooLarge = v.RowSizeTooLarge
	if err := rowSizeTooLarge.Check(r.KV); err != nil {
		r.Summary = append(r.Summary, err.Error())
	}
}
