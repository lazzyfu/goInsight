package services

import (
	"testing"

	"gorm.io/datatypes"
)

func TestNormalizeClaimUsers(t *testing.T) {
	t.Run("dedupe and trim", func(t *testing.T) {
		users, err := normalizeClaimUsers([]string{" alice ", "bob", "alice", ""})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(users) != 2 || users[0] != "alice" || users[1] != "bob" {
			t.Fatalf("unexpected users: %#v", users)
		}
	})

	t.Run("empty should fail", func(t *testing.T) {
		if _, err := normalizeClaimUsers([]string{"", "   "}); err == nil {
			t.Fatal("expected error for empty claim users")
		}
	})
}

func TestCanUserClaim(t *testing.T) {
	raw := datatypes.JSON([]byte(`["alice","bob"]`))

	ok, err := canUserClaim(raw, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("alice should be allowed to claim")
	}

	ok, err = canUserClaim(raw, "charlie")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("charlie should not be allowed to claim")
	}
}

func TestCanUserClaimInvalidJSON(t *testing.T) {
	if _, err := canUserClaim(datatypes.JSON([]byte(`{}`)), "alice"); err == nil {
		t.Fatal("expected error for invalid claim users json")
	}
}
