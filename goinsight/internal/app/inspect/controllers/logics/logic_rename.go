/*
@Time    :   2022/08/25 16:42:48
@Author  :   zongfei.fu
@Desc    :   None
*/

package logics

import (
	"fmt"
	"sqlSyntaxAudit/common/utils"
)

// LogicRenameTable
func LogicRenameTable(v *TraverseRenameTable, r *Rule) {
	if v.IsMatch == 0 {
		return
	}
	if !r.AuditConfig.ENABLE_RENAME_TABLE_NAME {
		r.Summary = append(r.Summary, "不允许执行RENAME TABLE操作")
		return
	}
	// 禁止审核指定的表
	if len(r.AuditConfig.DISABLE_AUDIT_DDL_TABLES) > 0 {
		for _, item := range r.AuditConfig.DISABLE_AUDIT_DDL_TABLES {
			for _, t := range v.tables {
				if item.DB == r.DB.Database && utils.IsContain(item.Tables, t.OldTable) {
					r.Summary = append(r.Summary, fmt.Sprintf("表`%s`.`%s`被限制进行DDL语法审核，原因: %s", r.DB.Database, t.OldTable, item.Reason))
				}
			}
		}
	}
	var oldTables []string
	// 旧表必须存在
	for _, t := range v.tables {
		if err, msg := DescTable(t.OldTable, r.DB); err != nil {
			r.Summary = append(r.Summary, msg)
		} else {
			oldTables = append(oldTables, t.OldTable)
		}
	}
	// 新表不能存在
	for _, t := range v.tables {
		if len(oldTables) > 0 && utils.IsContain(oldTables, t.NewTable) {
			continue
		}
		if err, msg := DescTable(t.NewTable, r.DB); err == nil {
			r.Summary = append(r.Summary, msg)
		}
	}
}
