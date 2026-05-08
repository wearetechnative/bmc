package awsconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/ini.v1"
)

// Profile represents a parsed AWS config profile entry.
type Profile struct {
	Name          string
	RoleARN       string
	SourceProfile string
	Group         string
	AccountID     string // extracted from RoleARN
	RoleName      string // extracted from RoleARN
}

// Credentials holds temporary or permanent AWS credentials from ~/.aws/credentials.
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Expiration      string // format: "YYYY-MM-DD HH:MM:SS"
	MFADevice       string // aws_mfa_device field (from long-term section)
}

// LoadProfiles parses ~/.aws/config and returns all profiles.
func LoadProfiles() ([]Profile, error) {
	path := filepath.Join(os.Getenv("HOME"), ".aws", "config")
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("~/.aws/config not found — run 'aws configure' to set up profiles")
		}
		return nil, fmt.Errorf("failed to read ~/.aws/config: %w", err)
	}

	var profiles []Profile
	for _, section := range cfg.Sections() {
		name := section.Name()
		// Config sections are named "profile <name>" except for [default]
		profileName := name
		if name == "DEFAULT" {
			continue
		}
		if len(name) > 8 && name[:8] == "profile " {
			profileName = name[8:]
		}

		roleARN := section.Key("role_arn").String()
		sourcePrf := section.Key("source_profile").String()
		group := section.Key("group").String()

		accountID, roleName := parseARN(roleARN)

		profiles = append(profiles, Profile{
			Name:          profileName,
			RoleARN:       roleARN,
			SourceProfile: sourcePrf,
			Group:         group,
			AccountID:     accountID,
			RoleName:      roleName,
		})
	}
	return profiles, nil
}

// Groups returns unique, sorted group names from the profiles list.
func Groups(profiles []Profile) []string {
	seen := map[string]bool{}
	var groups []string
	for _, p := range profiles {
		if p.Group != "" && !seen[p.Group] {
			seen[p.Group] = true
			groups = append(groups, p.Group)
		}
	}
	sort.Strings(groups)
	return groups
}

// ByGroup returns profiles belonging to the given group.
func ByGroup(profiles []Profile, group string) []Profile {
	var result []Profile
	for _, p := range profiles {
		if p.Group == group {
			result = append(result, p)
		}
	}
	return result
}

// FindProfile finds a profile by name.
func FindProfile(profiles []Profile, name string) (Profile, bool) {
	for _, p := range profiles {
		if p.Name == name {
			return p, true
		}
	}
	return Profile{}, false
}

// ResolveSourceProfile determines the source profile for MFA.
// For role profiles: uses source_profile from config.
// For credentials-only profiles: uses the profile name itself.
func ResolveSourceProfile(profile Profile) (string, error) {
	if profile.SourceProfile != "" {
		return profile.SourceProfile, nil
	}

	// Check if it exists directly in credentials
	creds, err := loadCredentialsFile()
	if err != nil {
		return "", err
	}
	if creds.HasSection(profile.Name) {
		return profile.Name, nil
	}

	return "", fmt.Errorf("cannot resolve source profile for %q", profile.Name)
}

// LoadCredentials reads credentials for a given profile name from ~/.aws/credentials.
func LoadCredentials(profileName string) (Credentials, error) {
	creds, err := loadCredentialsFile()
	if err != nil {
		return Credentials{}, err
	}

	if !creds.HasSection(profileName) {
		return Credentials{}, nil
	}

	sec := creds.Section(profileName)
	return Credentials{
		AccessKeyID:     sec.Key("aws_access_key_id").String(),
		SecretAccessKey: sec.Key("aws_secret_access_key").String(),
		SessionToken:    sec.Key("aws_session_token").String(),
		Expiration:      sec.Key("expiration").String(),
	}, nil
}

// LoadLongTermCredentials reads the [profileName-long-term] section for MFA device.
func LoadLongTermCredentials(profileName string) (Credentials, error) {
	creds, err := loadCredentialsFile()
	if err != nil {
		return Credentials{}, err
	}

	ltSection := profileName + "-long-term"
	if !creds.HasSection(ltSection) {
		return Credentials{}, nil
	}

	sec := creds.Section(ltSection)
	return Credentials{
		AccessKeyID:     sec.Key("aws_access_key_id").String(),
		SecretAccessKey: sec.Key("aws_secret_access_key").String(),
		MFADevice:       sec.Key("aws_mfa_device").String(),
	}, nil
}

// WriteSessionCredentials writes temporary credentials to [profileName] in ~/.aws/credentials.
// Format is compatible with aws-mfa (broamski).
func WriteSessionCredentials(profileName, accessKeyID, secretKey, sessionToken, expiration string) error {
	credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

	f, err := lockAndLoad(credPath)
	if err != nil {
		return err
	}
	defer f.unlock()

	sec, err := f.cfg.GetSection(profileName)
	if err != nil {
		sec, err = f.cfg.NewSection(profileName)
		if err != nil {
			return fmt.Errorf("failed to create credentials section: %w", err)
		}
	}

	mustSet(sec, "aws_access_key_id", accessKeyID)
	mustSet(sec, "aws_secret_access_key", secretKey)
	mustSet(sec, "aws_session_token", sessionToken)
	mustSet(sec, "aws_security_token", sessionToken) // legacy alias used by aws-mfa and older SDKs
	mustSet(sec, "expiration", expiration)

	return f.save(credPath)
}

// loadCredentialsFile parses ~/.aws/credentials (read-only).
func loadCredentialsFile() (*ini.File, error) {
	path := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, path)
	if err != nil {
		if os.IsNotExist(err) {
			return ini.Empty(), nil
		}
		return nil, fmt.Errorf("failed to read ~/.aws/credentials: %w", err)
	}
	return cfg, nil
}

// mustSet sets a key in an ini section, creating it if absent.
func mustSet(sec *ini.Section, key, value string) {
	if sec.HasKey(key) {
		sec.Key(key).SetValue(value)
	} else {
		_, _ = sec.NewKey(key, value)
	}
}

// parseARN extracts account ID and role name from an IAM role ARN.
// arn:aws:iam::123456789012:role/MyRole → ("123456789012", "MyRole")
func parseARN(arn string) (accountID, roleName string) {
	// arn:aws:iam::<account>:role/<name>
	var prefix, resource string
	if _, err := fmt.Sscanf(arn, "arn:aws:iam::%s", &prefix); err != nil {
		return "", ""
	}
	// prefix is "<account>:role/<name>"
	n := len(prefix)
	colonIdx := -1
	for i := 0; i < n; i++ {
		if prefix[i] == ':' {
			colonIdx = i
			break
		}
	}
	if colonIdx < 0 {
		return prefix, ""
	}
	accountID = prefix[:colonIdx]
	resource = prefix[colonIdx+1:]
	if len(resource) > 5 && resource[:5] == "role/" {
		roleName = resource[5:]
	}
	return accountID, roleName
}
