/*
@Time    :   2022/06/29 15:30:31
@Author  :   xff
*/

package controllers

import (
	"goInsight/internal/apps/inspect/config"
	"goInsight/internal/apps/inspect/controllers/dao"
	"goInsight/internal/pkg/kv"
)

type RuleHint struct {
	Summary        []string `json:"summary"` // 摘要
	AffectedRows   int      `json:"affected_rows"`
	IsSkipNextStep bool     // 是否跳过接下来的检查步骤
	DB             *dao.DB
	KV             *kv.KVCache
	Query          string // 原始SQL
	MergeAlter     string
	InspectParams  *config.InspectParams
}
