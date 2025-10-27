package logics

import (
	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/process"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicRenameTable
func LogicAnalyzeTable(v *traverses.TraverseAnalyzeTable, r *controllers.RuleHint) {
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
		if msg, err := dao.CheckIfTableExists(table, r.DB); err != nil {
			r.Summary = append(r.Summary, msg)
		}
	}
}
