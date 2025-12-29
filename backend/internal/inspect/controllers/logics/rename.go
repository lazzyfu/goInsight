package logics

import (
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
