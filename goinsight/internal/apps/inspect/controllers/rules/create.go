/*
@Time    :   2022/07/06 10:12:42
@Author  :   zongfei.fu
@Desc    :   None
*/

package rules

import (
	"goInsight/internal/apps/inspect/controllers/logics"
	"goInsight/internal/apps/inspect/controllers/traverses"

	"github.com/pingcap/tidb/parser/ast"
)

func CreateTableRules() []Rule {
	return []Rule{
		{
			Hint:      "CreateTable#检查表是否存在",
			CheckFunc: (*Rule).RuleCreateTableIsExist,
		},
		{
			Hint:      "CreateTable#检查CreateTableAs语法",
			CheckFunc: (*Rule).RuleCreateTableAs,
		},
		{
			Hint:      "CreateTable#检查CreateTableLike语法",
			CheckFunc: (*Rule).RuleCreateTableLike,
		},
		{
			Hint:      "CreateTable#表Options检查",
			CheckFunc: (*Rule).RuleCreateTableOptions,
		},
		{
			Hint:      "CreateTable#主键检查",
			CheckFunc: (*Rule).RuleCreateTablePrimaryKey,
		},
		{
			Hint:      "CreateTable#约束检查",
			CheckFunc: (*Rule).RuleCreateTableConstraint,
		},
		{
			Hint:      "CreateTable#审计字段检查",
			CheckFunc: (*Rule).RuleCreateTableAuditCols,
		},
		{
			Hint:      "CreateTable#列Options检查",
			CheckFunc: (*Rule).RuleCreateTableColsOptions,
		},
		{
			Hint:      "CreateTable#列重复定义检查",
			CheckFunc: (*Rule).RuleCreateTableColsRepeatDefine,
		},
		{

			Hint:      "CreateTable#列字符集检查",
			CheckFunc: (*Rule).RuleCreateTableColsCharset,
		},
		{
			Hint: "CreateTable#索引前缀检查",

			CheckFunc: (*Rule).RuleCreateTableIndexesPrefix,
		},
		{
			Hint:      "CreateTable#索引数量检查",
			CheckFunc: (*Rule).RuleCreateTableIndexesCount,
		},
		{
			Hint:      "CreateTable#索引重复定义检查",
			CheckFunc: (*Rule).RuleCreateTableIndexesRepeatDefine,
		},
		{
			Hint:      "CreateTable#冗余索引检查",
			CheckFunc: (*Rule).RuleCreateTableRedundantIndexes,
		},
		{
			Hint:      "CreateTable#BLOB/TEXT类型不能设置为索引",
			CheckFunc: (*Rule).RuleCreateTableDisabledIndexes,
		},
		{
			Hint:      "CreateTable#索引InnodbLargePrefix",
			CheckFunc: (*Rule).RuleCreateTableInnodbLargePrefix,
		},
		{
			Hint:      "CreateTable#检查InnoDB表定义的行大小",
			CheckFunc: (*Rule).RuleCreateTableInnoDBRowSize,
		},
	}
}

// RuleCreateTableIsExist
func (r *Rule) RuleCreateTableIsExist(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableIsExist{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableIsExist(v, r.RuleHint)
}

// RuleCreateTableAs
func (r *Rule) RuleCreateTableAs(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableAs{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableAs(v, r.RuleHint)
}

// RuleCreateTableLike
func (r *Rule) RuleCreateTableLike(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableLike{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableLike(v, r.RuleHint)
}

// RuleCreateTableOptions
func (r *Rule) RuleCreateTableOptions(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableOptions{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableOptions(v, r.RuleHint)
}

// RuleCreateTablePrimaryKey
func (r *Rule) RuleCreateTablePrimaryKey(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTablePrimaryKey{}
	(*tistmt).Accept(v)
	logics.LogicCreateTablePrimaryKey(v, r.RuleHint)
}

// RuleCreateTableConstraint
func (r *Rule) RuleCreateTableConstraint(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableConstraint{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableConstraint(v, r.RuleHint)
}

// RuleCreateTableAuditCols
func (r *Rule) RuleCreateTableAuditCols(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableAuditCols{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableAuditCols(v, r.RuleHint)
}

// RuleCreateTableColsOptions
func (r *Rule) RuleCreateTableColsOptions(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableColsOptions{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableColsOptions(v, r.RuleHint)
}

// RuleCreateTableColsRepeatDefine
func (r *Rule) RuleCreateTableColsRepeatDefine(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableColsRepeatDefine{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableColsRepeatDefine(v, r.RuleHint)
}

// RuleCreateTableColsCharset
func (r *Rule) RuleCreateTableColsCharset(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableColsCharset{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableColsCharset(v, r.RuleHint)
}

// RuleCreateTableIndexesPrefix
func (r *Rule) RuleCreateTableIndexesPrefix(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableIndexesPrefix{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableIndexesPrefix(v, r.RuleHint)
}

// RuleCreateTableIndexesCount
func (r *Rule) RuleCreateTableIndexesCount(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableIndexesCount{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableIndexesCount(v, r.RuleHint)
}

// RuleCreateTableIndexesRepeatDefine
func (r *Rule) RuleCreateTableIndexesRepeatDefine(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableIndexesRepeatDefine{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableIndexesRepeatDefine(v, r.RuleHint)
}

// RuleCreateTableRedundantIndexes
func (r *Rule) RuleCreateTableRedundantIndexes(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableRedundantIndexes{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableRedundantIndexes(v, r.RuleHint)
}

// RuleCreateTableDisabledIndexes
func (r *Rule) RuleCreateTableDisabledIndexes(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableDisabledIndexes{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableDisabledIndexes(v, r.RuleHint)
}

// RuleCreateTableInnodbLargePrefix
func (r *Rule) RuleCreateTableInnodbLargePrefix(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableInnodbLargePrefix{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableInnodbLargePrefix(v, r.RuleHint)
}

// RuleCreateTableInnoDBRowSize
func (r *Rule) RuleCreateTableInnoDBRowSize(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateTableInnoDBRowSize{}
	(*tistmt).Accept(v)
	logics.LogicCreateTableInnoDBRowSize(v, r.RuleHint)
}
