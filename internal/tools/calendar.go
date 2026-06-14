package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	"google.golang.org/api/calendar/v3"
)

// ListCalendarEventsInput defines input for list_calendar_events tool
type ListCalendarEventsInput struct {
	Email string `json:"email" jsonschema:"Email address to access calendar"`
}

// ListCalendarEventsOutput defines output for list_calendar_events tool
type ListCalendarEventsOutput struct {
	Events string `json:"events" jsonschema:"List of upcoming events"`
}

// CreateCalendarEventInput defines input for create_calendar_event tool
type CreateCalendarEventInput struct {
	Email       string `json:"email" jsonschema:"Email address to access calendar"`
	Summary     string `json:"summary" jsonschema:"Event title/summary"`
	Description string `json:"description,omitempty" jsonschema:"Event description"`
	StartTime   string `json:"startTime" jsonschema:"Start time in RFC3339 format"`
	EndTime     string `json:"endTime" jsonschema:"End time in RFC3339 format"`
}

// CreateCalendarEventOutput defines output for create_calendar_event tool
type CreateCalendarEventOutput struct {
	Result string `json:"result" jsonschema:"Result of event creation"`
}

// ListCalendarEvents handles the list_calendar_events tool call
func ListCalendarEvents(ctx context.Context, req *mcp.CallToolRequest, input ListCalendarEventsInput) (*mcp.CallToolResult, ListCalendarEventsOutput, error) {
	srv, err := utils.NewCalendarClient(input.Email)
	if err != nil {
		return nil, ListCalendarEventsOutput{}, err
	}

	timeMin := time.Now().Format(time.RFC3339)
	timeMax := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)

	events, err := srv.Events.List("primary").TimeMin(timeMin).TimeMax(timeMax).MaxResults(10).OrderBy("startTime").SingleEvents(true).Do()
	if err != nil {
		return nil, ListCalendarEventsOutput{}, err
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

	return nil, ListCalendarEventsOutput{Events: resp}, nil
}

// CreateCalendarEvent handles the create_calendar_event tool call
func CreateCalendarEvent(ctx context.Context, req *mcp.CallToolRequest, input CreateCalendarEventInput) (*mcp.CallToolResult, CreateCalendarEventOutput, error) {
	srv, err := utils.NewCalendarClient(input.Email)
	if err != nil {
		return nil, CreateCalendarEventOutput{}, err
	}

	event := &calendar.Event{
		Summary:     input.Summary,
		Description: input.Description,
		Start: &calendar.EventDateTime{
			DateTime: input.StartTime,
			TimeZone: "Asia/Shanghai",
		},
		End: &calendar.EventDateTime{
			DateTime: input.EndTime,
			TimeZone: "Asia/Shanghai",
		},
	}

	createdEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, CreateCalendarEventOutput{}, fmt.Errorf("failed to create calendar event: %w", err)
	}

	resp := fmt.Sprintf("Event created successfully:\nTitle: %s\nStart: %s\nEnd: %s\nLink: %s",
		createdEvent.Summary,
		createdEvent.Start.DateTime,
		createdEvent.End.DateTime,
		createdEvent.HtmlLink)

	return nil, CreateCalendarEventOutput{Result: resp}, nil
}

// RegisterCalendarTools registers all calendar-related tools with the MCP server
func RegisterCalendarTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_calendar_events",
		Description: "List Calendar Events",
	}, ListCalendarEvents)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_calendar_event",
		Description: "Create a new calendar event",
	}, CreateCalendarEvent)
}
