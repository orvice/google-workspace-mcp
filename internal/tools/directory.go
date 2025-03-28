package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	admin "google.golang.org/api/admin/directory/v1"
)

type Directory struct {
	client *admin.Service
}

func NewDirectory() *Directory {
	client, err := utils.DefaultClient()
	if err != nil {
		log.Fatal("failed to create admin service", "error", err)
	}
	return &Directory{client: client}
}

func (d *Directory) Users(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	domain, ok := request.Params.Arguments["domain"].(string)
	if !ok {
		return nil, fmt.Errorf("domain is required")
	}
	users, err := d.client.Users.List().Domain(domain).Do()
	if err != nil {
		return nil, err
	}
	var resp string
	for _, user := range users.Users {
		resp += fmt.Sprintf("Email: %s Name: %s \n", user.PrimaryEmail, user.Name.FullName)
	}
	return mcp.NewToolResultText(resp), nil
}

func (d *Directory) Toolls() []Tool {
	return []Tool{
		{
			Tool: mcp.NewTool("directory_users",
				mcp.WithDescription("List Directory Users"),
				mcp.WithString("domain",
					mcp.Required(),
					mcp.Description("domain"),
				),
			),
			Handler: d.Users,
		},
	}
}
