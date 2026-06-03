package watcher

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"sync"
)

// Server is the local HTTP server used to perform invisible browser session
// refreshes. It maps short-lived tokens to federation URLs and serves a page
// that fetches the URL (setting the AWS session cookie inside the browser
// container) and then closes itself.
type Server struct {
	port   int
	tokens map[string]string
	mu     sync.Mutex
}

// refreshPageTmpl serves the refresh page. The federation URL is embedded in a
// <meta> tag so that html/template HTML-encodes it safely; JavaScript reads it
// back via getAttribute which reverses the HTML encoding.
var refreshPageTmpl = template.Must(template.New("refresh").Parse(`<!DOCTYPE html>
<html>
<head>
  <title>Refreshing AWS session...</title>
  <meta id="u" content="{{.}}">
</head>
<body>
<p>Refreshing AWS session...</p>
<script>
var u = document.getElementById('u').getAttribute('content');
fetch(u, {credentials: 'include', mode: 'no-cors', redirect: 'follow', keepalive: true})
  .catch(function() {});
window.close();
</script>
</body>
</html>`))

// StartServer starts a local HTTP server bound to a random available port on
// 127.0.0.1.
func StartServer() (*Server, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("watcher: failed to bind local port: %w", err)
	}
	srv := &Server{
		port:   ln.Addr().(*net.TCPAddr).Port,
		tokens: make(map[string]string),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/refresh", srv.handleRefresh)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go http.Serve(ln, mux) //nolint:errcheck
	return srv, nil
}

// Port returns the port the server is listening on.
func (s *Server) Port() int {
	return s.port
}

// RefreshURL generates a one-time token for the given federation URL and
// returns the localhost URL that will trigger the invisible refresh.
func (s *Server) RefreshURL(federationURL string) string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)

	s.mu.Lock()
	s.tokens[token] = federationURL
	s.mu.Unlock()

	return fmt.Sprintf("http://127.0.0.1:%d/refresh?t=%s", s.port, url.QueryEscape(token))
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("t")

	s.mu.Lock()
	federationURL, ok := s.tokens[token]
	if ok {
		delete(s.tokens, token) // one-time use
	}
	s.mu.Unlock()

	if !ok {
		http.Error(w, "invalid or expired token", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = refreshPageTmpl.Execute(w, federationURL)
}
