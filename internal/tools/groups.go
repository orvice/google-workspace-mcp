package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	admin "google.golang.org/api/admin/directory/v1"
)

// ListGroupsInput defines input for list_groups tool
type ListGroupsInput struct {
	Domain string `json:"domain" jsonschema:"required,description=Domain to list groups from"`
}

// ListGroupsOutput defines output for list_groups tool
type ListGroupsOutput struct {
	Groups string `json:"groups" jsonschema:"description=List of groups with email and name"`
}

// GetGroupInput defines input for get_group tool
type GetGroupInput struct {
	GroupKey string `json:"groupKey" jsonschema:"required,description=Group's email address or unique group ID"`
}

// GetGroupOutput defines output for get_group tool
type GetGroupOutput struct {
	Group string `json:"group" jsonschema:"description=Detailed information about the group"`
}

// CreateGroupInput defines input for create_group tool
type CreateGroupInput struct {
	Email       string `json:"email" jsonschema:"required,description=Email address for the new group"`
	Name        string `json:"name" jsonschema:"required,description=Display name of the group"`
	Description string `json:"description,omitempty" jsonschema:"description=Description of the group"`
}

// CreateGroupOutput defines output for create_group tool
type CreateGroupOutput struct {
	Result string `json:"result" jsonschema:"description=Result of group creation"`
}

// DeleteGroupInput defines input for delete_group tool
type DeleteGroupInput struct {
	GroupKey string `json:"groupKey" jsonschema:"required,description=Group's email address or unique group ID"`
}

// DeleteGroupOutput defines output for delete_group tool
type DeleteGroupOutput struct {
	Result string `json:"result" jsonschema:"description=Result of group deletion"`
}

// ListGroupMembersInput defines input for list_group_members tool
type ListGroupMembersInput struct {
	GroupKey string `json:"groupKey" jsonschema:"required,description=Group's email address or unique group ID"`
}

// ListGroupMembersOutput defines output for list_group_members tool
type ListGroupMembersOutput struct {
	Members string `json:"members" jsonschema:"description=List of group members"`
}

// AddGroupMemberInput defines input for add_group_member tool
type AddGroupMemberInput struct {
	GroupKey string `json:"groupKey" jsonschema:"required,description=Group's email address or unique group ID"`
	Email    string `json:"email" jsonschema:"required,description=Email address of the member to add"`
	Role     string `json:"role,omitempty" jsonschema:"description=Member role: MEMBER, MANAGER, or OWNER (defaults to MEMBER)"`
}

// AddGroupMemberOutput defines output for add_group_member tool
type AddGroupMemberOutput struct {
	Result string `json:"result" jsonschema:"description=Result of adding the member"`
}

// RemoveGroupMemberInput defines input for remove_group_member tool
type RemoveGroupMemberInput struct {
	GroupKey  string `json:"groupKey" jsonschema:"required,description=Group's email address or unique group ID"`
	MemberKey string `json:"memberKey" jsonschema:"required,description=Email address or ID of the member to remove"`
}

// RemoveGroupMemberOutput defines output for remove_group_member tool
type RemoveGroupMemberOutput struct {
	Result string `json:"result" jsonschema:"description=Result of removing the member"`
}

// ListGroups handles the list_groups tool call
func ListGroups(ctx context.Context, req *mcp.CallToolRequest, input ListGroupsInput) (*mcp.CallToolResult, ListGroupsOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, ListGroupsOutput{}, err
	}

	groups, err := client.Groups.List().Domain(input.Domain).Do()
	if err != nil {
		return nil, ListGroupsOutput{}, fmt.Errorf("failed to list groups: %w", err)
	}

	var resp string
	if len(groups.Groups) == 0 {
		resp = "No groups found."
	} else {
		for _, g := range groups.Groups {
			resp += fmt.Sprintf("Email: %s Name: %s Members: %d\n", g.Email, g.Name, g.DirectMembersCount)
		}
	}

	return nil, ListGroupsOutput{Groups: resp}, nil
}

// GetGroup handles the get_group tool call
func GetGroup(ctx context.Context, req *mcp.CallToolRequest, input GetGroupInput) (*mcp.CallToolResult, GetGroupOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, GetGroupOutput{}, err
	}

	group, err := client.Groups.Get(input.GroupKey).Do()
	if err != nil {
		return nil, GetGroupOutput{}, fmt.Errorf("failed to get group: %w", err)
	}

	resp := fmt.Sprintf("Email: %s\nName: %s\nID: %s\nDescription: %s\nDirect Members: %d",
		group.Email,
		group.Name,
		group.Id,
		group.Description,
		group.DirectMembersCount)

	return nil, GetGroupOutput{Group: resp}, nil
}

