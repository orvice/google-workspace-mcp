package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/tools"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Google Workspace MCP",
		Version: "0.0.1",
	}, nil)

	// Register all tools
	tools.RegisterAll(server)

	// Run server over stdio
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
