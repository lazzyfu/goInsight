package logics

import (
	"fmt"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicRenameTable
func LogicRenameTable(v *traverses.TraverseRenameTable, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if !r.InspectParams.ENABLE_RENAME_TABLE_NAME {
		r.Warn("禁止执行 `RENAME TABLE`")
		return
	}
	// 黑名单表：禁止对指定表执行 DDL（常用于核心表保护）。
	if len(r.InspectParams.DISABLE_AUDIT_DDL_TABLES) > 0 {
		for _, item := range r.InspectParams.DISABLE_AUDIT_DDL_TABLES {
			for _, t := range v.Tables {
				if item.DB == r.DB.Database && utils.IsContain(item.Tables, t.OldTable) {
					r.Warn(fmt.Sprintf("禁止对表`%s`.`%s`执行 DDL 审核：%s", r.DB.Database, t.OldTable, item.Reason))
				}
			}
		}
	}
	var oldTables []string
	// 旧表必须存在：避免把不存在的表重命名。
	for _, t := range v.Tables {
		if msg, err := dao.CheckIfTableExists(t.OldTable, r.DB); err != nil {
			r.Warn(msg)
		} else {
			oldTables = append(oldTables, t.OldTable)
		}
	}
	// 新表名不能已存在：避免覆盖现有对象导致数据风险。
	for _, t := range v.Tables {
		if len(oldTables) > 0 && utils.IsContain(oldTables, t.NewTable) {
			continue
		}
		if msg, err := dao.CheckIfTableExists(t.NewTable, r.DB); err == nil {
			r.Warn(msg)
		}
	}
}
