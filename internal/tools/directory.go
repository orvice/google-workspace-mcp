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

// GetUserInput defines input for get_user tool
type GetUserInput struct {
	UserKey string `json:"userKey" jsonschema:"required,description=User's primary email address or unique user ID"`
}

// GetUserOutput defines output for get_user tool
type GetUserOutput struct {
	User string `json:"user" jsonschema:"description=Detailed information about the user"`
}

// UpdateUserInput defines input for update_user tool
type UpdateUserInput struct {
	UserKey   string `json:"userKey" jsonschema:"required,description=User's primary email address or unique user ID"`
	FirstName string `json:"firstName,omitempty" jsonschema:"description=New first (given) name"`
	LastName  string `json:"lastName,omitempty" jsonschema:"description=New last (family) name"`
	Password  string `json:"password,omitempty" jsonschema:"description=New password for the user"`
	OrgUnit   string `json:"orgUnitPath,omitempty" jsonschema:"description=Organizational unit path to move the user to (e.g. /Sales)"`
}

// UpdateUserOutput defines output for update_user tool
type UpdateUserOutput struct {
	Result string `json:"result" jsonschema:"description=Result of user update"`
}

// DeleteUserInput defines input for delete_user tool
type DeleteUserInput struct {
	UserKey string `json:"userKey" jsonschema:"required,description=User's primary email address or unique user ID"`
}

// DeleteUserOutput defines output for delete_user tool
type DeleteUserOutput struct {
	Result string `json:"result" jsonschema:"description=Result of user deletion"`
}

// SuspendUserInput defines input for suspend_user tool
type SuspendUserInput struct {
	UserKey   string `json:"userKey" jsonschema:"required,description=User's primary email address or unique user ID"`
	Suspended bool   `json:"suspended" jsonschema:"required,description=Set true to suspend the account, false to restore it"`
}

// SuspendUserOutput defines output for suspend_user tool
type SuspendUserOutput struct {
	Result string `json:"result" jsonschema:"description=Result of suspending or restoring the user"`
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

// GetUser handles the get_user tool call
func GetUser(ctx context.Context, req *mcp.CallToolRequest, input GetUserInput) (*mcp.CallToolResult, GetUserOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, GetUserOutput{}, err
	}

	user, err := client.Users.Get(input.UserKey).Do()
	if err != nil {
		return nil, GetUserOutput{}, fmt.Errorf("failed to get user: %w", err)
	}

	resp := fmt.Sprintf("Email: %s\nName: %s\nID: %s\nAdmin: %t\nSuspended: %t\nOrg Unit: %s\nLast Login: %s\nCreated: %s",
		user.PrimaryEmail,
		user.Name.FullName,
		user.Id,
		user.IsAdmin,
		user.Suspended,
		user.OrgUnitPath,
		user.LastLoginTime,
		user.CreationTime)

	return nil, GetUserOutput{User: resp}, nil
}

// UpdateUser handles the update_user tool call
func UpdateUser(ctx context.Context, req *mcp.CallToolRequest, input UpdateUserInput) (*mcp.CallToolResult, UpdateUserOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, UpdateUserOutput{}, err
	}

	user := &admin.User{}
	if input.FirstName != "" || input.LastName != "" {
		user.Name = &admin.UserName{
			GivenName:  input.FirstName,
			FamilyName: input.LastName,
		}
	}
	if input.Password != "" {
		user.Password = input.Password
	}
	if input.OrgUnit != "" {
		user.OrgUnitPath = input.OrgUnit
	}

	updatedUser, err := client.Users.Patch(input.UserKey, user).Do()
	if err != nil {
		return nil, UpdateUserOutput{}, fmt.Errorf("failed to update user: %w", err)
	}

	resp := fmt.Sprintf("User updated successfully:\nEmail: %s\nName: %s\nOrg Unit: %s",
		updatedUser.PrimaryEmail,
		updatedUser.Name.FullName,
		updatedUser.OrgUnitPath)

	return nil, UpdateUserOutput{Result: resp}, nil
}

// DeleteUser handles the delete_user tool call
func DeleteUser(ctx context.Context, req *mcp.CallToolRequest, input DeleteUserInput) (*mcp.CallToolResult, DeleteUserOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, DeleteUserOutput{}, err
	}

	if err := client.Users.Delete(input.UserKey).Do(); err != nil {
		return nil, DeleteUserOutput{}, fmt.Errorf("failed to delete user: %w", err)
	}

	return nil, DeleteUserOutput{Result: fmt.Sprintf("User %s deleted successfully", input.UserKey)}, nil
}

// SuspendUser handles the suspend_user tool call
func SuspendUser(ctx context.Context, req *mcp.CallToolRequest, input SuspendUserInput) (*mcp.CallToolResult, SuspendUserOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, SuspendUserOutput{}, err
	}

	user := &admin.User{
		Suspended:       input.Suspended,
		ForceSendFields: []string{"Suspended"},
	}

	updatedUser, err := client.Users.Patch(input.UserKey, user).Do()
	if err != nil {
		return nil, SuspendUserOutput{}, fmt.Errorf("failed to update user suspension state: %w", err)
	}

	action := "restored"
	if input.Suspended {
		action = "suspended"
	}

	return nil, SuspendUserOutput{Result: fmt.Sprintf("User %s %s successfully", updatedUser.PrimaryEmail, action)}, nil
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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_user",
		Description: "Get detailed information about a specific user",
	}, GetUser)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_user",
		Description: "Update an existing user's name, password, or organizational unit",
	}, UpdateUser)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_user",
		Description: "Delete a user from Google Workspace",
	}, DeleteUser)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "suspend_user",
		Description: "Suspend or restore a user account",
	}, SuspendUser)
}
