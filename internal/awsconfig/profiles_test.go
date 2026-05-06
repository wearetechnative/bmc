package awsconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseARN(t *testing.T) {
	tests := []struct {
		arn       string
		accountID string
		roleName  string
	}{
		{"arn:aws:iam::123456789012:role/MyRole", "123456789012", "MyRole"},
		{"arn:aws:iam::000000000000:role/Admin", "000000000000", "Admin"},
		{"", "", ""},
		{"not-an-arn", "", ""},
	}
	for _, tt := range tests {
		a, r := parseARN(tt.arn)
		if a != tt.accountID || r != tt.roleName {
			t.Errorf("parseARN(%q) = (%q, %q), want (%q, %q)", tt.arn, a, r, tt.accountID, tt.roleName)
		}
	}
}

func TestGroups(t *testing.T) {
	profiles := []Profile{
		{Name: "a", Group: "dev"},
		{Name: "b", Group: "prod"},
		{Name: "c", Group: "dev"},
		{Name: "d", Group: ""},
	}
	got := Groups(profiles)
	if len(got) != 2 || got[0] != "dev" || got[1] != "prod" {
		t.Errorf("Groups() = %v, want [dev prod]", got)
	}
}

func TestByGroup(t *testing.T) {
	profiles := []Profile{
		{Name: "a", Group: "dev"},
		{Name: "b", Group: "prod"},
		{Name: "c", Group: "dev"},
	}
	devProfiles := ByGroup(profiles, "dev")
	if len(devProfiles) != 2 {
		t.Errorf("ByGroup(dev) returned %d profiles, want 2", len(devProfiles))
	}
}

func TestWriteReadSessionCredentials(t *testing.T) {
	// Use a temp dir to avoid touching real ~/.aws/credentials
	tmpDir := t.TempDir()
	credPath := filepath.Join(tmpDir, "credentials")
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	// Point HOME at a tmp dir structure
	awsDir := filepath.Join(tmpDir, ".aws")
	if err := os.MkdirAll(awsDir, 0700); err != nil {
		t.Fatal(err)
	}
	os.Setenv("HOME", tmpDir)
	credPath = filepath.Join(tmpDir, ".aws", "credentials")

	err := WriteSessionCredentials("testprofile",
		"AKIATEST", "secretkey", "sessiontoken", "2026-05-06 12:00:00")
	if err != nil {
		t.Fatalf("WriteSessionCredentials error: %v", err)
	}

	creds, err := LoadCredentials("testprofile")
	if err != nil {
		t.Fatalf("LoadCredentials error: %v", err)
	}

	if creds.AccessKeyID != "AKIATEST" {
		t.Errorf("AccessKeyID = %q, want AKIATEST", creds.AccessKeyID)
	}
	if creds.Expiration != "2026-05-06 12:00:00" {
		t.Errorf("Expiration = %q, want 2026-05-06 12:00:00", creds.Expiration)
	}
	_ = credPath
}
