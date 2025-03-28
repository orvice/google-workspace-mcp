package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Tool struct {
	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}

