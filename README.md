# Google Workspace MCP

A Model Context Protocol (MCP) server for Google Workspace that provides tools for managing Google Workspace resources through the Admin SDK, Gmail, Calendar, and Drive APIs.

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

> **Note:** Make sure your service account has the necessary API access enabled in Google Cloud Console (Admin SDK API, Gmail API, Calendar API, and Drive API).

## Available Tools

### Directory Tools
- `directory_users` - List all users in your Google Workspace directory
- `create_user` - Create a new user in Google Workspace

### Gmail Tools
- `list_gmail` - List recent Gmail messages (requires Gmail API access)

### Calendar Tools
- `list_calendar_events` - List upcoming calendar events for a user (requires Calendar API access)
- `create_calendar_event` - Create a new calendar event (requires Calendar API access)

### Drive Tools
- `list_drive_files` - List files in Google Drive (requires Drive API access)
- `search_drive_files` - Search for files in Google Drive (requires Drive API access)
- `get_drive_file` - Get detailed information about a specific Drive file (requires Drive API access)
- `create_drive_folder` - Create a new folder in Google Drive (requires Drive API access)
- `upload_drive_file` - Upload a file to Google Drive (requires Drive API access)
- `share_drive_file` - Share a Drive file with another user (requires Drive API access)

## Required OAuth Scopes

When setting up domain-wide delegation for your service account, ensure you grant the following OAuth scopes:

- `https://www.googleapis.com/auth/admin.directory.user` - For accessing and managing directory user information
- `https://www.googleapis.com/auth/gmail.readonly` - For reading Gmail messages
- `https://www.googleapis.com/auth/calendar` - For reading and writing calendar events
- `https://www.googleapis.com/auth/drive` - For full access to Google Drive (reading, writing, and managing files)