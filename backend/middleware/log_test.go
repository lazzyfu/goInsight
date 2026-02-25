package middleware

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestSanitizeRequestHeaders(t *testing.T) {
	headers := http.Header{
		"Authorization": []string{"JWT abc"},
		"Cookie":        []string{"jwt=abc"},
		"X-Trace-Id":    []string{"trace-1"},
	}

	got := sanitizeRequestHeaders(headers)
	if got["Authorization"][0] != "***" {
		t.Fatalf("authorization should be masked, got %q", got["Authorization"][0])
	}
	if got["Cookie"][0] != "***" {
		t.Fatalf("cookie should be masked, got %q", got["Cookie"][0])
	}
	if got["X-Trace-Id"][0] != "trace-1" {
		t.Fatalf("non-sensitive header should keep value, got %q", got["X-Trace-Id"][0])
	}
}

func TestSanitizeRequestURI(t *testing.T) {
	u, err := url.Parse("/api/v1/user/otp-auth-url?username=alice&password=secret&token=t1")
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}

	got := sanitizeRequestURI(u)
	if got == "" {
		t.Fatal("sanitized URI should not be empty")
	}
	if got == "/api/v1/user/otp-auth-url?username=alice&password=secret&token=t1" {
		t.Fatal("sensitive query fields should be masked")
	}
	if want := "username=alice"; !strings.Contains(got, want) {
		t.Fatalf("expected non-sensitive query to remain, got %q", got)
	}
	if want := "password=%2A%2A%2A"; !strings.Contains(got, want) {
		t.Fatalf("expected password to be masked, got %q", got)
	}
	if want := "token=%2A%2A%2A"; !strings.Contains(got, want) {
		t.Fatalf("expected token to be masked, got %q", got)
	}
}
