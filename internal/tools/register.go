package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// RegisterAll registers all tools with the MCP server.
// This function should be called after creating the server to add all
// Google Workspace tools (Directory, Gmail, Calendar, Drive, Sheets).
func RegisterAll(server *mcp.Server) {
	RegisterDirectoryTools(server)
	RegisterGmailTools(server)
	RegisterCalendarTools(server)
	RegisterDriveTools(server)
	RegisterSheetsTools(server)
	RegisterTasksTools(server)
}
