package logics

import (
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicCreateViewIsExist
func LogicCreateViewIsExist(v *traverses.TraverseCreateViewIsExist, r *controllers.RuleHint) {
	if !r.InspectParams.ENABLE_CREATE_VIEW {
		r.Warn(fmt.Sprintf("禁止创建视图`%s`", v.View))
		r.IsBreak = true
		return
	}
	if !v.OrReplace {
		// `CREATE VIEW`（非 OR REPLACE）要求视图不存在，否则执行会失败。
		if msg, err := dao.CheckIfTableExists(v.View, r.DB); err == nil {
			newMsg := strings.Join([]string{msg, "；建议在 TiDB 使用 `CREATE OR REPLACE VIEW` 更新视图"}, "")
			r.Warn(newMsg)
			r.IsBreak = true
		}
	}
	for _, table := range v.Tables {
		// 依赖表校验：引用的表/视图必须存在（排除当前要创建的视图名本身）。
		if v.View != table {
			if msg, err := dao.CheckIfDatabaseExists(table, r.DB); err != nil {
				r.Warn(msg)
				r.IsBreak = true
			}
		}
	}
}
