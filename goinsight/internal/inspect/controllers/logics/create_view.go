/*
@Time    :   2022/07/06 10:12:14
@Author  :   xff
@Desc    :   None
*/

package logics

import (
	"fmt"
	"strings"

	"goInsight/internal/inspect/controllers"
	"goInsight/internal/inspect/controllers/dao"
	"goInsight/internal/inspect/controllers/traverses"
)

// LogicCreateViewIsExist
func LogicCreateViewIsExist(v *traverses.TraverseCreateViewIsExist, r *controllers.RuleHint) {
	if !r.InspectParams.ENABLE_CREATE_VIEW {
		r.Summary = append(r.Summary, fmt.Sprintf("不允许创建视图`%s`", v.View))
		r.IsSkipNextStep = true
		return
	}
	if !v.OrReplace {
		// create view，需要确保视图不存在
		if msg, err := dao.CheckIfTableExists(v.View, r.DB); err == nil {
			newMsg := strings.Join([]string{msg, "【TiDB可以使用`CREATE OR REPLACE VIEW`语法】"}, "")
			r.Summary = append(r.Summary, newMsg)
			r.IsSkipNextStep = true
		}
	}
	for _, table := range v.Tables {
		// 检查除视图名外的表是否存在
		if v.View != table {
			if msg, err := dao.CheckIfDatabaseExists(table, r.DB); err != nil {
				r.Summary = append(r.Summary, msg)
				r.IsSkipNextStep = true
			}
		}
	}
}
