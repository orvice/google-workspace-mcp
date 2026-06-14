package tools

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	"google.golang.org/api/drive/v3"
)

// ListDriveFilesInput defines input for list_drive_files tool
type ListDriveFilesInput struct {
	Email      string `json:"email" jsonschema:"Email address to access Drive"`
	MaxResults int64  `json:"maxResults,omitempty" jsonschema:"Maximum number of files to return (default 10)"`
	FolderID   string `json:"folderId,omitempty" jsonschema:"Optional folder ID to list files from"`
}

// ListDriveFilesOutput defines output for list_drive_files tool
type ListDriveFilesOutput struct {
	Files string `json:"files" jsonschema:"List of files"`
}

// SearchDriveFilesInput defines input for search_drive_files tool
type SearchDriveFilesInput struct {
	Email      string `json:"email" jsonschema:"Email address to access Drive"`
	Query      string `json:"query" jsonschema:"Search query"`
	MaxResults int64  `json:"maxResults,omitempty" jsonschema:"Maximum number of files to return (default 10)"`
}

// SearchDriveFilesOutput defines output for search_drive_files tool
type SearchDriveFilesOutput struct {
	Files string `json:"files" jsonschema:"Search results"`
}

// GetDriveFileInput defines input for get_drive_file tool
type GetDriveFileInput struct {
	Email  string `json:"email" jsonschema:"Email address to access Drive"`
	FileID string `json:"fileId" jsonschema:"File ID to retrieve"`
}

// GetDriveFileOutput defines output for get_drive_file tool
type GetDriveFileOutput struct {
	FileInfo string `json:"fileInfo" jsonschema:"File information"`
}

// CreateDriveFolderInput defines input for create_drive_folder tool
type CreateDriveFolderInput struct {
	Email      string `json:"email" jsonschema:"Email address to access Drive"`
	Name       string `json:"name" jsonschema:"Folder name"`
	ParentID   string `json:"parentId,omitempty" jsonschema:"Parent folder ID (optional)"`
}

// CreateDriveFolderOutput defines output for create_drive_folder tool
type CreateDriveFolderOutput struct {
	Result string `json:"result" jsonschema:"Created folder information"`
}

// UploadDriveFileInput defines input for upload_drive_file tool
type UploadDriveFileInput struct {
	Email      string `json:"email" jsonschema:"Email address to access Drive"`
	FilePath   string `json:"filePath" jsonschema:"Local file path to upload"`
	Name       string `json:"name,omitempty" jsonschema:"Name for the file in Drive (uses filename if not specified)"`
	ParentID   string `json:"parentId,omitempty" jsonschema:"Parent folder ID (optional)"`
}

// UploadDriveFileOutput defines output for upload_drive_file tool
type UploadDriveFileOutput struct {
	Result string `json:"result" jsonschema:"Uploaded file information"`
}

// ShareDriveFileInput defines input for share_drive_file tool
type ShareDriveFileInput struct {
	Email      string `json:"email" jsonschema:"Email address to access Drive"`
	FileID     string `json:"fileId" jsonschema:"File ID to share"`
	UserEmail  string `json:"userEmail" jsonschema:"Email address to share with"`
	Role       string `json:"role,omitempty" jsonschema:"Permission role: reader, writer, commenter (default: reader)"`
}

// ShareDriveFileOutput defines output for share_drive_file tool
type ShareDriveFileOutput struct {
	Result string `json:"result" jsonschema:"Share result"`
}

// ListDriveFiles handles the list_drive_files tool call
func ListDriveFiles(ctx context.Context, req *mcp.CallToolRequest, input ListDriveFilesInput) (*mcp.CallToolResult, ListDriveFilesOutput, error) {
	srv, err := utils.NewDriveClient(input.Email)
	if err != nil {
		return nil, ListDriveFilesOutput{}, err
	}

	maxResults := input.MaxResults
	if maxResults == 0 {
		maxResults = 10
	}

	query := "trashed = false"
	if input.FolderID != "" {
		query = fmt.Sprintf("'%s' in parents and trashed = false", input.FolderID)
	}

	files, err := srv.Files.List().
		PageSize(maxResults).
		Q(query).
		Fields("files(id, name, mimeType, modifiedTime, size, webViewLink)").
		Do()
	if err != nil {
		return nil, ListDriveFilesOutput{}, err
	}

	var resp strings.Builder
	if len(files.Files) == 0 {
		resp.WriteString("No files found.\n")
	} else {
		resp.WriteString("Files:\n")
		for _, file := range files.Files {
			fileType := "File"
			if file.MimeType == "application/vnd.google-apps.folder" {
				fileType = "Folder"
			}
			resp.WriteString(fmt.Sprintf("[%s] %s\n  ID: %s\n  Type: %s\n  Modified: %s\n  Link: %s\n\n",
				fileType, file.Name, file.Id, file.MimeType, file.ModifiedTime, file.WebViewLink))
		}
	}

	return nil, ListDriveFilesOutput{Files: resp.String()}, nil
}

