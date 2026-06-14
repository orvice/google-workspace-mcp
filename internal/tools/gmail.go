package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
)

// ListGmailInput defines input for list_gmail tool
type ListGmailInput struct {
	Email string `json:"email" jsonschema:"Email address to access Gmail"`
}

// ListGmailOutput defines output for list_gmail tool
type ListGmailOutput struct {
	Messages string `json:"messages" jsonschema:"List of recent messages"`
}

// ListGmail handles the list_gmail tool call
func ListGmail(ctx context.Context, req *mcp.CallToolRequest, input ListGmailInput) (*mcp.CallToolResult, ListGmailOutput, error) {
	srv, err := utils.NewGmailClient(input.Email)
	if err != nil {
		return nil, ListGmailOutput{}, err
	}

	messages, err := srv.Users.Messages.List("me").MaxResults(10).Do()
	if err != nil {
		return nil, ListGmailOutput{}, err
	}

	var resp string
	for _, msg := range messages.Messages {
		fullMsg, err := srv.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue
		}

		var subject string
		for _, header := range fullMsg.Payload.Headers {
			if header.Name == "Subject" {
				subject = header.Value
				break
			}
		}

		resp += fmt.Sprintf("ID: %s, Subject: %s\n", msg.Id, subject)
	}

	return nil, ListGmailOutput{Messages: resp}, nil
}

// RegisterGmailTools registers all Gmail-related tools with the MCP server
func RegisterGmailTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_gmail",
		Description: "List Gmail Messages",
	}, ListGmail)
}
