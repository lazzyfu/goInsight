package controllers

import (
	"github.com/lazzyfu/goinsight/pkg/kv"

	"github.com/lazzyfu/goinsight/internal/inspect/config"
	"github.com/lazzyfu/goinsight/internal/inspect/controllers/dao"
)

// level: INFO/WARN/ERROR
type SummaryItem struct {
	Level   string `json:"level"` // 级别,INFO/WARN/ERROR
	Message string `json:"message"`
}

type RuleHint struct {
	Summary      []SummaryItem `json:"summary"`        // 摘要
	AffectedRows int           `json:"affected_rows"`  // 默认为0
	IsBreak      bool          `json:"skip_next_step"` // 是否跳过接下来的检查步骤
	MergeAlter   string        `json:"merge_alter,omitempty"`

	// 内部执行上下文：禁止输出到 JSON（避免泄露/序列化失败）
	DB            *dao.DB               `json:"-"`
	KV            *kv.KVCache           `json:"-"`
	Query         string                `json:"-"` // 原始SQL
	InspectParams *config.InspectParams `json:"-"`
}

func (r *RuleHint) AddSummary(level, message string) {
	if r == nil {
		return
	}
	r.Summary = append(r.Summary, SummaryItem{Message: message, Level: level})
}

func (r *RuleHint) Info(message string) { r.AddSummary("INFO", message) }
func (r *RuleHint) Warn(message string) { r.AddSummary("WARN", message) }