// SearchDriveFiles handles the search_drive_files tool call
func SearchDriveFiles(ctx context.Context, req *mcp.CallToolRequest, input SearchDriveFilesInput) (*mcp.CallToolResult, SearchDriveFilesOutput, error) {
	srv, err := utils.NewDriveClient(input.Email)
	if err != nil {
		return nil, SearchDriveFilesOutput{}, err
	}

	maxResults := input.MaxResults
	if maxResults == 0 {
		maxResults = 10
	}

	// Build search query
	query := fmt.Sprintf("name contains '%s' and trashed = false", input.Query)

	files, err := srv.Files.List().
		PageSize(maxResults).
		Q(query).
		Fields("files(id, name, mimeType, modifiedTime, size, webViewLink)").
		Do()
	if err != nil {
		return nil, SearchDriveFilesOutput{}, err
	}

	var resp strings.Builder
	if len(files.Files) == 0 {
		resp.WriteString(fmt.Sprintf("No files found matching '%s'.\n", input.Query))
	} else {
		resp.WriteString(fmt.Sprintf("Found %d file(s) matching '%s':\n\n", len(files.Files), input.Query))
		for _, file := range files.Files {
			fileType := "File"
			if file.MimeType == "application/vnd.google-apps.folder" {
				fileType = "Folder"
			}
			resp.WriteString(fmt.Sprintf("[%s] %s\n  ID: %s\n  Type: %s\n  Modified: %s\n  Link: %s\n\n",
				fileType, file.Name, file.Id, file.MimeType, file.ModifiedTime, file.WebViewLink))
		}
	}

	return nil, SearchDriveFilesOutput{Files: resp.String()}, nil
}

// GetDriveFile handles the get_drive_file tool call
func GetDriveFile(ctx context.Context, req *mcp.CallToolRequest, input GetDriveFileInput) (*mcp.CallToolResult, GetDriveFileOutput, error) {
	srv, err := utils.NewDriveClient(input.Email)
	if err != nil {
		return nil, GetDriveFileOutput{}, err
	}

	file, err := srv.Files.Get(input.FileID).
		Fields("id, name, mimeType, modifiedTime, size, webViewLink, description, owners, permissions").
		Do()
	if err != nil {
		return nil, GetDriveFileOutput{}, err
	}

	var resp strings.Builder
	resp.WriteString("File Information:\n")
	resp.WriteString(fmt.Sprintf("  Name: %s\n", file.Name))
	resp.WriteString(fmt.Sprintf("  ID: %s\n", file.Id))
	resp.WriteString(fmt.Sprintf("  Type: %s\n", file.MimeType))
	resp.WriteString(fmt.Sprintf("  Modified: %s\n", file.ModifiedTime))
	if file.Size > 0 {
		resp.WriteString(fmt.Sprintf("  Size: %d bytes\n", file.Size))
	}
	if file.Description != "" {
		resp.WriteString(fmt.Sprintf("  Description: %s\n", file.Description))
	}
	resp.WriteString(fmt.Sprintf("  Link: %s\n", file.WebViewLink))
	
	if len(file.Owners) > 0 {
		resp.WriteString(fmt.Sprintf("  Owner: %s\n", file.Owners[0].EmailAddress))
	}

	return nil, GetDriveFileOutput{FileInfo: resp.String()}, nil
}

// CreateDriveFolder handles the create_drive_folder tool call
func CreateDriveFolder(ctx context.Context, req *mcp.CallToolRequest, input CreateDriveFolderInput) (*mcp.CallToolResult, CreateDriveFolderOutput, error) {
	srv, err := utils.NewDriveClient(input.Email)
	if err != nil {
		return nil, CreateDriveFolderOutput{}, err
	}

	fileMetadata := &drive.File{
		Name:     input.Name,
		MimeType: "application/vnd.google-apps.folder",
	}

	if input.ParentID != "" {
		fileMetadata.Parents = []string{input.ParentID}
	}

	folder, err := srv.Files.Create(fileMetadata).
		Fields("id, name, webViewLink").
		Do()
	if err != nil {
		return nil, CreateDriveFolderOutput{}, fmt.Errorf("failed to create folder: %w", err)
	}

	resp := fmt.Sprintf("Folder created successfully:\n  Name: %s\n  ID: %s\n  Link: %s",
		folder.Name, folder.Id, folder.WebViewLink)

	return nil, CreateDriveFolderOutput{Result: resp}, nil
}

