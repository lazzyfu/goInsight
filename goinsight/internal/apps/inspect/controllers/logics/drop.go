/*
@Time    :   2022/07/06 10:12:27
@Author  :   xff
@Desc    :   None
*/

package logics

import (
	"fmt"
	"goInsight/internal/apps/inspect/controllers"
	"goInsight/internal/apps/inspect/controllers/dao"
	"goInsight/internal/apps/inspect/controllers/traverses"
	"goInsight/internal/pkg/utils"
)

// LogicDropTable
func LogicDropTable(v *traverses.TraverseDropTable, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	if v.IsHasDropTable {
		if !r.InspectParams.ENABLE_DROP_TABLE {
			r.Summary = append(r.Summary, fmt.Sprintf("禁止DROP[表%s]", v.Tables))
			return
		}
		// 禁止审核指定的表
		if len(r.InspectParams.DISABLE_AUDIT_DDL_TABLES) > 0 {
			for _, item := range r.InspectParams.DISABLE_AUDIT_DDL_TABLES {
				for _, table := range v.Tables {
					if item.DB == r.DB.Database && utils.IsContain(item.Tables, table) {
						r.Summary = append(r.Summary, fmt.Sprintf("表`%s`.`%s`被限制进行DDL语法审核，原因: %s", r.DB.Database, table, item.Reason))
					}
				}
			}
		}
		// 检查表是否存在
		for _, table := range v.Tables {
			if err, msg := dao.DescTable(table, r.DB); err != nil {
				r.Summary = append(r.Summary, msg)
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
			r.Summary = append(r.Summary, fmt.Sprintf("禁止TRUNCATE[表%s]", v.Table))
			return
		}
		// 禁止审核指定的表
		if len(r.InspectParams.DISABLE_AUDIT_DDL_TABLES) > 0 {
			for _, item := range r.InspectParams.DISABLE_AUDIT_DDL_TABLES {
				if item.DB == r.DB.Database && utils.IsContain(item.Tables, v.Table) {
					r.Summary = append(r.Summary, fmt.Sprintf("表`%s`.`%s`被限制进行DDL语法审核，原因: %s", r.DB.Database, v.Table, item.Reason))
				}
			}
		}
		// 检查表是否存在
		if err, msg := dao.DescTable(v.Table, r.DB); err != nil {
			r.Summary = append(r.Summary, msg)
		}
	}
}
