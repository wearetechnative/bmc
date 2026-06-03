package awsops

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/wearetechnative/bmc/internal/config"
)

// BuildFederationURL loads credentials for the given profile, calls the AWS
// federation endpoint, and returns a ready-to-use signin URL along with the
// credential expiry time. The expiry reflects the underlying STS session
// duration (typically 1 hour for role-assumed credentials).
func BuildFederationURL(profile, service string, bmcCfg config.Config) (string, time.Time, error) {
	ctx := context.Background()

	awsCfg, err := awscfg.LoadDefaultConfig(ctx,
		awscfg.WithSharedConfigProfile(profile),
		awscfg.WithRegion(getRegionOrDefault(ctx, profile)),
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to load AWS config for profile %q: %w", profile, err)
	}

	creds, err := awsCfg.Credentials.Retrieve(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	signinToken, err := getFederationToken(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get federation token: %w", err)
	}

	consoleURL := buildConsoleURL(signinToken, service, awsCfg.Region)
	return consoleURL, creds.Expires, nil
}

// OpenConsole opens the AWS Management Console for the given profile in the default browser.
// If service is non-empty (e.g. "ec2"), the console will open at that service page.
// If cfg.Console.FirefoxContainers is true, the URL is opened via the Granted Firefox extension.
func OpenConsole(profile, service string, bmcCfg config.Config) error {
	consoleURL, _, err := BuildFederationURL(profile, service, bmcCfg)
	if err != nil {
		return err
	}
	return openBrowser(consoleURL, profile, bmcCfg.Console)
}

// OpenURLInBrowser opens any URL in the appropriate browser based on the
// console config (Firefox container, Chrome profile, or system default).
// containerName is used as the container/profile name for browser isolation.
func OpenURLInBrowser(u, containerName string, consoleCfg config.ConsoleConfig) error {
	return openBrowser(u, containerName, consoleCfg)
}

type federationTokenResponse struct {
	SigninToken string
}

func getFederationToken(accessKey, secretKey, sessionToken string) (string, error) {
	session := map[string]string{
		"sessionId":    accessKey,
		"sessionKey":   secretKey,
		"sessionToken": sessionToken,
	}
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("Action", "getSigninToken")
	params.Set("DurationSeconds", fmt.Sprintf("%d", int(time.Hour.Seconds()*12)))
	params.Set("Session", string(sessionJSON))

	endpoint := "https://signin.aws.amazon.com/federation?" + params.Encode()
	resp, err := http.Get(endpoint) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("federation endpoint unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("federation endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp federationTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode federation response: %w", err)
	}
	return tokenResp.SigninToken, nil
}

func buildConsoleURL(signinToken, service, region string) string {
	destination := "https://" + region + ".console.aws.amazon.com/"
	if service != "" {
		if strings.Contains(service, "/") {
			destination = "https://" + region + ".console.aws.amazon.com/" + service
		} else {
			destination = "https://" + region + ".console.aws.amazon.com/" + service + "/home"
		}
	}

	params := url.Values{}
	params.Set("Action", "login")
	params.Set("Issuer", "bmc")
	params.Set("Destination", destination)
	params.Set("SigninToken", signinToken)

	return "https://signin.aws.amazon.com/federation?" + params.Encode()
}

func openBrowser(u, containerName string, consoleCfg config.ConsoleConfig) error {
	if consoleCfg.ChromeProfiles && consoleCfg.FirefoxContainers {
		return fmt.Errorf("config conflict: chrome_profiles and firefox_containers cannot both be true — disable one in ~/.config/bmc/config.json")
	}
	if consoleCfg.ChromeProfiles {
		return openChromeProfile(u, containerName, consoleCfg)
	}
	if consoleCfg.FirefoxContainers {
		firefoxPath, err := exec.LookPath("firefox")
		if err != nil {
			return fmt.Errorf("firefox not found in PATH (required for firefox_containers): install Firefox or set firefox_containers = false in config")
		}
		grantedURL := fmt.Sprintf("ext+granted-containers:name=%s&url=%s&color=blue&icon=circle",
			url.QueryEscape(containerName), url.QueryEscape(u))
		cmd := exec.Command(firefoxPath, grantedURL)
		cmd.Stderr = nil
		return cmd.Start()
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", u)
	case "darwin":
		cmd = exec.Command("open", u)
	default:
		fmt.Println("Open this URL in your browser:")
		fmt.Println(u)
		return nil
	}
	cmd.Stderr = nil
	return cmd.Start()
}

func getRegionOrDefault(ctx context.Context, profile string) string {
	cfg, err := awscfg.LoadDefaultConfig(ctx, awscfg.WithSharedConfigProfile(profile))
	if err != nil {
		return "eu-west-1"
	}
	region := cfg.Region
	if region == "" {
		return "eu-west-1"
	}
	return region
}

// UseOrSelectProfile returns the current AWS_PROFILE if set, or prompts for selection.
func ProfileForConsole(profile string, profiles interface{ FindByName(string) (string, bool) }) string {
	return strings.TrimSpace(profile)
}
