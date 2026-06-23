package mfa

import (
	"os"
	"testing"
	"time"

	"github.com/wearetechnative/bmc/internal/config"
)

func TestAcquireTOTP_ProfileScript(t *testing.T) {
	cfg := config.Config{
		MFA: config.MFAConfig{
			TOTPScript: `echo "global"`,
			ProfileScripts: map[string]string{
				"myprofile": `echo "profile-specific"`,
			},
		},
	}
	code, err := acquireTOTP(cfg, "myprofile", os.Stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "profile-specific" {
		t.Errorf("expected profile-specific script output, got %q", code)
	}
}

func TestAcquireTOTP_GlobalFallback(t *testing.T) {
	cfg := config.Config{
		MFA: config.MFAConfig{
			TOTPScript: `echo "global"`,
			ProfileScripts: map[string]string{
				"other": `echo "other"`,
			},
		},
	}
	code, err := acquireTOTP(cfg, "myprofile", os.Stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "global" {
		t.Errorf("expected global script output, got %q", code)
	}
}

func TestAcquireTOTP_NilMap(t *testing.T) {
	cfg := config.Config{
		MFA: config.MFAConfig{
			TOTPScript:     `echo "global"`,
			ProfileScripts: nil,
		},
	}
	code, err := acquireTOTP(cfg, "myprofile", os.Stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "global" {
		t.Errorf("expected global fallback on nil map, got %q", code)
	}
}

func TestIsValid(t *testing.T) {
	future := time.Now().Add(time.Hour).UTC().Format(expirationLayout)
	past := time.Now().Add(-time.Hour).UTC().Format(expirationLayout)

	if !isValid(future) {
		t.Error("expected future expiration to be valid")
	}
	if isValid(past) {
		t.Error("expected past expiration to be invalid")
	}
	if isValid("") {
		t.Error("expected empty expiration to be invalid")
	}
	if isValid("1970-01-01 01:00:00") {
		t.Error("expected epoch to be invalid")
	}
}
