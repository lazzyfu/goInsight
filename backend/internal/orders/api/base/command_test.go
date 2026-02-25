package base

import (
	"strings"
	"testing"
)

func TestRenderCommandForLog(t *testing.T) {
	cmd := RenderCommandForLog("/usr/local/bin/gh-ost", []string{
		"--user=test",
		"--password=secret",
		"--host=127.0.0.1",
	})
	if cmd == "" {
		t.Fatal("log command should not be empty")
	}
	if strings.Contains(cmd, "secret") {
		t.Fatalf("password should be masked, got %q", cmd)
	}
	if !strings.Contains(cmd, "--password=...") {
		t.Fatalf("masked password should appear, got %q", cmd)
	}
}
