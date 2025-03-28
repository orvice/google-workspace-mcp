# Google Workspace MCP

A Model Context Protocol (MCP) server for Google Workspace that provides tools for managing Google Workspace resources through the Admin SDK.

## Prerequisites

1. A Google Cloud Platform project with the Admin SDK API enabled
2. A service account with appropriate permissions
3. A Google Workspace admin user to impersonate

## Setup

### Service Account Configuration

1. Create a service account in the Google Cloud Console
2. Grant the service account appropriate permissions for Google Workspace Admin SDK
3. Create and download a JSON key file for the service account
4. Enable domain-wide delegation for the service account
5. Grant the necessary OAuth scopes to the service account in your Google Workspace Admin Console

## Environment Variables

The application requires the following environment variables to be set:

| Variable | Description |
|----------|-------------|
| `GOOGLE_SERVICE_ACCOUNT` | The path to the service account JSON key file |
| `GOOGLE_ADMIN_EMAIL` | The email address of the Google Workspace admin user to impersonate |

## Usage

### Build

```bash
make build
```

 
 ### config
 
 ```json
 {
  "mcpServers": {
     "googleworkspace-mcp": {
      "command": "/go/bin/google-workspace-mcp",
      "args": [],
      "env": {
        "GOOGLE_SERVICE_ACCOUNT": "test.json",
        "GOOGLE_ADMIN_EMAIL": "admin@yourdomain.com"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

> **Note:** Make sure your service account has the necessary API access enabled in Google Cloud Console (Admin SDK API, Gmail API, and Calendar API).

## Available Tools

- `directory_users` - List all users in your Google Workspace directory
- `list_gmail` - List recent Gmail messages (requires Gmail API access)
- `list_calendar_events` - List upcoming calendar events for a user (requires Calendar API access)

## Required OAuth Scopes

When setting up domain-wide delegation for your service account, ensure you grant the following OAuth scopes:

- `https://www.googleapis.com/auth/admin.directory.user` - For accessing directory user information
- `https://www.googleapis.com/auth/gmail.readonly` - For reading Gmail messages
- `https://www.googleapis.com/auth/calendar.readonly` - For reading calendar events
