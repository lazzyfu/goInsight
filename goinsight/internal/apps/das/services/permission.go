/*
@Time    :   2023/03/23 16:45:28
@Author  :   zongfei.fu
@Desc    :   检查用户的库表权限
*/

package services

import (
	"errors"
	"fmt"
	"goInsight/global"
	"goInsight/internal/apps/das/parser"
	"goInsight/internal/pkg/utils"
	"strings"

	"github.com/google/uuid"
)

type CheckUserPerm struct {
	UserName   string
	InstanceID uuid.UUID
	Tables     []parser.Table
}

func (p CheckUserPerm) checkHasAllTablePerms() bool {
	// 表das_user_table_permissions中没有找到username=? and instance_id=? and `schema`=?的记录
	for _, t := range p.Tables {
		type countResult struct {
			Count int
		}
		var result countResult
		global.App.DB.Table("`insight_das_user_table_permissions`").
			Select("count(*) as count").
			Where("username=? and instance_id=? and `schema`=?", p.UserName, p.InstanceID, t.Schema).
			Scan(&result)
		if result.Count == 0 {
			return true
		}
	}
	return false
}

func (p CheckUserPerm) checkOnlyHasAllowTablePerms() []string {
	// 仅有allow规则的记录
	var errMsg []string
	for _, t := range p.Tables {
		type ruleResult struct {
			Table string
		}
		var result []ruleResult
		global.App.DB.Table("`insight_das_user_table_permissions`").
			Select("*").
			Where("username=? and instance_id=? and `schema`=? and rule='allow'", p.UserName, p.InstanceID, t.Schema).
			Scan(&result)
		if len(result) > 0 {
			// 结果集大于0，表示有allow规则的记录
			var tmpTables []string
			for _, r := range result {
				tmpTables = append(tmpTables, r.Table)
			}
			// 判断提取的表名t1是否在匹配出来的表集合里面["work_days", "t1"]
			if !utils.IsContain(tmpTables, t.Table) {
				errMsg = append(errMsg, fmt.Sprintf("您没有表`%s`.`%s`权限", t.Schema, t.Table))
			}
		}
	}
	return errMsg
}

func (p CheckUserPerm) checkOnlyHasDenyTablePerms() []string {
	var errMsg []string
	for _, t := range p.Tables {
		type ruleResult struct {
			Table string
		}
		var result []ruleResult
		global.App.DB.Table("`insight_das_user_table_permissions`").
			Select("*").
			Where("username=? and instance_id=? and `schema`=? and rule='deny'", p.UserName, p.InstanceID, t.Schema).
			Where("not exists(select * from insight_das_user_table_permissions where username=? and instance_id=? and `schema`=? and rule='allow')", p.UserName, p.InstanceID, t.Schema).
			Scan(&result)
		if len(result) > 0 {
			// 结果集大于0，表示仅有deny规则的记录
			var tmpTables []string
			for _, r := range result {
				tmpTables = append(tmpTables, r.Table)
			}
			// 判断提取的表名t1是否在匹配出来的表集合里面["work_days", "t1"]
			if utils.IsContain(tmpTables, t.Table) {
				errMsg = append(errMsg, fmt.Sprintf("您没有表`%s`.`%s`权限", t.Schema, t.Table))
			}
		}
	}
	return errMsg
}

// 检查是否有库访问权限
func (p CheckUserPerm) HasSchemaPerms() error {
	if err := func() error {
		for _, t := range p.Tables {
			// information_schema库默认有权限
			if strings.ToLower(t.Schema) == "information_schema" {
				continue
			}
			// 鉴权
			type countResult struct {
				Count int
			}
			var result countResult
			global.App.DB.Table("`insight_das_user_schema_permissions` p").
				Select("count(*) as count").
				Joins("join `insight_db_schemas` s on p.instance_id = s.instance_id and p.`schema` = s.`schema`").
				Where("s.is_deleted = 0 and p.username=? and p.instance_id=? and p.`schema`=?", p.UserName, p.InstanceID, t.Schema).
				Scan(&result)
			if result.Count == 0 {
				return fmt.Errorf("您没有库`%s`的访问权限", t.Schema)
			}
		}
		return nil
	}(); err != nil {
		return err
	}
	return nil
}

func (p CheckUserPerm) HasTablePerms() error {
	/*
		优先级：allow > deny
		前提：insight_das_user_schema_permissions表有zhangsan访问test库的授权记录
		Case分析：
			Case1:
				Desc:
					表das_user_table_permissions中用户zhangsan的记录为空
				Result:
					用户zhangsan可以访问当前拥有test库所有的表权限
			Case2：
				Desc:
					表das_user_table_permissions中用户zhangsan有如下记录
					zhangsan | test | work_days | allow
					zhangsan | test | st_week   | allow
					zhangsan | test | st_total  | deny
				Result:
					当allow和deny规则并存时，仅allow规则生效，deny规则不生效
					用户zhangsan可以访问test.work_days、test.st_week表，deny这条规则不生效，且不用添加
					如果不希望访问表st_week，删除这条记录（zhangsan | test | st_week | allow）
			Case3：
				Desc:
					表das_user_table_permissions中用户zhangsan有如下记录
					zhangsan | test | work_days | deny
					zhangsan | test | st_total  | deny
				Result:
					用户zhangsan不可以访问test.work_days、test.st_total，但可以访问test库其他的表
	*/
	// Case1
	if ok := p.checkHasAllTablePerms(); ok {
		return nil
	}
	// Case2
	if errMsg := p.checkOnlyHasAllowTablePerms(); len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ";"))
	}
	// Case3
	if errMsg := p.checkOnlyHasDenyTablePerms(); len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ";"))
	}
	return nil
}
