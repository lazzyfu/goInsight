package logics

import (
	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/traverses"
)

// LogicCreateDatabaseIsExist
func LogicCreateDatabaseIsExist(v *traverses.TraverseCreateDatabaseIsExist, r *controllers.RuleHint) {
	if msg, err := dao.CheckIfDatabaseExists(v.Name, r.DB); err == nil {
		r.Warn("数据库已存在：" + msg)
		r.IsBreak = true
	}
}
