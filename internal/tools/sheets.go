package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	"google.golang.org/api/sheets/v4"
)

// ListSpreadsheetsInput defines input for list_spreadsheets tool
type ListSpreadsheetsInput struct {
	Email      string `json:"email" jsonschema:"required,description=Email address to access Sheets"`
	MaxResults int64  `json:"maxResults,omitempty" jsonschema:"description=Maximum number of spreadsheets to return (default 10)"`
}

// ListSpreadsheetsOutput defines output for list_spreadsheets tool
type ListSpreadsheetsOutput struct {
	Spreadsheets string `json:"spreadsheets" jsonschema:"description=List of spreadsheets"`
}

// GetSpreadsheetInput defines input for get_spreadsheet tool
type GetSpreadsheetInput struct {
	Email         string `json:"email" jsonschema:"required,description=Email address to access Sheets"`
	SpreadsheetID string `json:"spreadsheetId" jsonschema:"required,description=Spreadsheet ID"`
}

// GetSpreadsheetOutput defines output for get_spreadsheet tool
type GetSpreadsheetOutput struct {
	Info string `json:"info" jsonschema:"description=Spreadsheet information"`
}

// ReadSheetRangeInput defines input for read_sheet_range tool
type ReadSheetRangeInput struct {
	Email         string `json:"email" jsonschema:"required,description=Email address to access Sheets"`
	SpreadsheetID string `json:"spreadsheetId" jsonschema:"required,description=Spreadsheet ID"`
	Range         string `json:"range" jsonschema:"required,description=A1 notation range (e.g. Sheet1!A1:B10)"`
}

// ReadSheetRangeOutput defines output for read_sheet_range tool
type ReadSheetRangeOutput struct {
	Data string `json:"data" jsonschema:"description=Cell data in table format"`
}

// WriteSheetRangeInput defines input for write_sheet_range tool
type WriteSheetRangeInput struct {
	Email         string          `json:"email" jsonschema:"required,description=Email address to access Sheets"`
	SpreadsheetID string          `json:"spreadsheetId" jsonschema:"required,description=Spreadsheet ID"`
	Range         string          `json:"range" jsonschema:"required,description=A1 notation range (e.g. Sheet1!A1:B10)"`
	Values        [][]interface{} `json:"values" jsonschema:"required,description=2D array of values to write"`
}

// WriteSheetRangeOutput defines output for write_sheet_range tool
type WriteSheetRangeOutput struct {
	Result string `json:"result" jsonschema:"description=Write operation result"`
}


// AppendSheetRowsInput defines input for append_sheet_rows tool
type AppendSheetRowsInput struct {
	Email         string          `json:"email" jsonschema:"required,description=Email address to access Sheets"`
	SpreadsheetID string          `json:"spreadsheetId" jsonschema:"required,description=Spreadsheet ID"`
	Range         string          `json:"range" jsonschema:"required,description=A1 notation range to append after (e.g. Sheet1!A:A)"`
	Values        [][]interface{} `json:"values" jsonschema:"required,description=2D array of rows to append"`
}

// AppendSheetRowsOutput defines output for append_sheet_rows tool
type AppendSheetRowsOutput struct {
	Result string `json:"result" jsonschema:"description=Append operation result"`
}

// CreateSpreadsheetInput defines input for create_spreadsheet tool
type CreateSpreadsheetInput struct {
	Email      string   `json:"email" jsonschema:"required,description=Email address to access Sheets"`
	Title      string   `json:"title" jsonschema:"required,description=Spreadsheet title"`
	SheetNames []string `json:"sheetNames,omitempty" jsonschema:"description=Optional list of sheet names to create"`
}

// CreateSpreadsheetOutput defines output for create_spreadsheet tool
type CreateSpreadsheetOutput struct {
	Result string `json:"result" jsonschema:"description=Created spreadsheet information"`
}

