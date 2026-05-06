package mfa

import (
	"testing"
	"time"
)

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
