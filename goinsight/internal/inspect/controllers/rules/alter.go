/*
@Time    :   2022/06/28 10:21:41
@Author  :   xff
@Desc    :   alter规则逻辑，Level初始化为INFO
*/

package rules

import (
	"goInsight/internal/inspect/controllers/logics"
	"goInsight/internal/inspect/controllers/traverses"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

func AlterTableRules() []Rule {
	return []Rule{
		{
			Hint:      "AlterTable#检查表是否存在",
			CheckFunc: (*Rule).RuleAlterTableIsExist,
		},
		{
			Hint:      "AlterTable#检查TiDBMergeAlter",
			CheckFunc: (*Rule).RuleAlterTiDBMerge,
		},
		{
			Hint:      "AlterTable#DROP列和索引检查",
			CheckFunc: (*Rule).RuleAlterTableDropColsOrIndexes,
		},
		{
			Hint:      "AlterTable#DropTiDBColWithCoveredIndex检查",
			CheckFunc: (*Rule).RuleAlterTableDropTiDBColWithCoveredIndex,
		},
		{
			Hint:      "AlterTable#表Options检查",
			CheckFunc: (*Rule).RuleAlterTableOptions,
		},
		{
			Hint:      "AlterTable#列字符集检查",
			CheckFunc: (*Rule).RuleAlterTableColCharset,
		},
		{
			Hint:      "AlterTable#Add列After检查",
			CheckFunc: (*Rule).RuleAlterTableAddColAfter,
		},
		{
			Hint:      "AlterTable#Add列Options检查",
			CheckFunc: (*Rule).RuleAlterTableAddColOptions,
		},
		{
			Hint:      "AlterTable#Add主键检查",
			CheckFunc: (*Rule).RuleAlterTableAddPrimaryKey,
		},
		{
			Hint:      "AlterTable#Add重复列检查",
			CheckFunc: (*Rule).RuleAlterTableAddColRepeatDefine,
		},
		{
			Hint:      "AlterTable#Add索引前缀检查",
			CheckFunc: (*Rule).RuleAlterTableAddIndexPrefix,
		},
		{
			Hint:      "AlterTable#Add索引数量检查",
			CheckFunc: (*Rule).RuleAlterTableAddIndexCount,
		},
		{
			Hint:      "AlterTable#AddConstraint检查",
			CheckFunc: (*Rule).RuleAlterTableAddConstraint,
		},
		{
			Hint:      "AlterTable#Add重复索引检查",
			CheckFunc: (*Rule).RuleAlterTableAddIndexRepeatDefine,
		},
		{
			Hint:      "AlterTable#Add冗余索引检查",
			CheckFunc: (*Rule).RuleAlterTableRedundantIndexes,
		},
		{
			Hint:      "AlterTable#BLOB/TEXT类型不能设置为索引",
			CheckFunc: (*Rule).RuleAlterTableDisabledIndexes,
		},
		{
			Hint:      "AlterTable#Modify列Options检查",
			CheckFunc: (*Rule).RuleAlterTableModifyColOptions,
		},
		{
			Hint:      "AlterTable#Change列Options检查",
			CheckFunc: (*Rule).RuleAlterTableChangeColOptions,
		},
		{
			Hint:      "AlterTable#RenameIndex检查",
			CheckFunc: (*Rule).RuleAlterTableRenameIndex,
		},
		{
			Hint:      "AlterTable#RenameTblName检查",
			CheckFunc: (*Rule).RuleAlterTableRenameTblName,
		},
		{
			Hint:      "AlterTable#索引InnodbLargePrefix",
			CheckFunc: (*Rule).RuleAlterTableInnodbLargePrefix,
		},
		{
			Hint:      "AlterTable#检查表定义的行是否超过65535",
			CheckFunc: (*Rule).RuleAlterTableInnoDBRowSize,
		},
	}
}

// RuleAlterTableIsExist
func (r *Rule) RuleAlterTableIsExist(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableIsExist{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableIsExist(v, r.RuleHint)
}

// RuleAlterTiDBMerge
func (r *Rule) RuleAlterTiDBMerge(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTiDBMerge{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableTiDBMerge(v, r.RuleHint)
}

// RuleAlterTableDropCols
func (r *Rule) RuleAlterTableDropColsOrIndexes(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableDropColsOrIndexes{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableDropColsOrIndexes(v, r.RuleHint)
}

// RuleAlterTableDropTiDBColWithCoveredIndex
func (r *Rule) RuleAlterTableDropTiDBColWithCoveredIndex(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableDropTiDBColWithCoveredIndex{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableDropTiDBColWithCoveredIndex(v, r.RuleHint)
}

// RuleAlterTableOptions
func (r *Rule) RuleAlterTableOptions(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableOptions{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableOptions(v, r.RuleHint)
}

// RuleAlterTableColCharset
func (r *Rule) RuleAlterTableColCharset(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableColCharset{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableColCharset(v, r.RuleHint)
}

// RuleAlterTableAddColAfter
func (r *Rule) RuleAlterTableAddColAfter(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddColAfter{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddColAfter(v, r.RuleHint)
}

// RuleAlterTableAddColOptions
func (r *Rule) RuleAlterTableAddColOptions(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddColOptions{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddColOptions(v, r.RuleHint)
}

// RuleAlterTableAddColWithPrimaryKey
func (r *Rule) RuleAlterTableAddPrimaryKey(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddPrimaryKey{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddPrimaryKey(v, r.RuleHint)
}

// RuleAlterTableAddColRepeatDefine
func (r *Rule) RuleAlterTableAddColRepeatDefine(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddColRepeatDefine{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddColRepeatDefine(v, r.RuleHint)
}

// RuleAlterTableAddIndexPrefix
func (r *Rule) RuleAlterTableAddIndexPrefix(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddIndexPrefix{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddIndexPrefix(v, r.RuleHint)
}

// RuleAlterTableAddIndexCount
func (r *Rule) RuleAlterTableAddIndexCount(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddIndexCount{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddIndexCount(v, r.RuleHint)
}

// RuleAlterTableAddConstraint
func (r *Rule) RuleAlterTableAddConstraint(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddConstraint{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddConstraint(v, r.RuleHint)
}

// RuleAlterTableAddIndexRepeatDefine
func (r *Rule) RuleAlterTableAddIndexRepeatDefine(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableAddIndexRepeatDefine{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableAddIndexRepeatDefine(v, r.RuleHint)
}

// RuleAlterTableRedundantIndexes
func (r *Rule) RuleAlterTableRedundantIndexes(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableRedundantIndexes{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableRedundantIndexes(v, r.RuleHint)
}

// RuleAlterTableDisabledIndexes
func (r *Rule) RuleAlterTableDisabledIndexes(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableDisabledIndexes{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableDisabledIndexes(v, r.RuleHint)
}

// RuleAlterTableModifyColOptions
func (r *Rule) RuleAlterTableModifyColOptions(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableModifyColOptions{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableModifyColOptions(v, r.RuleHint)
}

// RuleAlterTableChangeColOptions
func (r *Rule) RuleAlterTableChangeColOptions(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableChangeColOptions{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableChangeColOptions(v, r.RuleHint)
}

// RuleAlterTableRenameIndex
func (r *Rule) RuleAlterTableRenameIndex(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableRenameIndex{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableRenameIndex(v, r.RuleHint)
}

// RuleAlterTableRenameTblName
func (r *Rule) RuleAlterTableRenameTblName(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableRenameTblName{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableRenameTblName(v, r.RuleHint)
}

// RuleAlterTableInnodbLargePrefix
func (r *Rule) RuleAlterTableInnodbLargePrefix(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableInnodbLargePrefix{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableInnodbLargePrefix(v, r.RuleHint)
}

// RuleAlterTableInnoDBRowSize
func (r *Rule) RuleAlterTableInnoDBRowSize(tistmt *ast.StmtNode) {
	v := &traverses.TraverseAlterTableInnoDBRowSize{}
	(*tistmt).Accept(v)
	logics.LogicAlterTableInnoDBRowSize(v, r.RuleHint)
}
