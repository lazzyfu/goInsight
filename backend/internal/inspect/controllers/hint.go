package controllers

import (
	"github.com/lazzyfu/goinsight/pkg/kv"

	"github.com/lazzyfu/goinsight/internal/inspect/config"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
)

type RuleHint struct {
	Summary        []string `json:"summary"`       // 摘要
	AffectedRows   int      `json:"affected_rows"` // 默认为0
	IsSkipNextStep bool     // 是否跳过接下来的检查步骤
	DB             *dao.DB
	KV             *kv.KVCache
	Query          string // 原始SQL
	MergeAlter     string
	InspectParams  *config.InspectParams
}
