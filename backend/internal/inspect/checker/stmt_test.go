package checker

import (
	"testing"

	"github.com/lazzyfu/goinsight/internal/inspect/controllers"
	"github.com/pingcap/tidb/pkg/parser/ast"
)

func TestAlterTableStmt_AddDefaultSummaryWhenEmpty(t *testing.T) {
	s := &Stmt{}
	data, _ := s.AlterTableStmt(&ast.AlterTableStmt{}, nil, "finger")

	if len(data.Summary) != 1 {
		t.Fatalf("expected 1 summary item, got %d", len(data.Summary))
	}
	if data.Summary[0] != (controllers.SummaryItem{Level: LevelInfo, Message: "ALTER语句检查通过"}) {
		t.Fatalf("unexpected summary item: %#v", data.Summary[0])
	}
}

func TestAlterTableStmt_DoNotOverrideExistingSummary(t *testing.T) {
	s := &Stmt{}
	data, _ := s.AlterTableStmt(&mockAlterStmt{text: "ALTER TABLE t ADD CONSTRAINT c CHECK (1=1)"}, nil, "finger")

	if len(data.Summary) == 0 {
		t.Fatalf("expected non-empty summary")
	}
	if got := data.Summary[len(data.Summary)-1]; got.Level == LevelInfo && got.Message == "ALTER语句检查通过" {
		t.Fatalf("should not inject default pass message when summary already exists")
	}
}

type mockAlterStmt struct {
	ast.AlterTableStmt
	text string
}

func (m *mockAlterStmt) Text() string { return m.text }
