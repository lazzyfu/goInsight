package dao

import "testing"

func TestValidateIdentifier(t *testing.T) {
	if err := validateIdentifier("test_db", "database"); err != nil {
		t.Fatalf("expected valid identifier, got %v", err)
	}
	if err := validateIdentifier("test-db", "database"); err == nil {
		t.Fatal("expected invalid identifier to fail")
	}
}
