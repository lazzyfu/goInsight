package services

import (
	"fmt"
	"goInsight/internal/pkg/kv"
	"goInsight/internal/pkg/utils"
	"regexp"

	"goInsight/internal/app/inspect/controllers/process"
	"goInsight/internal/app/inspect/controllers/rules"

	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

type Stmt struct {
	*SyntaxInspectService
}

func (s *Stmt) CreateTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	// 建表语句
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DDL", Level: "INFO"}
	for _, rule := range rules.CreateTableRules() {
		rule.DB = s.DB
		rule.KV = kv
		rule.AuditConfig = &s.AuditConfig
		rule.CheckFunc(&rule, &stmt)
		if len(rule.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.Summary...)
		}
		if rule.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data
}

func (s *Stmt) CreateViewStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	// 建视图语句
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DDL", Level: "INFO"}
	for _, rule := range rules.CreateViewRules() {
		rule.DB = s.DB
		rule.KV = kv
		rule.AuditConfig = &s.AuditConfig
		rule.CheckFunc(&rule, &stmt)
		if len(rule.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.Summary...)
		}
		if rule.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data
}

func (s *Stmt) AlterTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) (ReturnData, string) {
	// alter语句
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DDL", Level: "INFO"}
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
		rule.DB = s.DB
		rule.KV = kv
		rule.AuditConfig = &s.AuditConfig
		rule.CheckFunc(&rule, &stmt)
		if len(rule.MergeAlter) > 0 && len(mergeAlter) == 0 {
			mergeAlter = rule.MergeAlter
		}
		if len(rule.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.Summary...)
		}
		if rule.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data, mergeAlter
}

func (s *Stmt) RenameTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	// rename table语句
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DDL", Level: "INFO"}
	for _, rule := range rules.RenameTableRules() {
		rule.DB = s.DB
		rule.KV = kv
		rule.AuditConfig = &s.AuditConfig
		rule.CheckFunc(&rule, &stmt)
		if len(rule.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.Summary...)
		}
		if rule.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data
}

func (s *Stmt) AnalyzeTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	// analyze table语句
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DDL", Level: "INFO"}
	for _, rule := range rules.AnalyzeTableRules() {
		rule.DB = s.DB
		rule.KV = kv
		rule.AuditConfig = &s.AuditConfig
		rule.CheckFunc(&rule, &stmt)
		if len(rule.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.Summary...)
		}
		if rule.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data
}

func (s *Stmt) DropTableStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	// drop/truncate语句
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DDL", Level: "INFO"}
	for _, rule := range rules.DropTableRules() {
		rule.DB = s.DB
		rule.KV = kv
		rule.AuditConfig = &s.AuditConfig
		rule.CheckFunc(&rule, &stmt)
		if len(rule.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.Summary...)
		}
		if rule.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data
}

func (s *Stmt) DMLStmt(stmt ast.StmtNode, kv *kv.KVCache, fingerId string) ReturnData {
	// delete/update/insert语句
	var data ReturnData = ReturnData{FingerId: fingerId, Query: stmt.Text(), Type: "DML", Level: "INFO"}
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
	for _, rule := range rules.DMLRules() {
		rule.DB = s.DB
		rule.KV = kv
		rule.AuditConfig = &s.AuditConfig
		rule.Query = stmt.Text()
		rule.CheckFunc(&rule, &stmt)
		data.AffectedRows = rule.AffectedRows
		if len(rule.Summary) > 0 {
			// 检查不通过
			data.Level = "WARN"
			data.Summary = append(data.Summary, rule.Summary...)
		}
		if rule.IsSkipNextStep {
			// 如果IsSkipNextStep为true，跳过接下来的检查步骤
			break
		}
	}
	return data
}

func (s *Stmt) MergeAlter(kv *kv.KVCache, mergeAlters []string) ReturnData {
	// 检查mysql merge操作
	var data ReturnData = ReturnData{Level: "INFO"}
	dbVersionIns := process.DbVersion{Version: kv.Get("dbVersion").(string)}
	if s.AuditConfig.ENABLE_MYSQL_MERGE_ALTER_TABLE && !dbVersionIns.IsTiDB() {
		if ok, val := utils.IsRepeat(mergeAlters); ok {
			for _, v := range val {
				data.Summary = append(data.Summary, fmt.Sprintf("[MySQL数据库]表`%s`的多条ALTER操作，请合并为一条ALTER语句", v))
			}
		}
	}
	if len(data.Summary) > 0 {
		data.Level = "WARN"
	}
	return data
}