// UploadDriveFile handles the upload_drive_file tool call
func UploadDriveFile(ctx context.Context, req *mcp.CallToolRequest, input UploadDriveFileInput) (*mcp.CallToolResult, UploadDriveFileOutput, error) {
	srv, err := utils.NewDriveClient(input.Email)
	if err != nil {
		return nil, UploadDriveFileOutput{}, err
	}

	// Open the file
	file, err := os.Open(input.FilePath)
	if err != nil {
		return nil, UploadDriveFileOutput{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, UploadDriveFileOutput{}, fmt.Errorf("failed to get file info: %w", err)
	}

	// Use provided name or default to filename
	fileName := input.Name
	if fileName == "" {
		fileName = fileInfo.Name()
	}

	fileMetadata := &drive.File{
		Name: fileName,
	}

	if input.ParentID != "" {
		fileMetadata.Parents = []string{input.ParentID}
	}

	uploadedFile, err := srv.Files.Create(fileMetadata).
		Media(file).
		Fields("id, name, mimeType, size, webViewLink").
		Do()
	if err != nil {
		return nil, UploadDriveFileOutput{}, fmt.Errorf("failed to upload file: %w", err)
	}

	resp := fmt.Sprintf("File uploaded successfully:\n  Name: %s\n  ID: %s\n  Type: %s\n  Size: %d bytes\n  Link: %s",
		uploadedFile.Name, uploadedFile.Id, uploadedFile.MimeType, uploadedFile.Size, uploadedFile.WebViewLink)

	return nil, UploadDriveFileOutput{Result: resp}, nil
}

// ShareDriveFile handles the share_drive_file tool call
func ShareDriveFile(ctx context.Context, req *mcp.CallToolRequest, input ShareDriveFileInput) (*mcp.CallToolResult, ShareDriveFileOutput, error) {
	srv, err := utils.NewDriveClient(input.Email)
	if err != nil {
		return nil, ShareDriveFileOutput{}, err
	}

	role := input.Role
	if role == "" {
		role = "reader"
	}

	// Validate role
	validRoles := map[string]bool{
		"reader":    true,
		"writer":    true,
		"commenter": true,
		"owner":     true,
	}
	if !validRoles[role] {
		return nil, ShareDriveFileOutput{}, fmt.Errorf("invalid role: %s (must be reader, writer, commenter, or owner)", role)
	}

	permission := &drive.Permission{
		Type:         "user",
		Role:         role,
		EmailAddress: input.UserEmail,
	}

	_, err = srv.Permissions.Create(input.FileID, permission).
		SendNotificationEmail(true).
		Do()
	if err != nil {
		return nil, ShareDriveFileOutput{}, fmt.Errorf("failed to share file: %w", err)
	}

	// Get file info
	file, err := srv.Files.Get(input.FileID).Fields("name, webViewLink").Do()
	if err != nil {
		return nil, ShareDriveFileOutput{}, fmt.Errorf("failed to get file info: %w", err)
	}

	resp := fmt.Sprintf("File shared successfully:\n  File: %s\n  Shared with: %s\n  Role: %s\n  Link: %s",
		file.Name, input.UserEmail, role, file.WebViewLink)

	return nil, ShareDriveFileOutput{Result: resp}, nil
}

// RegisterDriveTools registers all Drive-related tools with the MCP server
func RegisterDriveTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_drive_files",
		Description: "List files in Google Drive",
	}, ListDriveFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_drive_files",
		Description: "Search for files in Google Drive",
	}, SearchDriveFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_drive_file",
		Description: "Get detailed information about a specific Drive file",
	}, GetDriveFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_drive_folder",
		Description: "Create a new folder in Google Drive",
	}, CreateDriveFolder)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "upload_drive_file",
		Description: "Upload a file to Google Drive",
	}, UploadDriveFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "share_drive_file",
		Description: "Share a Drive file with another user",
	}, ShareDriveFile)
}
