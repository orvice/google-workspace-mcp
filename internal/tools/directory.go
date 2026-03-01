package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	admin "google.golang.org/api/admin/directory/v1"
)

// ListUsersInput defines input for directory_users tool
type ListUsersInput struct {
	Domain string `json:"domain" jsonschema:"required,description=Domain to list users from"`
}

// ListUsersOutput defines output for directory_users tool
type ListUsersOutput struct {
	Users string `json:"users" jsonschema:"description=List of users with email and name"`
}

// CreateUserInput defines input for create_user tool
type CreateUserInput struct {
	Email     string `json:"email" jsonschema:"required,description=Email address for the new user"`
	FirstName string `json:"firstName" jsonschema:"required,description=First name of the user"`
	LastName  string `json:"lastName" jsonschema:"required,description=Last name of the user"`
	Password  string `json:"password" jsonschema:"required,description=Initial password for the user"`
}

// CreateUserOutput defines output for create_user tool
type CreateUserOutput struct {
	Result string `json:"result" jsonschema:"description=Result of user creation"`
}

// ListUsers handles the directory_users tool call
func ListUsers(ctx context.Context, req *mcp.CallToolRequest, input ListUsersInput) (*mcp.CallToolResult, ListUsersOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, ListUsersOutput{}, err
	}

	users, err := client.Users.List().Domain(input.Domain).Do()
	if err != nil {
		return nil, ListUsersOutput{}, err
	}

	var result string
	for _, user := range users.Users {
		result += fmt.Sprintf("Email: %s Name: %s\n", user.PrimaryEmail, user.Name.FullName)
	}

	return nil, ListUsersOutput{Users: result}, nil
}

// CreateUser handles the create_user tool call
func CreateUser(ctx context.Context, req *mcp.CallToolRequest, input CreateUserInput) (*mcp.CallToolResult, CreateUserOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, CreateUserOutput{}, err
	}

	// Create user object
	user := &admin.User{
		PrimaryEmail: input.Email,
		Name: &admin.UserName{
			GivenName:  input.FirstName,
			FamilyName: input.LastName,
			FullName:   input.FirstName + " " + input.LastName,
		},
		Password: input.Password,
	}

	// Create user in Google Workspace
	createdUser, err := client.Users.Insert(user).Do()
	if err != nil {
		return nil, CreateUserOutput{}, fmt.Errorf("failed to create user: %w", err)
	}

	resp := fmt.Sprintf("User created successfully:\nEmail: %s\nName: %s\nID: %s",
		createdUser.PrimaryEmail,
		createdUser.Name.FullName,
		createdUser.Id)

	return nil, CreateUserOutput{Result: resp}, nil
}

// RegisterDirectoryTools registers all directory-related tools with the MCP server
func RegisterDirectoryTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "directory_users",
		Description: "List Directory Users",
	}, ListUsers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_user",
		Description: "Create a new user in Google Workspace",
	}, CreateUser)
}
