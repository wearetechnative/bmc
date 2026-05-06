package mfa

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

const expirationLayout = "2006-01-02 15:04:05"

// EnsureValid checks if the MFA session for sourceProfile is still valid.
// If not, it refreshes it. Skips entirely if mfa.enabled is false.
func EnsureValid(sourceProfile string, cfg config.Config, outfd *os.File) error {
	if !cfg.MFA.Enabled {
		return nil
	}

	if outfd == nil {
		outfd = os.Stderr
	}

	fmt.Fprintf(outfd, "-- Using AWS source-profile: %s\n", sourceProfile)

	creds, err := awsconfig.LoadCredentials(sourceProfile)
	if err != nil {
		return err
	}

	ltCreds, err := awsconfig.LoadLongTermCredentials(sourceProfile)
	if err != nil {
		return err
	}

	if ltCreds.MFADevice == "" {
		fmt.Fprintln(outfd, "!! AWS MFA Device not found. Can't renew session")
		return nil
	}

	if isValid(creds.Expiration) {
		exp, _ := time.Parse(expirationLayout, creds.Expiration)
		fmt.Fprintf(outfd, "Current MFA Session Valid, until: %s\n\n", exp.Format(expirationLayout))
		return nil
	}

	fmt.Fprintf(outfd, "-- Refreshing MFA session for %s...\n", sourceProfile)

	totpCode, err := acquireTOTP(cfg, outfd)
	if err != nil {
		return fmt.Errorf("failed to get TOTP code: %w", err)
	}

	newCreds, err := callGetSessionToken(ltCreds, ltCreds.MFADevice, totpCode)
	if err != nil {
		return fmt.Errorf("!!  Error with AWS MFA code for device. Wrong TOTP? %w", err)
	}

	accessKey, secretKey, sessionToken, expTime := extractCredentials(newCreds)
	expStr := expTime.UTC().Format(expirationLayout)

	if err := awsconfig.WriteSessionCredentials(sourceProfile, accessKey, secretKey, sessionToken, expStr); err != nil {
		return fmt.Errorf("failed to write session credentials: %w", err)
	}

	fmt.Fprintf(outfd, "-- MFA session refreshed, valid until: %s\n\n", expStr)
	return nil
}

func extractCredentials(c *types.Credentials) (accessKey, secretKey, sessionToken string, expiration time.Time) {
	if c == nil {
		return
	}
	if c.AccessKeyId != nil {
		accessKey = *c.AccessKeyId
	}
	if c.SecretAccessKey != nil {
		secretKey = *c.SecretAccessKey
	}
	if c.SessionToken != nil {
		sessionToken = *c.SessionToken
	}
	if c.Expiration != nil {
		expiration = *c.Expiration
	}
	return
}

// isValid returns true if the expiration string represents a future time.
func isValid(expiration string) bool {
	if expiration == "" {
		return false
	}
	exp, err := time.Parse(expirationLayout, expiration)
	if err != nil {
		return false
	}
	return exp.After(time.Now())
}

// acquireTOTP gets the TOTP code via totp_script or interactive prompt.
func acquireTOTP(cfg config.Config, outfd *os.File) (string, error) {
	if cfg.MFA.TOTPScript != "" {
		fmt.Fprintln(outfd, "-- Executing TOTP script...")
		code, err := runTOTPScript(cfg.MFA.TOTPScript)
		if err != nil {
			return "", err
		}
		code = strings.TrimSpace(code)
		fmt.Fprintln(outfd, code)

		if cfg.MFA.ClipboardCommand != "" {
			copyToClipboard(cfg.MFA.ClipboardCommand, code, outfd)
		}
		return code, nil
	}

	fmt.Fprintln(outfd, "-- No TOTP script configured. Please enter MFA code manually.")
	code, err := ui.Input("MFA code:", true)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

func runTOTPScript(script string) (string, error) {
	parts := strings.Fields(script)
	if len(parts) == 0 {
		return "", fmt.Errorf("totp_script is empty")
	}
	out, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return "", fmt.Errorf("totp_script failed: %w", err)
	}
	return string(out), nil
}

func copyToClipboard(cmd, text string, outfd *os.File) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}
	c := exec.Command(parts[0], parts[1:]...)
	c.Stdin = strings.NewReader(text)
	if err := c.Run(); err != nil {
		fmt.Fprintf(outfd, "-- Note: Clipboard copy failed (%v)\n", err)
	} else {
		fmt.Fprintln(outfd, "-- Copied to clipboard")
	}
}

// callGetSessionToken calls STS GetSessionToken using long-term credentials.
func callGetSessionToken(ltCreds awsconfig.Credentials, mfaDevice, tokenCode string) (*types.Credentials, error) {
	stsClient := sts.New(sts.Options{
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			ltCreds.AccessKeyID,
			ltCreds.SecretAccessKey,
			"",
		)),
		Region: "us-east-1",
	})

	resp, err := stsClient.GetSessionToken(context.Background(), &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int32(43200),
		SerialNumber:    aws.String(mfaDevice),
		TokenCode:       aws.String(tokenCode),
	})
	if err != nil {
		return nil, err
	}
	return resp.Credentials, nil
}
