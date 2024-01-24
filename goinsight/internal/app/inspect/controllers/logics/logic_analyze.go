/*
@Time    :   2023/04/19 15:09:38
@Author  :   zongfei.fu
@Desc    :
*/

package logics

import (
	"goInsight/internal/app/inspect/controllers"
	"goInsight/internal/app/inspect/controllers/dao"
	"goInsight/internal/app/inspect/controllers/traverses"
	"sqlSyntaxAudit/controllers/process"
)

// LogicRenameTable
func LogicAnalyzeTable(v *traverses.TraverseAnalyzeTable, r *controllers.Rule) {
	if v.IsMatch == 0 {
		return
	}
	dbVersionIns := process.DbVersion{Version: r.KV.Get("dbVersion").(string)}
	if !dbVersionIns.IsTiDB() {
		r.Summary = append(r.Summary, "仅允许TiDB提交Analyze table语法")
		return
	}
	// 表必须存在
	for _, table := range v.TableNames {
		if err, msg := dao.DescTable(table, r.DB); err != nil {
			r.Summary = append(r.Summary, msg)
		}
	}
}
