package tools

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (d *Directory) CreateUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameters
	email, ok := request.Params.Arguments["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email is required")
	}
	
	firstName, ok := request.Params.Arguments["firstName"].(string)
	if !ok {
		return nil, fmt.Errorf("firstName is required")
	}
	
	lastName, ok := request.Params.Arguments["lastName"].(string)
	if !ok {
		return nil, fmt.Errorf("lastName is required")
	}
	
	password, ok := request.Params.Arguments["password"].(string)
	if !ok {
		return nil, fmt.Errorf("password is required")
	}
	
	// Create user object
	user := &admin.User{
		PrimaryEmail: email,
		Name: &admin.UserName{
			GivenName:  firstName,
			FamilyName: lastName,
			FullName:   firstName + " " + lastName,
		},
		Password: password,
	}
	
	// Create user in Google Workspace
	createdUser, err := d.client.Users.Insert(user).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	resp := fmt.Sprintf("User created successfully:\nEmail: %s\nName: %s\nID: %s", 
		createdUser.PrimaryEmail, 
		createdUser.Name.FullName,
		createdUser.Id)
	
	return mcp.NewToolResultText(resp), nil
}

func (d *Directory) ListEmail(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	email, ok := request.Params.Arguments["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email is required")
	}
	srv, err := utils.NewGmailClient(email)
	if err != nil {
		return nil, err
	}
	// Get list of message IDs
	messages, err := srv.Users.Messages.List("me").MaxResults(10).Do()
	if err != nil {
		return nil, err
	}

	var resp string
	// For each message ID, get the full message details
	for _, msg := range messages.Messages {
		// Get the full message
		fullMsg, err := srv.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue
		}

		// Extract subject from headers
		var subject string
		for _, header := range fullMsg.Payload.Headers {
			if header.Name == "Subject" {
				subject = header.Value
				break
			}
		}

		resp += fmt.Sprintf("ID: %s, Subject: %s\n", msg.Id, subject)
	}

	return mcp.NewToolResultText(resp), nil
}

func (d *Directory) ListCalendarEvents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	email, ok := request.Params.Arguments["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email is required")
	}
	
	// Create calendar client
	srv, err := utils.NewCalendarClient(email)
	if err != nil {
		return nil, err
	}
	
	// Get calendar events for primary calendar
	timeMin := time.Now().Format(time.RFC3339)
	timeMax := time.Now().AddDate(0, 0, 7).Format(time.RFC3339) // Get events for next 7 days
	
	events, err := srv.Events.List("primary").TimeMin(timeMin).TimeMax(timeMax).MaxResults(10).OrderBy("startTime").SingleEvents(true).Do()
	if err != nil {
		return nil, err
	}
	
	var resp string
	if len(events.Items) == 0 {
		resp = "No upcoming events found."
	} else {
		resp = "Upcoming events:\n"
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			resp += fmt.Sprintf("%s (%s)\n", item.Summary, date)
		}
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
		{
			Tool: mcp.NewTool("create_user",
				mcp.WithDescription("Create a new user in Google Workspace"),
				mcp.WithString("email",
					mcp.Required(),
					mcp.Description("Email address for the new user"),
				),
				mcp.WithString("firstName",
					mcp.Required(),
					mcp.Description("First name of the user"),
				),
				mcp.WithString("lastName",
					mcp.Required(),
					mcp.Description("Last name of the user"),
				),
				mcp.WithString("password",
					mcp.Required(),
					mcp.Description("Initial password for the user"),
				),
			),
			Handler: d.CreateUser,
		},
		{
			Tool: mcp.NewTool("list_gmail",
				mcp.WithDescription("List Gmail Messages"),
				mcp.WithString("email",
					mcp.Required(),
					mcp.Description("Email address to access Gmail"),
				),
			),
			Handler: d.ListEmail,
		},
		{
			Tool: mcp.NewTool("list_calendar_events",
				mcp.WithDescription("List Calendar Events"),
				mcp.WithString("email",
					mcp.Required(),
					mcp.Description("Email address to access calendar"),
				),
			),
			Handler: d.ListCalendarEvents,
		},
	}
}
