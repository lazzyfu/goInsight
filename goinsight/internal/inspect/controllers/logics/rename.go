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
		r.Summary = append(r.Summary, "不允许执行RENAME TABLE操作")
		return
	}
	// 禁止审核指定的表
	if len(r.InspectParams.DISABLE_AUDIT_DDL_TABLES) > 0 {
		for _, item := range r.InspectParams.DISABLE_AUDIT_DDL_TABLES {
			for _, t := range v.Tables {
				if item.DB == r.DB.Database && utils.IsContain(item.Tables, t.OldTable) {
					r.Summary = append(r.Summary, fmt.Sprintf("表`%s`.`%s`被限制进行DDL语法审核，原因: %s", r.DB.Database, t.OldTable, item.Reason))
				}
			}
		}
	}
	var oldTables []string
	// 旧表必须存在
	for _, t := range v.Tables {
		if msg, err := dao.CheckIfTableExists(t.OldTable, r.DB); err != nil {
			r.Summary = append(r.Summary, msg)
		} else {
			oldTables = append(oldTables, t.OldTable)
		}
	}
	// 新表不能存在
	for _, t := range v.Tables {
		if len(oldTables) > 0 && utils.IsContain(oldTables, t.NewTable) {
			continue
		}
		if msg, err := dao.CheckIfTableExists(t.NewTable, r.DB); err == nil {
			r.Summary = append(r.Summary, msg)
		}
	}
}
