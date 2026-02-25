package services

import "testing"

func TestValidatePasswordConfirmation(t *testing.T) {
	t.Run("mismatch should fail", func(t *testing.T) {
		err := validatePasswordConfirmation("new-password-1", "new-password-2")
		if err == nil {
			t.Fatal("expected mismatch to fail")
		}
	})

	t.Run("same should pass", func(t *testing.T) {
		err := validatePasswordConfirmation("new-password", "new-password")
		if err != nil {
			t.Fatalf("expected matching passwords to pass, got %v", err)
		}
	})
}