// ListSpreadsheets handles the list_spreadsheets tool call
func ListSpreadsheets(ctx context.Context, req *mcp.CallToolRequest, input ListSpreadsheetsInput) (*mcp.CallToolResult, ListSpreadsheetsOutput, error) {
	// Use Drive API to list spreadsheets by MIME type
	driveSrv, err := utils.NewDriveClient(input.Email)
	if err != nil {
		return nil, ListSpreadsheetsOutput{}, err
	}

	maxResults := input.MaxResults
	if maxResults == 0 {
		maxResults = 10
	}

	// Query for Google Sheets MIME type
	query := "mimeType = 'application/vnd.google-apps.spreadsheet' and trashed = false"

	files, err := driveSrv.Files.List().
		PageSize(maxResults).
		Q(query).
		Fields("files(id, name, modifiedTime)").
		Do()
	if err != nil {
		return nil, ListSpreadsheetsOutput{}, fmt.Errorf("failed to list spreadsheets: %w", err)
	}

	var resp strings.Builder
	if len(files.Files) == 0 {
		resp.WriteString("No spreadsheets found.\n")
	} else {
		resp.WriteString("Spreadsheets:\n")
		for _, file := range files.Files {
			resp.WriteString(fmt.Sprintf("- %s\n  ID: %s\n  Modified: %s\n\n",
				file.Name, file.Id, file.ModifiedTime))
		}
	}

	return nil, ListSpreadsheetsOutput{Spreadsheets: resp.String()}, nil
}


// GetSpreadsheet handles the get_spreadsheet tool call
func GetSpreadsheet(ctx context.Context, req *mcp.CallToolRequest, input GetSpreadsheetInput) (*mcp.CallToolResult, GetSpreadsheetOutput, error) {
	srv, err := utils.NewSheetsClient(input.Email)
	if err != nil {
		return nil, GetSpreadsheetOutput{}, err
	}

	spreadsheet, err := srv.Spreadsheets.Get(input.SpreadsheetID).Do()
	if err != nil {
		return nil, GetSpreadsheetOutput{}, fmt.Errorf("failed to get spreadsheet: %w", err)
	}

	var resp strings.Builder
	resp.WriteString(fmt.Sprintf("Spreadsheet: %s\n\n", spreadsheet.Properties.Title))
	resp.WriteString("Sheets:\n")
	for _, sheet := range spreadsheet.Sheets {
		props := sheet.Properties
		resp.WriteString(fmt.Sprintf("- %s\n  ID: %d\n  Rows: %d\n  Columns: %d\n\n",
			props.Title, props.SheetId, props.GridProperties.RowCount, props.GridProperties.ColumnCount))
	}

	return nil, GetSpreadsheetOutput{Info: resp.String()}, nil
}

// ReadSheetRange handles the read_sheet_range tool call
func ReadSheetRange(ctx context.Context, req *mcp.CallToolRequest, input ReadSheetRangeInput) (*mcp.CallToolResult, ReadSheetRangeOutput, error) {
	srv, err := utils.NewSheetsClient(input.Email)
	if err != nil {
		return nil, ReadSheetRangeOutput{}, err
	}

	valueRange, err := srv.Spreadsheets.Values.Get(input.SpreadsheetID, input.Range).Do()
	if err != nil {
		return nil, ReadSheetRangeOutput{}, fmt.Errorf("failed to read range: %w", err)
	}

	if len(valueRange.Values) == 0 {
		return nil, ReadSheetRangeOutput{Data: "No data found in the specified range."}, nil
	}

	var resp strings.Builder
	resp.WriteString(fmt.Sprintf("Data from %s:\n\n", input.Range))

	// Format as table
	for i, row := range valueRange.Values {
		resp.WriteString(fmt.Sprintf("Row %d: ", i+1))
		for j, cell := range row {
			if j > 0 {
				resp.WriteString(" | ")
			}
			resp.WriteString(fmt.Sprintf("%v", cell))
		}
		resp.WriteString("\n")
	}

	return nil, ReadSheetRangeOutput{Data: resp.String()}, nil
}


