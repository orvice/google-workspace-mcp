package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	"google.golang.org/api/calendar/v3"
)

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

func (d *Directory) CreateCalendarEvent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameters
	email, ok := request.Params.Arguments["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email is required")
	}
	
	summary, ok := request.Params.Arguments["summary"].(string)
	if !ok {
		return nil, fmt.Errorf("summary is required")
	}
	
	// Optional parameters
	description, _ := request.Params.Arguments["description"].(string)
	
	// Start and end times
	startTimeStr, ok := request.Params.Arguments["startTime"].(string)
	if !ok {
		return nil, fmt.Errorf("startTime is required (format: YYYY-MM-DDThh:mm:ss)")
	}
	
	endTimeStr, ok := request.Params.Arguments["endTime"].(string)
	if !ok {
		return nil, fmt.Errorf("endTime is required (format: YYYY-MM-DDThh:mm:ss)")
	}
	
	// Create calendar client
	srv, err := utils.NewCalendarClient(email)
	if err != nil {
		return nil, err
	}
	
	// Create calendar event
	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startTimeStr,
			TimeZone: "Asia/Shanghai", // Default timezone
		},
		End: &calendar.EventDateTime{
			DateTime: endTimeStr,
			TimeZone: "Asia/Shanghai", // Default timezone
		},
	}
	
	// Insert the event to the user's primary calendar
	createdEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar event: %w", err)
	}
	
	resp := fmt.Sprintf("Event created successfully:\nTitle: %s\nStart: %s\nEnd: %s\nLink: %s", 
		createdEvent.Summary,
		createdEvent.Start.DateTime,
		createdEvent.End.DateTime,
		createdEvent.HtmlLink)
	
	return mcp.NewToolResultText(resp), nil
}

// RegisterCalendarTools registers all calendar-related tools
func (d *Directory) RegisterCalendarTools() []Tool {
	return []Tool{
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
		{
			Tool: mcp.NewTool("create_calendar_event",
				mcp.WithDescription("Create a new calendar event"),
				mcp.WithString("email",
					mcp.Required(),
					mcp.Description("Email address to access calendar"),
				),
				mcp.WithString("summary",
					mcp.Required(),
					mcp.Description("Event title/summary"),
				),
				mcp.WithString("description",
					mcp.Description("Event description"),
				),
				mcp.WithString("startTime",
					mcp.Required(),
					mcp.Description("Start time in RFC3339 format (YYYY-MM-DDThh:mm:ss)"),
				),
				mcp.WithString("endTime",
					mcp.Required(),
					mcp.Description("End time in RFC3339 format (YYYY-MM-DDThh:mm:ss)"),
				),
			),
			Handler: d.CreateCalendarEvent,
		},
	}
}
