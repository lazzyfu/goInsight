package logics

import (
	"fmt"

	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicDropTable
func LogicDropTable(v *traverses.TraverseDropTable, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if v.IsHasDropTable {
		if !r.InspectParams.ENABLE_DROP_TABLE {
			r.Warn(fmt.Sprintf("禁止执行 `DROP TABLE`：%v", v.Tables))
			return
		}
		// 语句校验：目标表必须存在。
		for _, table := range v.Tables {
			if msg, err := dao.CheckIfTableExists(table, r.DB); err != nil {
				r.Warn(msg)
			}
		}
	}
}

// LogicTruncateTable
func LogicTruncateTable(v *traverses.TraverseTruncateTable, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if v.IsHasTruncateTable {
		if !r.InspectParams.ENABLE_TRUNCATE_TABLE {
			r.Warn(fmt.Sprintf("禁止执行 `TRUNCATE TABLE`：`%s`", v.Table))
			return
		}
		// 语句校验：目标表必须存在。
		if msg, err := dao.CheckIfTableExists(v.Table, r.DB); err != nil {
			r.Warn(msg)
		}
	}
}
