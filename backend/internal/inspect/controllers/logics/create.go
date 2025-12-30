package logics

import (
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/process"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicCreateTableIsExist
func LogicCreateTableIsExist(v *traverses.TraverseCreateTableIsExist, r *controllers.RuleHint) {
	// 建表前置校验：表已存在时，本次 CREATE 会失败，直接终止后续规则检查。
	if msg, err := dao.CheckIfTableExists(v.Table, r.DB); err == nil {
		r.Warn(msg)
		r.IsBreak = true
	}
}

// LogicCreateTableAs
func LogicCreateTableAs(v *traverses.TraverseCreateTableAs, r *controllers.RuleHint) {
	if v.IsCreateAs {
		// `CREATE TABLE ... AS SELECT ...`：通常会绕过字段/索引/约束等细粒度规范，默认需要显式放开。
		if !r.InspectParams.ENABLE_CREATE_TABLE_AS {
			r.Warn(fmt.Sprintf("禁止使用 `CREATE TABLE ... AS SELECT ...`（表`%s`）", v.Table))
			r.IsBreak = true
		}
	}
}

// LogicCreateTableLike
func LogicCreateTableLike(v *traverses.TraverseCreateTableLike, r *controllers.RuleHint) {
	if v.IsCreateLike {
		if !r.InspectParams.ENABLE_CREATE_TABLE_LIKE {
			r.Warn(fmt.Sprintf("禁止使用 `CREATE TABLE ... LIKE ...`（表`%s`）", v.Table))
			r.IsBreak = true
		}
	}
}

// LogicCreateTableOptions
func LogicCreateTableOptions(v *traverses.TraverseCreateTableOptions, r *controllers.RuleHint) {
	v.Type = "create"
	v.TableOptions.InspectParams = r.InspectParams
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
			r.Warn(err.Error())
		}
	}
}

// LogicCreateTablePrimaryKey
func LogicCreateTablePrimaryKey(v *traverses.TraverseCreateTablePrimaryKey, r *controllers.RuleHint) {
	// 主键规则：没有主键会导致 UPDATE/DELETE 定位困难、复制/分库分表不友好。
	if r.InspectParams.CHECK_TABLE_PRIMARY_KEY {
		if len(v.PrimaryKeys) == 0 {
			r.Warn(fmt.Sprintf("表`%s`必须定义主键（建议使用 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT）", v.Table))
		}
		if len(v.PrimaryKeys) > 1 {
			r.Warn(fmt.Sprintf("表`%s`主键定义异常：只能定义一个主键", v.Table))
		}
	}
	// 检查主键列的类型/属性是否符合规范。
	for _, item := range v.PrimaryKeys {
		var p process.PrimaryKey = item
		p.InspectParams = r.InspectParams
		fns := []func() error{
			p.CheckBigint,
			p.CheckUnsigned,
			p.CheckAutoIncrement,
			p.CheckNotNull,
		}
		for _, fn := range fns {
			if err := fn(); err != nil {
				r.Warn(err.Error())
			}
		}
	}
}

// LogicCreateTableConstraint
func LogicCreateTableConstraint(v *traverses.TraverseCreateTableConstraint, r *controllers.RuleHint) {
	if !r.InspectParams.ENABLE_FOREIGN_KEY && v.IsForeignKey {
		// 外键会带来额外的锁/性能开销与跨系统耦合，默认不允许。
		r.Warn(fmt.Sprintf("表`%s`禁止定义外键", v.Table))
	}
}

// LogicCreateTableAuditCols
func LogicCreateTableAuditCols(v *traverses.TraverseCreateTableAuditCols, r *controllers.RuleHint) {
	if r.InspectParams.CHECK_TABLE_AUDIT_TYPE_COLUMNS {
		// 审计字段：建议至少包含“创建时间/更新时间”两类字段。
		// 这里不强制字段名，但会检查是否具备典型的默认值/自动更新时间语义。
		var colsOptionsArray []string
		for _, item := range v.AuditCols {
			for _, value := range item {
				colsOptionsArray = append(colsOptionsArray, value)
			}
		}
		if !utils.IsContain(colsOptionsArray, "DEFAULT CURRENT_TIMESTAMP") {
			r.Warn(fmt.Sprintf("表`%s`缺少创建时间类审计字段：建议包含 `DEFAULT CURRENT_TIMESTAMP`（如 `created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP`）", v.Table))
		}
		if !utils.IsContain(colsOptionsArray, "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP") {
			r.Warn(fmt.Sprintf("表`%s`缺少更新时间类审计字段：建议包含 `DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`（如 `updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`）", v.Table))
		}
	}
}

// LogicCreateTableColsOptions
func LogicCreateTableColsOptions(v *traverses.TraverseCreateTableColsOptions, r *controllers.RuleHint) {
	for _, col := range v.Cols {
		col.InspectParams = r.InspectParams
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
				r.Warn(err.Error())
			}
		}
	}
}

// LogicCreateTableColsRepeatDefine
func LogicCreateTableColsRepeatDefine(v *traverses.TraverseCreateTableColsRepeatDefine, r *controllers.RuleHint) {
	// 列名重复通常意味着建表语句拷贝/合并出错。
	if ok, data := utils.IsRepeat(v.Cols); ok {
		r.Warn(fmt.Sprintf("表`%s`存在重复列名：`%s`", v.Table, strings.Join(data, ",")))
	}
}

