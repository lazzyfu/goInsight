package logics

import (
	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/process"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicAnalyzeTable
// 审核 ANALYZE TABLE：目前仅允许在 TiDB 环境执行（避免在 MySQL 上触发不可控的统计信息更新/性能影响）。
func LogicAnalyzeTable(v *traverses.TraverseAnalyzeTable, r *controllers.RuleHint) {
	if v.IsMatch == 0 {
		return
	}
	dbVersionIns := process.DbVersion{Version: r.KV.Get("dbVersion").(string)}
	if !dbVersionIns.IsTiDB() {
		r.Warn("仅允许在 TiDB 环境提交 `ANALYZE TABLE` 语句")
		return
	}
	// 审核前先校验目标表是否存在：避免对不存在的表提交 analyze 导致执行失败。
	for _, table := range v.TableNames {
		if msg, err := dao.CheckIfTableExists(table, r.DB); err != nil {
			r.Warn(msg)
		}
	}
}
