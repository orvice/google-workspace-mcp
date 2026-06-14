package main

import (
	"context"
	"log"
	"os"

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

	// Select transport via MCP_TRANSPORT (stdio by default).
	if os.Getenv("MCP_TRANSPORT") == "http" {
		if err := runHTTP(server); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Run server over stdio
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