// WriteSheetRange handles the write_sheet_range tool call
func WriteSheetRange(ctx context.Context, req *mcp.CallToolRequest, input WriteSheetRangeInput) (*mcp.CallToolResult, WriteSheetRangeOutput, error) {
	srv, err := utils.NewSheetsClient(input.Email)
	if err != nil {
		return nil, WriteSheetRangeOutput{}, err
	}

	valueRange := &sheets.ValueRange{
		Values: input.Values,
	}

	updateResp, err := srv.Spreadsheets.Values.Update(input.SpreadsheetID, input.Range, valueRange).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return nil, WriteSheetRangeOutput{}, fmt.Errorf("failed to write range: %w", err)
	}

	resp := fmt.Sprintf("Write successful:\n  Updated range: %s\n  Updated cells: %d\n  Updated rows: %d\n  Updated columns: %d",
		updateResp.UpdatedRange, updateResp.UpdatedCells, updateResp.UpdatedRows, updateResp.UpdatedColumns)

	return nil, WriteSheetRangeOutput{Result: resp}, nil
}

// AppendSheetRows handles the append_sheet_rows tool call
func AppendSheetRows(ctx context.Context, req *mcp.CallToolRequest, input AppendSheetRowsInput) (*mcp.CallToolResult, AppendSheetRowsOutput, error) {
	srv, err := utils.NewSheetsClient(input.Email)
	if err != nil {
		return nil, AppendSheetRowsOutput{}, err
	}

	valueRange := &sheets.ValueRange{
		Values: input.Values,
	}

	appendResp, err := srv.Spreadsheets.Values.Append(input.SpreadsheetID, input.Range, valueRange).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		return nil, AppendSheetRowsOutput{}, fmt.Errorf("failed to append rows: %w", err)
	}

	resp := fmt.Sprintf("Append successful:\n  Updated range: %s\n  Appended rows: %d",
		appendResp.Updates.UpdatedRange, appendResp.Updates.UpdatedRows)

	return nil, AppendSheetRowsOutput{Result: resp}, nil
}


// CreateSpreadsheet handles the create_spreadsheet tool call
func CreateSpreadsheet(ctx context.Context, req *mcp.CallToolRequest, input CreateSpreadsheetInput) (*mcp.CallToolResult, CreateSpreadsheetOutput, error) {
	srv, err := utils.NewSheetsClient(input.Email)
	if err != nil {
		return nil, CreateSpreadsheetOutput{}, err
	}

	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: input.Title,
		},
	}

	// Add optional sheet names
	if len(input.SheetNames) > 0 {
		spreadsheet.Sheets = make([]*sheets.Sheet, len(input.SheetNames))
		for i, name := range input.SheetNames {
			spreadsheet.Sheets[i] = &sheets.Sheet{
				Properties: &sheets.SheetProperties{
					Title: name,
				},
			}
		}
	}

	created, err := srv.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		return nil, CreateSpreadsheetOutput{}, fmt.Errorf("failed to create spreadsheet: %w", err)
	}

	resp := fmt.Sprintf("Spreadsheet created successfully:\n  Title: %s\n  ID: %s\n  URL: %s",
		created.Properties.Title, created.SpreadsheetId, created.SpreadsheetUrl)

	return nil, CreateSpreadsheetOutput{Result: resp}, nil
}

// RegisterSheetsTools registers all Sheets-related tools with the MCP server
func RegisterSheetsTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_spreadsheets",
		Description: "List Google Sheets spreadsheets in Drive",
	}, ListSpreadsheets)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_spreadsheet",
		Description: "Get detailed information about a spreadsheet including its sheets",
	}, GetSpreadsheet)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_sheet_range",
		Description: "Read data from a specific range in a spreadsheet",
	}, ReadSheetRange)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "write_sheet_range",
		Description: "Write data to a specific range in a spreadsheet",
	}, WriteSheetRange)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "append_sheet_rows",
		Description: "Append rows of data to a spreadsheet",
	}, AppendSheetRows)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_spreadsheet",
		Description: "Create a new Google Sheets spreadsheet",
	}, CreateSpreadsheet)
}
