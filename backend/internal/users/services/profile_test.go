package services

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

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

func TestValidateOTPBindingCode(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatalf("generate otp code failed: %v", err)
	}

	t.Run("valid otp code should pass", func(t *testing.T) {
		if err := validateOTPBindingCode(secret, code); err != nil {
			t.Fatalf("expected valid otp code, got error: %v", err)
		}
	})

	t.Run("invalid format should fail", func(t *testing.T) {
		if err := validateOTPBindingCode(secret, "12a456"); err == nil {
			t.Fatal("expected invalid format to fail")
		}
	})

	t.Run("wrong code should fail", func(t *testing.T) {
		if err := validateOTPBindingCode(secret, "000000"); err == nil {
			t.Fatal("expected wrong otp code to fail")
		}
	})
}
