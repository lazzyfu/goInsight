package logics

import (
	"fmt"

	"github.com/lazzyfu/goinsight/pkg/utils"

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
		// 黑名单表：禁止对指定表执行 DDL（常用于核心表保护）。
		if len(r.InspectParams.DISABLE_AUDIT_DDL_TABLES) > 0 {
			for _, item := range r.InspectParams.DISABLE_AUDIT_DDL_TABLES {
				for _, table := range v.Tables {
					if item.DB == r.DB.Database && utils.IsContain(item.Tables, table) {
						r.Warn(fmt.Sprintf("禁止对表`%s`.`%s`执行 DDL 审核：%s", r.DB.Database, table, item.Reason))
					}
				}
			}
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
		// 黑名单表：禁止对指定表执行 DDL（常用于核心表保护）。
		if len(r.InspectParams.DISABLE_AUDIT_DDL_TABLES) > 0 {
			for _, item := range r.InspectParams.DISABLE_AUDIT_DDL_TABLES {
				if item.DB == r.DB.Database && utils.IsContain(item.Tables, v.Table) {
					r.Warn(fmt.Sprintf("禁止对表`%s`.`%s`执行 DDL 审核：%s", r.DB.Database, v.Table, item.Reason))
				}
			}
		}
		// 语句校验：目标表必须存在。
		if msg, err := dao.CheckIfTableExists(v.Table, r.DB); err != nil {
			r.Warn(msg)
		}
	}
}
