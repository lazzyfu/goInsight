package checker

import (
	"regexp"

	"goInsight/internal/app/inspect/controllers"
	"goInsight/internal/app/inspect/controllers/rules"
	"goInsight/internal/pkg/kv"

	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

type Stmt struct {
	*SyntaxInspectService
}

func (s *Stmt) commonCheck(stmt ast.StmtNode, kv *kv.KVCache, fingerId string, sqlType string, rulesFunc func() []rules.Rule) ReturnData {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: sqlType, Level: "INFO"}

	for _, rule := range rulesFunc() {
		var ruleHint *controllers.RuleHint = &controllers.RuleHint{
			DB:            s.DB,
			KV:            kv,
			Query:         stmt.Text(),
			InspectParams: &s.InspectParams,
		}
		rule.RuleHint = ruleHint
		rule.CheckFunc(&rule, &stmt)

		if len(rule.RuleHint.Summary) > 0 {
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.RuleHint.Summary...)
		}
		if rule.RuleHint.IsSkipNextStep {
			break
		}
	}

	return data
}

func (s *Stmt) CreateTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	return s.commonCheck(stmt, kv, fingerId, "DDL/CreateTable", rules.CreateTableRules)
}

func (s *Stmt) CreateViewStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	return s.commonCheck(stmt, kv, fingerId, "DDL/CreateView", rules.CreateTableRules)
}

func (s *Stmt) RenameTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	return s.commonCheck(stmt, kv, fingerId, "DDL/RenameTable", rules.CreateTableRules)
}

func (s *Stmt) AnalyzeTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	return s.commonCheck(stmt, kv, fingerId, "DDL/AnalyzeTable", rules.CreateTableRules)
}

func (s *Stmt) DropTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	return s.commonCheck(stmt, kv, fingerId, "DDL/DropTable", rules.CreateTableRules)
}

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
	return s.commonCheck(stmt, kv, fingerId, "DML", rules.CreateTableRules)
}

func (s *Stmt) AlterTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) (ReturnData, string) {
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DDL/AlterTable", Level: "INFO"}
	var mergeAlter string
	// 禁止使用ALTER TABLE...ADD CONSTRAINT...语法
	tmpCompile := regexp.MustCompile(`(?is:.*alter.*table.*add.*constraint.*)`)
	match := tmpCompile.MatchString(stmt.Text())
	if match {
		data.Level = "WARN"
		data.Summary = append(data.Summary, "禁止使用ALTER TABLE...ADD CONSTRAINT...语法")
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
		if len(rule.RuleHint.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.RuleHint.Summary...)
		}
		if rule.RuleHint.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data, mergeAlter
}
