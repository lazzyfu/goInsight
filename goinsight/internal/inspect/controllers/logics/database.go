package logics

import (
	"goInsight/internal/inspect/controllers"
	"goInsight/internal/inspect/controllers/dao"
	"goInsight/internal/inspect/controllers/traverses"
)

// LogicCreateDatabaseIsExist
func LogicCreateDatabaseIsExist(v *traverses.TraverseCreateDatabaseIsExist, r *controllers.RuleHint) {
	if msg, err := dao.CheckIfDatabaseExists(v.Name, r.DB); err == nil {
		r.Summary = append(r.Summary, msg)
		r.IsSkipNextStep = true
	}
}
