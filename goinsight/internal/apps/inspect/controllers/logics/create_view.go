/*
@Time    :   2022/07/06 10:12:14
@Author  :   zongfei.fu
@Desc    :   None
*/

package logics

import (
	"fmt"
	"strings"

	"goInsight/internal/apps/inspect/controllers"
	"goInsight/internal/apps/inspect/controllers/dao"
	"goInsight/internal/apps/inspect/controllers/traverses"
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
		if err, msg := dao.DescTable(v.View, r.DB); err == nil {
			newMsg := strings.Join([]string{msg, "【TiDB可以使用`CREATE OR REPLACE VIEW`语法】"}, "")
			r.Summary = append(r.Summary, newMsg)
			r.IsSkipNextStep = true
		}
	}
	for _, table := range v.Tables {
		// 检查除视图名外的表是否存在
		if v.View != table {
			if err, msg := dao.VerifyTable(table, r.DB); err != nil {
				r.Summary = append(r.Summary, msg)
				r.IsSkipNextStep = true
			}
		}
	}
}
