package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"go.orx.me/mcp/google-workspace/internal/tools"
)

const (
	version = "0.0.1"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Google Workspace MCP",
		version,
	)

	directory := tools.NewDirectory()
	for _, tool := range directory.Toolls() {
		s.AddTool(tool.Tool, tool.Handler)
	}

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