// LogicCreateTableColsCharset
func LogicCreateTableColsCharset(v *traverses.TraverseCreateTableColsCharset, r *controllers.RuleHint) {
	// 列字符集检查
	if r.InspectParams.CHECK_COLUMN_CHARSET {
		if len(v.Cols) > 0 {
			if err := v.CheckColumn(); err != nil {
				r.Warn(err.Error())
			}
		}
	}
}

// LogicCreateTableIndexesPrefix
func LogicCreateTableIndexesPrefix(v *traverses.TraverseCreateTableIndexesPrefix, r *controllers.RuleHint) {
	// 检查唯一索引前缀、如唯一索引必须以uniq_为前缀
	var indexPrefixCheck process.IndexPrefix = v.Prefix
	indexPrefixCheck.InspectParams = r.InspectParams
	if r.InspectParams.CHECK_UNIQ_INDEX_PREFIX {
		if err := indexPrefixCheck.CheckUniquePrefix(); err != nil {
			r.Warn(err.Error())
		}
	}
	// 检查二级索引前缀、如二级索引必须以idx_为前缀
	if r.InspectParams.CHECK_SECONDARY_INDEX_PREFIX {
		if err := indexPrefixCheck.CheckSecondaryPrefix(); err != nil {
			r.Warn(err.Error())
		}
	}
	// 检查全文索引前缀、如全文索引必须以full_为前缀
	if r.InspectParams.CHECK_FULLTEXT_INDEX_PREFIX {
		if err := indexPrefixCheck.CheckFulltextPrefix(); err != nil {
			r.Warn(err.Error())
		}
	}
}

// LogicCreateTableIndexesCount
func LogicCreateTableIndexesCount(v *traverses.TraverseCreateTableIndexesCount, r *controllers.RuleHint) {
	// 检查二级索引数量
	var indexNumberCheck process.IndexNumber = v.Number
	indexNumberCheck.InspectParams = r.InspectParams
	if err := indexNumberCheck.CheckSecondaryIndexesNum(); err != nil {
		r.Warn(err.Error())
	}
	if err := indexNumberCheck.CheckPrimaryKeyColsNum(); err != nil {
		r.Warn(err.Error())
	}
}

// LogicCreateTableIndexesRepeatDefine
func LogicCreateTableIndexesRepeatDefine(v *traverses.TraverseCreateTableIndexesRepeatDefine, r *controllers.RuleHint) {
	// 索引名重复会导致建表失败；字段组合重复会造成冗余索引与写入成本上升。
	if ok, data := utils.IsRepeat(v.Indexes); ok {
		r.Warn(fmt.Sprintf("表`%s`存在重复索引：`%s`", v.Table, strings.Join(data, ",")))
	}
}

// LogicCreateTableRedundantIndexes
func LogicCreateTableRedundantIndexes(v *traverses.TraverseCreateTableRedundantIndexes, r *controllers.RuleHint) {
	if r.InspectParams.ENABLE_REDUNDANT_INDEX {
		return
	}
	// 冗余索引会增加写入成本、占用更多存储，且通常不会带来查询收益。
	// 典型场景：同列重复索引、(a) 与 (a,b) 这种前缀覆盖。
	var redundantIndexCheck process.RedundantIndex = v.Redundant
	if err := redundantIndexCheck.CheckRepeatCols(); err != nil {
		r.Warn(err.Error())
	}
	if err := redundantIndexCheck.CheckRepeatColsWithDiffIndexes(); err != nil {
		r.Warn(err.Error())
	}
	if err := redundantIndexCheck.CheckRedundantColsWithDiffIndexes(); err != nil {
		r.Warn(err.Error())
	}
}

// LogicCreateTableDisabledIndexes
func LogicCreateTableDisabledIndexes(v *traverses.TraverseCreateTableDisabledIndexes, r *controllers.RuleHint) {
	// BLOB/TEXT 等大字段不允许建索引（存储/性能成本高，且大多场景不可用）。
	var indexTypesCheck process.DisabledIndexes = v.DisabledIndexes
	if err := indexTypesCheck.Check(); err != nil {
		r.Warn(err.Error())
	}
}

// LogicCreateTableInnodbLargePrefix
func LogicCreateTableInnodbLargePrefix(v *traverses.TraverseCreateTableInnodbLargePrefix, r *controllers.RuleHint) {
	var largePrefix process.LargePrefix = v.LargePrefix
	if err := largePrefix.Check(r.KV); err != nil {
		r.Warn(err.Error())
	}
}

// LogicCreateTableInnoDBRowSize
func LogicCreateTableInnoDBRowSize(v *traverses.TraverseCreateTableInnoDBRowSize, r *controllers.RuleHint) {
	var rowSize process.InnoDBRowSize = v.InnoDBRowSize
	if err := rowSize.Check(r.KV); err != nil {
		r.Warn(err.Error())
	}
}

// LogicCreateTableInnoDBRowFormat
func LogicCreateTableInnoDBRowFormat(v *traverses.TraverseCreateTableOptions, r *controllers.RuleHint) {
	// 行格式影响行存储与溢出页策略；这里只允许白名单中的 ROW_FORMAT。
	var rowFormat string = v.RowFormat
	if v.RowFormat == "DEFAULT" {
		rowFormat = r.KV.Get("innodbDefaultRowFormat").(string)
	}
	if !utils.IsContain(r.InspectParams.INNODB_ROW_FORMAT, rowFormat) {
		r.Warn(fmt.Sprintf("表`%s`的 ROW_FORMAT=%s 不在允许列表中（允许：%s）", v.Table, rowFormat, strings.Join(r.InspectParams.INNODB_ROW_FORMAT, ",")))
	}
}
