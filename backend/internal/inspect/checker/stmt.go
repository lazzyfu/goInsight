package checker

import (
	"regexp"

	"github.com/lazzyfu/goinsight/pkg/kv"

	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/rules"

	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

type Stmt struct {
	*SyntaxInspectService
}

// CreateTableStmt 检查 CreateTable 语句
func (s *Stmt) CreateTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "CreateTable"}

	for _, rule := range rules.CreateTableRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}

		if rule.RuleHint.IsBreak {
			break
		}
	}

	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}

	return data
}

// CreateViewStmt 检查 CreateView 语句
func (s *Stmt) CreateViewStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "CreateView"}

	for _, rule := range rules.CreateViewRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}

		if rule.RuleHint.IsBreak {
			break
		}
	}

	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}

	return data
}

// CreateDatabaseStmt 检查 CreateDatabase 语句
func (s *Stmt) CreateDatabaseStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "CreateDatabase"}

	for _, rule := range rules.CreateDatabaseRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}

		if rule.RuleHint.IsBreak {
			break
		}
	}

	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}

	return data
}

// RenameTableStmt 检查 RenameTable 语句
func (s *Stmt) RenameTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "RenameTable"}

	for _, rule := range rules.RenameTableRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}
		if rule.RuleHint.IsBreak {
			break
		}
	}

	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}

	return data
}

// AnalyzeTableStmt 检查 AnalyzeTable 语句
func (s *Stmt) AnalyzeTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "AnalyzeTable"}

	for _, rule := range rules.AnalyzeTableRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}

		if rule.RuleHint.IsBreak {
			break
		}
	}

	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}

	return data
}

// DropTableStmt 检查 DropTable 语句
func (s *Stmt) DropTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DropTable"}

	for _, rule := range rules.DropTableRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}
		if rule.RuleHint.IsBreak {
			break
		}
	}

	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}

	return data
}

// AlterTableStmt 检查 AlterTable 语句
func (s *Stmt) AlterTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) (ReturnData, string) {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "AlterTable"}
	var mergeAlter string
	// 禁止使用ALTER TABLE...ADD CONSTRAINT...语法
	tmpCompile := regexp.MustCompile(`(?is:.*alter.*table.*add.*constraint.*)`)
	match := tmpCompile.MatchString(stmt.Text())
	if match {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelWarn, Message: "禁止使用ALTER TABLE...ADD CONSTRAINT...语法"})
		return data, mergeAlter
	}

	for _, rule := range rules.AlterTableRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)
		if len(rule.RuleHint.MergeAlter) > 0 && len(mergeAlter) == 0 {
			mergeAlter = rule.RuleHint.MergeAlter
		}

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}

		if rule.RuleHint.IsBreak {
			// 如果IsBreak为true，跳过接下来的检查步骤
			break
		}
	}
	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}
	return data, mergeAlter
}

// DMLStmt 检查 DML 语句
func (s *Stmt) DMLStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	// delete/update/insert语句
	/*
		DML语句真的需要对同一个指纹的SQL跳过校验？
		1. DML规则并不多，对实际校验性能影响不大
		2. 每条DML都需要进行Explain，由于考虑传值不一样，因此指纹一样并不能代表Explain的影响行数一样
		3. 实际测试1000条update校验仅需800ms,2000条update校验仅需1500ms
		finger := kv.Get(fingerId)
		var IsSkipAudit bool
		if finger != nil {
			IsSkipAudit = true
		}
	*/
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DML"}

	for _, rule := range rules.DMLRules() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		// 当为DML语句时，赋值AffectedRows
		data.AffectedRows = rule.RuleHint.AffectedRows

		if len(ruleHint.Summary) > 0 {
			data.Summary = append(data.Summary, ruleHint.Summary...)
		}

		if rule.RuleHint.IsBreak {
			break
		}
	}

	// 若无任何规则命中/提示，则注入一条默认通过信息
	if len(data.Summary) == 0 {
		data.Summary = append(data.Summary, controllers.SummaryItem{Level: LevelInfo, Message: "语句检查通过"})
	}

	return data
}
