package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// runHTTP serves the MCP server over the Streamable HTTP transport.
//
// Configuration (environment variables):
//   - MCP_HTTP_ADDR:  listen address, defaults to ":8080"
//   - MCP_AUTH_TOKEN: when set, requests must carry "Authorization: Bearer <token>"
func runHTTP(server *mcp.Server) error {
	addr := os.Getenv("MCP_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	var h http.Handler = handler
	if token := os.Getenv("MCP_AUTH_TOKEN"); token != "" {
		h = withTokenAuth(h, token)
		log.Printf("Google Workspace MCP listening on %s (streamable HTTP, bearer auth enabled)", addr)
	} else {
		log.Printf("Google Workspace MCP listening on %s (streamable HTTP, no auth)", addr)
	}

	return http.ListenAndServe(addr, h)
}

// withTokenAuth wraps next, requiring an "Authorization: Bearer <token>" header
// that matches token. The comparison is constant-time to avoid leaking the
// token through response timing.
func withTokenAuth(next http.Handler, token string) http.Handler {
	expected := []byte("Bearer " + token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := []byte(r.Header.Get("Authorization"))
		if subtle.ConstantTimeCompare(got, expected) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