// CreateGroup handles the create_group tool call
func CreateGroup(ctx context.Context, req *mcp.CallToolRequest, input CreateGroupInput) (*mcp.CallToolResult, CreateGroupOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, CreateGroupOutput{}, err
	}

	group := &admin.Group{
		Email:       input.Email,
		Name:        input.Name,
		Description: input.Description,
	}

	createdGroup, err := client.Groups.Insert(group).Do()
	if err != nil {
		return nil, CreateGroupOutput{}, fmt.Errorf("failed to create group: %w", err)
	}

	resp := fmt.Sprintf("Group created successfully:\nEmail: %s\nName: %s\nID: %s",
		createdGroup.Email,
		createdGroup.Name,
		createdGroup.Id)

	return nil, CreateGroupOutput{Result: resp}, nil
}

// DeleteGroup handles the delete_group tool call
func DeleteGroup(ctx context.Context, req *mcp.CallToolRequest, input DeleteGroupInput) (*mcp.CallToolResult, DeleteGroupOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, DeleteGroupOutput{}, err
	}

	if err := client.Groups.Delete(input.GroupKey).Do(); err != nil {
		return nil, DeleteGroupOutput{}, fmt.Errorf("failed to delete group: %w", err)
	}

	return nil, DeleteGroupOutput{Result: fmt.Sprintf("Group %s deleted successfully", input.GroupKey)}, nil
}

// ListGroupMembers handles the list_group_members tool call
func ListGroupMembers(ctx context.Context, req *mcp.CallToolRequest, input ListGroupMembersInput) (*mcp.CallToolResult, ListGroupMembersOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, ListGroupMembersOutput{}, err
	}

	members, err := client.Members.List(input.GroupKey).Do()
	if err != nil {
		return nil, ListGroupMembersOutput{}, fmt.Errorf("failed to list group members: %w", err)
	}

	var resp string
	if len(members.Members) == 0 {
		resp = "No members found."
	} else {
		for _, m := range members.Members {
			resp += fmt.Sprintf("Email: %s Role: %s Status: %s\n", m.Email, m.Role, m.Status)
		}
	}

	return nil, ListGroupMembersOutput{Members: resp}, nil
}

// AddGroupMember handles the add_group_member tool call
func AddGroupMember(ctx context.Context, req *mcp.CallToolRequest, input AddGroupMemberInput) (*mcp.CallToolResult, AddGroupMemberOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, AddGroupMemberOutput{}, err
	}

	role := input.Role
	if role == "" {
		role = "MEMBER"
	}

	member := &admin.Member{
		Email: input.Email,
		Role:  role,
	}

	addedMember, err := client.Members.Insert(input.GroupKey, member).Do()
	if err != nil {
		return nil, AddGroupMemberOutput{}, fmt.Errorf("failed to add group member: %w", err)
	}

	resp := fmt.Sprintf("Member added successfully:\nGroup: %s\nEmail: %s\nRole: %s",
		input.GroupKey,
		addedMember.Email,
		addedMember.Role)

	return nil, AddGroupMemberOutput{Result: resp}, nil
}

// RemoveGroupMember handles the remove_group_member tool call
func RemoveGroupMember(ctx context.Context, req *mcp.CallToolRequest, input RemoveGroupMemberInput) (*mcp.CallToolResult, RemoveGroupMemberOutput, error) {
	client, err := utils.DefaultClient()
	if err != nil {
		return nil, RemoveGroupMemberOutput{}, err
	}

	if err := client.Members.Delete(input.GroupKey, input.MemberKey).Do(); err != nil {
		return nil, RemoveGroupMemberOutput{}, fmt.Errorf("failed to remove group member: %w", err)
	}

	return nil, RemoveGroupMemberOutput{Result: fmt.Sprintf("Member %s removed from group %s successfully", input.MemberKey, input.GroupKey)}, nil
}

// RegisterGroupsTools registers all group-related tools with the MCP server
func RegisterGroupsTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_groups",
		Description: "List groups in a domain",
	}, ListGroups)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_group",
		Description: "Get detailed information about a specific group",
	}, GetGroup)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_group",
		Description: "Create a new group in Google Workspace",
	}, CreateGroup)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_group",
		Description: "Delete a group from Google Workspace",
	}, DeleteGroup)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_group_members",
		Description: "List members of a group",
	}, ListGroupMembers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_group_member",
		Description: "Add a member to a group",
	}, AddGroupMember)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_group_member",
		Description: "Remove a member from a group",
	}, RemoveGroupMember)
}
