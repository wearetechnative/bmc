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
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/wearetechnative/bmc/internal/config"
)

// OpenConsole opens the AWS Management Console for the given profile in the default browser.
// If service is non-empty (e.g. "ec2"), the console will open at that service page.
// If cfg.Console.FirefoxContainers is true, the URL is opened via the Granted Firefox extension.
func OpenConsole(profile, service string, bmcCfg config.Config) error {
	ctx := context.Background()

	awsCfg, err := awscfg.LoadDefaultConfig(ctx,
		awscfg.WithSharedConfigProfile(profile),
		awscfg.WithRegion(getRegionOrDefault(ctx, profile)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config for profile %q: %w", profile, err)
	}

	// Assume the role to get temporary credentials
	stsClient := sts.NewFromConfig(awsCfg)
	callerID, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get caller identity: %w", err)
	}

	// Get the credentials from the config
	creds, err := awsCfg.Credentials.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	_ = callerID

	// Build sign-in token via federation endpoint
	signinToken, err := getFederationToken(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)
	if err != nil {
		return fmt.Errorf("failed to get federation token: %w", err)
	}

	consoleURL := buildConsoleURL(signinToken, service)
	return openBrowser(consoleURL, profile, bmcCfg.Console.FirefoxContainers)
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

func buildConsoleURL(signinToken, service string) string {
	destination := "https://console.aws.amazon.com/"
	if service != "" {
		destination = "https://" + service + ".console.aws.amazon.com/"
	}

	params := url.Values{}
	params.Set("Action", "login")
	params.Set("Issuer", "bmc")
	params.Set("Destination", destination)
	params.Set("SigninToken", signinToken)

	return "https://signin.aws.amazon.com/federation?" + params.Encode()
}

func openBrowser(u, containerName string, firefoxContainers bool) error {
	if firefoxContainers {
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
