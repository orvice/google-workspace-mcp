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

### Transport (optional)

By default the server runs over stdio. Set `MCP_TRANSPORT=http` to serve over the
Streamable HTTP transport instead.

| Variable | Description |
|----------|-------------|
| `MCP_TRANSPORT` | Transport to use: `stdio` (default) or `http` |
| `MCP_HTTP_ADDR` | Listen address for HTTP transport (default `:8080`) |
| `MCP_AUTH_TOKEN` | When set, HTTP requests must include `Authorization: Bearer <token>` |

```bash
MCP_TRANSPORT=http MCP_HTTP_ADDR=:8080 MCP_AUTH_TOKEN=your-secret-token google-workspace-mcp
```

Clients then connect to `http://<host>:8080/` and send the bearer token, e.g.:

```
Authorization: Bearer your-secret-token
```

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

> **Note:** Make sure your service account has the necessary API access enabled in Google Cloud Console (Admin SDK API, Gmail API, Calendar API, Drive API, Sheets API, and Tasks API).

## Available Tools

### Directory Tools
- `directory_users` - List all users in your Google Workspace directory
- `create_user` - Create a new user in Google Workspace
- `get_user` - Get detailed information about a specific user
- `update_user` - Update an existing user's name, password, or organizational unit
- `delete_user` - Delete a user from Google Workspace
- `suspend_user` - Suspend or restore a user account

### Group Tools
- `list_groups` - List groups in a domain
- `get_group` - Get detailed information about a specific group
- `create_group` - Create a new group in Google Workspace
- `delete_group` - Delete a group from Google Workspace
- `list_group_members` - List members of a group
- `add_group_member` - Add a member to a group
- `remove_group_member` - Remove a member from a group

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

### Sheets Tools
- `list_spreadsheets` - List Google Sheets spreadsheets in Drive (requires Sheets API access)
- `get_spreadsheet` - Get detailed information about a spreadsheet including its sheets (requires Sheets API access)
- `read_sheet_range` - Read data from a specific range in a spreadsheet (requires Sheets API access)
- `write_sheet_range` - Write data to a specific range in a spreadsheet (requires Sheets API access)
- `append_sheet_rows` - Append rows of data to a spreadsheet (requires Sheets API access)
- `create_spreadsheet` - Create a new Google Sheets spreadsheet (requires Sheets API access)

### Tasks Tools
- `list_task_lists` - List all Google Tasks task lists for a user (requires Tasks API access)
- `list_tasks` - List all tasks in a specific task list (requires Tasks API access)
- `create_task` - Create a new task in a task list (requires Tasks API access)
- `update_task` - Update an existing task (requires Tasks API access)
- `delete_task` - Delete a task from a task list (requires Tasks API access)
- `complete_task` - Mark a task as completed (requires Tasks API access)

## Required OAuth Scopes

When setting up domain-wide delegation for your service account, ensure you grant the following OAuth scopes:

- `https://www.googleapis.com/auth/admin.directory.user` - For accessing and managing directory user information
- `https://www.googleapis.com/auth/admin.directory.group` - For managing groups
- `https://www.googleapis.com/auth/admin.directory.group.member` - For managing group members
- `https://www.googleapis.com/auth/gmail.readonly` - For reading Gmail messages
- `https://www.googleapis.com/auth/calendar` - For reading and writing calendar events
- `https://www.googleapis.com/auth/drive` - For full access to Google Drive (reading, writing, and managing files)
- `https://www.googleapis.com/auth/spreadsheets` - For reading and writing Google Sheets spreadsheets
- `https://www.googleapis.com/auth/tasks` - For reading and writing Google Tasks