package watcher

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// CDPClient communicates with Firefox via the WebDriver BiDi protocol.
// Firefox 73+ exposes BiDi on the --remote-debugging-port WebSocket endpoint.
type CDPClient struct {
	host string
	port int
}

// NewCDPClient returns a CDPClient targeting the given host and port.
func NewCDPClient(host string, port int) *CDPClient {
	return &CDPClient{host: host, port: port}
}

// IsReachable returns true if something is listening on the configured port.
func (c *CDPClient) IsReachable() bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.host, c.port), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// bidiMsg is a WebDriver BiDi command or response.
// Error is interface{} because BiDi sends it as a string (e.g. "unknown command").
type bidiMsg struct {
	ID     int            `json:"id"`
	Method string         `json:"method,omitempty"`
	Params map[string]any `json:"params,omitempty"`
	Type   string         `json:"type,omitempty"`
	Result map[string]any `json:"result,omitempty"`
	Error  interface{}    `json:"error,omitempty"`
}

// dialBiDi opens a WebSocket to the Firefox BiDi endpoint.
// When started with --remote-debugging-port, Firefox exposes an existing
// session at /session — no session.new needed.
func (c *CDPClient) dialBiDi(ctx context.Context) (*websocket.Conn, error) {
	wsURL := fmt.Sprintf("ws://%s:%d/session", c.host, c.port)
	header := http.Header{"Sec-WebSocket-Protocol": {"webdriverBiDi"}}
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, header)
	if err != nil {
		return nil, fmt.Errorf("BiDi WebSocket dial: %w", err)
	}
	return conn, nil
}

// sendCmd sends a BiDi command and waits for the matching response.
func sendCmd(ctx context.Context, conn *websocket.Conn, msg bidiMsg) (bidiMsg, error) {
	if err := conn.WriteJSON(msg); err != nil {
		return bidiMsg{}, fmt.Errorf("BiDi write: %w", err)
	}
	conn.SetReadDeadline(time.Now().Add(4 * time.Second))
	for {
		select {
		case <-ctx.Done():
			return bidiMsg{}, fmt.Errorf("BiDi command timed out")
		default:
		}
		var resp bidiMsg
		if err := conn.ReadJSON(&resp); err != nil {
			return bidiMsg{}, fmt.Errorf("BiDi read: %w", err)
		}
		if resp.ID == msg.ID {
			return resp, nil
		}
	}
}

// FindConsoleTab returns the BiDi context ID of the first open AWS console tab.
func (c *CDPClient) FindConsoleTab() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := c.dialBiDi(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	resp, err := sendCmd(ctx, conn, bidiMsg{
		ID:     2,
		Method: "browsingContext.getTree",
		Params: map[string]any{},
	})
	if err != nil {
		return "", err
	}
	if resp.Type == "error" || resp.Error != nil {
		return "", fmt.Errorf("browsingContext.getTree error: %v", resp.Error)
	}

	contexts, ok := resp.Result["contexts"].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected getTree response format")
	}
	return findAWSConsoleContext(contexts)
}

func findAWSConsoleContext(contexts []interface{}) (string, error) {
	for _, item := range contexts {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		url, _ := m["url"].(string)
		id, _ := m["context"].(string)
		if isAWSConsoleURL(url) && id != "" {
			return id, nil
		}
		if children, ok := m["children"].([]interface{}); ok {
			if found, err := findAWSConsoleContext(children); err == nil {
				return found, nil
			}
		}
	}
	return "", fmt.Errorf("no AWS console tab found")
}

// isAWSConsoleURL returns true if the URL is an AWS console URL.
func isAWSConsoleURL(u string) bool {
	return strings.HasPrefix(u, "https://") &&
		strings.Contains(u, ".console.aws.amazon.com")
}

// Evaluate executes expression in the given BiDi context via script.evaluate.
func (c *CDPClient) Evaluate(contextID, expression string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := c.dialBiDi(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := sendCmd(ctx, conn, bidiMsg{
		ID:     2,
		Method: "script.evaluate",
		Params: map[string]any{
			"expression":   expression,
			"target":       map[string]any{"context": contextID},
			"awaitPromise": true,
		},
	})
	if err != nil {
		return err
	}
	if resp.Type == "error" || resp.Error != nil {
		return fmt.Errorf("script.evaluate error: %v", resp.Error)
	}
	return nil
}

// RefreshSession executes the federation fetch inside an existing AWS console
// tab via BiDi. Returns an error if no console tab is found or BiDi fails.
func (c *CDPClient) RefreshSession(federationURL string) error {
	contextID, err := c.FindConsoleTab()
	if err != nil {
		return err
	}
	expr := fmt.Sprintf(
		`fetch(%q, {credentials:'include', mode:'no-cors', redirect:'follow'})`,
		federationURL,
	)
	return c.Evaluate(contextID, expr)
}
