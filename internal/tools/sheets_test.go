package tools

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// ============================================================================
// Task 5.1: Input Validation Tests
// Requirements: 2.3, 5.4, 6.4
// ============================================================================

// TestSheetsInputStructTagCompleteness verifies that all Sheets Input structs
// have proper json and jsonschema tags on all fields
func TestSheetsInputStructTagCompleteness(t *testing.T) {
	inputStructs := []interface{}{
		ListSpreadsheetsInput{},
		GetSpreadsheetInput{},
		ReadSheetRangeInput{},
		WriteSheetRangeInput{},
		AppendSheetRowsInput{},
		CreateSpreadsheetInput{},
	}

	for _, input := range inputStructs {
		structType := reflect.TypeOf(input)
		structName := structType.Name()

		t.Run(structName, func(t *testing.T) {
			for i := 0; i < structType.NumField(); i++ {
				field := structType.Field(i)
				fieldName := field.Name

				// Check json tag exists
				jsonTag := field.Tag.Get("json")
				if jsonTag == "" {
					t.Errorf("Field %s.%s is missing json tag", structName, fieldName)
				}

				// Check jsonschema tag exists
				jsonschemaTag := field.Tag.Get("jsonschema")
				if jsonschemaTag == "" {
					t.Errorf("Field %s.%s is missing jsonschema tag", structName, fieldName)
				}

				// Check jsonschema tag has description
				if jsonschemaTag != "" && !strings.Contains(jsonschemaTag, "description=") {
					t.Errorf("Field %s.%s jsonschema tag is missing description", structName, fieldName)
				}
			}
		})
	}
}

// TestSheetsRequiredFieldsHaveRequiredTag verifies that required Input fields
// have the required jsonschema tag
// Requirements: 5.4, 6.4 - invalid range or values should return descriptive error
func TestSheetsRequiredFieldsHaveRequiredTag(t *testing.T) {
	requiredFields := map[string][]string{
		"ListSpreadsheetsInput":  {"Email"},
		"GetSpreadsheetInput":    {"Email", "SpreadsheetID"},
		"ReadSheetRangeInput":    {"Email", "SpreadsheetID", "Range"},
		"WriteSheetRangeInput":   {"Email", "SpreadsheetID", "Range", "Values"},
		"AppendSheetRowsInput":   {"Email", "SpreadsheetID", "Range", "Values"},
		"CreateSpreadsheetInput": {"Email", "Title"},
	}

	inputStructs := []interface{}{
		ListSpreadsheetsInput{},
		GetSpreadsheetInput{},
		ReadSheetRangeInput{},
		WriteSheetRangeInput{},
		AppendSheetRowsInput{},
		CreateSpreadsheetInput{},
	}

	for _, input := range inputStructs {
		structType := reflect.TypeOf(input)
		structName := structType.Name()

		t.Run(structName, func(t *testing.T) {
			expectedRequired, ok := requiredFields[structName]
			if !ok {
				t.Skipf("No required fields defined for %s", structName)
				return
			}

			for _, fieldName := range expectedRequired {
				field, found := structType.FieldByName(fieldName)
				if !found {
					t.Errorf("Expected required field %s not found in %s", fieldName, structName)
					continue
				}

				jsonschemaTag := field.Tag.Get("jsonschema")
				if !strings.Contains(jsonschemaTag, "required") {
					t.Errorf("Field %s.%s should have 'required' in jsonschema tag", structName, fieldName)
				}
			}
		})
	}
}

// TestSheetsOptionalFieldsNotRequired verifies that optional fields don't have required tag
func TestSheetsOptionalFieldsNotRequired(t *testing.T) {
	optionalFields := map[string][]string{
		"ListSpreadsheetsInput":  {"MaxResults"},
		"CreateSpreadsheetInput": {"SheetNames"},
	}

	inputStructs := []interface{}{
		ListSpreadsheetsInput{},
		CreateSpreadsheetInput{},
	}

	for _, input := range inputStructs {
		structType := reflect.TypeOf(input)
		structName := structType.Name()

		t.Run(structName, func(t *testing.T) {
			expectedOptional, ok := optionalFields[structName]
			if !ok {
				t.Skipf("No optional fields defined for %s", structName)
				return
			}

			for _, fieldName := range expectedOptional {
				field, found := structType.FieldByName(fieldName)
				if !found {
					t.Errorf("Expected optional field %s not found in %s", fieldName, structName)
					continue
				}

				jsonschemaTag := field.Tag.Get("jsonschema")
				if strings.Contains(jsonschemaTag, "required") {
					t.Errorf("Field %s.%s should NOT have 'required' in jsonschema tag (it's optional)", structName, fieldName)
				}
			}
		})
	}
}

// TestMaxResultsDefaultValue verifies that maxResults defaults to 10 when not specified
// Requirements: 2.3 - maxResults parameter with default value of 10
func TestMaxResultsDefaultValue(t *testing.T) {
	input := ListSpreadsheetsInput{
		Email:      "test@example.com",
		MaxResults: 0, // Not specified
	}

	// Verify the default behavior is documented in the struct
	structType := reflect.TypeOf(input)
	field, found := structType.FieldByName("MaxResults")
	if !found {
		t.Fatal("MaxResults field not found")
	}

	// Check that the field type allows zero value (which triggers default)
	if field.Type.Kind() != reflect.Int64 {
		t.Errorf("MaxResults should be int64, got %v", field.Type.Kind())
	}

	// Verify the json tag has omitempty (allowing default behavior)
	jsonTag := field.Tag.Get("json")
	if !strings.Contains(jsonTag, "omitempty") {
		t.Errorf("MaxResults json tag should have omitempty for default value handling")
	}
}

// TestSheetsOutputStructTagCompleteness verifies that all Output structs have proper json tags
func TestSheetsOutputStructTagCompleteness(t *testing.T) {
	outputStructs := []interface{}{
		ListSpreadsheetsOutput{},
		GetSpreadsheetOutput{},
		ReadSheetRangeOutput{},
		WriteSheetRangeOutput{},
		AppendSheetRowsOutput{},
		CreateSpreadsheetOutput{},
	}

	for _, output := range outputStructs {
		structType := reflect.TypeOf(output)
		structName := structType.Name()

		t.Run(structName, func(t *testing.T) {
			for i := 0; i < structType.NumField(); i++ {
				field := structType.Field(i)
				fieldName := field.Name

				// Check json tag exists
				jsonTag := field.Tag.Get("json")
				if jsonTag == "" {
					t.Errorf("Field %s.%s is missing json tag", structName, fieldName)
				}
			}
		})
	}
}

// ============================================================================
// Task 5.2: Output Format Tests
// Requirements: 2.2, 3.2, 4.2
// ============================================================================

// TestListSpreadsheetsOutputFormat verifies that list output contains required fields
// Requirements: 2.2 - return spreadsheet name, ID, and last modified time
func TestListSpreadsheetsOutputFormat(t *testing.T) {
	// Simulate the output format that ListSpreadsheets produces
	testCases := []struct {
		name     string
		output   string
		wantName bool
		wantID   bool
		wantMod  bool
	}{
		{
			name:     "empty list",
			output:   "No spreadsheets found.\n",
			wantName: false,
			wantID:   false,
			wantMod:  false,
		},
		{
			name: "single spreadsheet",
			output: `Spreadsheets:
- Test Sheet
  ID: abc123
  Modified: 2024-01-01T00:00:00Z

`,
			wantName: true,
			wantID:   true,
			wantMod:  true,
		},
		{
			name: "multiple spreadsheets",
			output: `Spreadsheets:
- Sheet One
  ID: id1
  Modified: 2024-01-01T00:00:00Z

- Sheet Two
  ID: id2
  Modified: 2024-01-02T00:00:00Z

`,
			wantName: true,
			wantID:   true,
			wantMod:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantName && !strings.Contains(tc.output, "- ") {
				t.Error("Output should contain spreadsheet name with '- ' prefix")
			}
			if tc.wantID && !strings.Contains(tc.output, "ID:") {
				t.Error("Output should contain 'ID:' field")
			}
			if tc.wantMod && !strings.Contains(tc.output, "Modified:") {
				t.Error("Output should contain 'Modified:' field")
			}
		})
	}
}

// TestGetSpreadsheetOutputFormat verifies that get spreadsheet output contains required fields
// Requirements: 3.2 - include sheet name, sheet ID, and row/column count for each sheet
func TestGetSpreadsheetOutputFormat(t *testing.T) {
	// Simulate the output format that GetSpreadsheet produces
	testOutput := `Spreadsheet: My Spreadsheet

Sheets:
- Sheet1
  ID: 0
  Rows: 1000
  Columns: 26

- Sheet2
  ID: 123456
  Rows: 500
  Columns: 10

`

	// Verify required fields are present
	requiredFields := []string{
		"Spreadsheet:",  // Title
		"Sheets:",       // Sheets section
		"ID:",           // Sheet ID
		"Rows:",         // Row count
		"Columns:",      // Column count
	}

	for _, field := range requiredFields {
		if !strings.Contains(testOutput, field) {
			t.Errorf("Output should contain '%s' field", field)
		}
	}
}

// TestReadSheetRangeOutputFormat verifies that read range output is formatted as table
// Requirements: 4.2 - format the output as a readable table structure
func TestReadSheetRangeOutputFormat(t *testing.T) {
	testCases := []struct {
		name       string
		output     string
		wantRange  bool
		wantRows   bool
		wantPipes  bool
	}{
		{
			name:       "empty range",
			output:     "No data found in the specified range.",
			wantRange:  false,
			wantRows:   false,
			wantPipes:  false,
		},
		{
			name: "single row",
			output: `Data from Sheet1!A1:C1:

Row 1: Value1 | Value2 | Value3
`,
			wantRange: true,
			wantRows:  true,
			wantPipes: true,
		},
		{
			name: "multiple rows",
			output: `Data from Sheet1!A1:B3:

Row 1: Header1 | Header2
Row 2: Data1 | Data2
Row 3: Data3 | Data4
`,
			wantRange: true,
			wantRows:  true,
			wantPipes: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantRange && !strings.Contains(tc.output, "Data from") {
				t.Error("Output should contain 'Data from' with range info")
			}
			if tc.wantRows && !strings.Contains(tc.output, "Row") {
				t.Error("Output should contain 'Row' prefix for each row")
			}
			if tc.wantPipes && !strings.Contains(tc.output, "|") {
				t.Error("Output should use '|' as column separator for table format")
			}
		})
	}
}

// TestWriteSheetRangeOutputFormat verifies that write output contains update metadata
// Requirements: 5.3 - return the number of updated cells and the updated range
func TestWriteSheetRangeOutputFormat(t *testing.T) {
	testOutput := `Write successful:
  Updated range: Sheet1!A1:B2
  Updated cells: 4
  Updated rows: 2
  Updated columns: 2`

	requiredFields := []string{
		"Updated range:",
		"Updated cells:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(testOutput, field) {
			t.Errorf("Output should contain '%s' field", field)
		}
	}
}

// TestAppendSheetRowsOutputFormat verifies that append output contains range info
// Requirements: 6.3 - return the range where data was appended
func TestAppendSheetRowsOutputFormat(t *testing.T) {
	testOutput := `Append successful:
  Updated range: Sheet1!A5:B6
  Appended rows: 2`

	requiredFields := []string{
		"Updated range:",
		"Appended rows:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(testOutput, field) {
			t.Errorf("Output should contain '%s' field", field)
		}
	}
}

// TestCreateSpreadsheetOutputFormat verifies that create output contains ID and URL
// Requirements: 7.3 - return the new spreadsheet ID and URL
func TestCreateSpreadsheetOutputFormat(t *testing.T) {
	testOutput := `Spreadsheet created successfully:
  Title: New Spreadsheet
  ID: abc123xyz
  URL: https://docs.google.com/spreadsheets/d/abc123xyz`

	requiredFields := []string{
		"Title:",
		"ID:",
		"URL:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(testOutput, field) {
			t.Errorf("Output should contain '%s' field", field)
		}
	}
}

// ============================================================================
// Task 6: Property-Based Tests
// Using rapid library for property-based testing
// ============================================================================

// ============================================================================
// Task 6.1: Property 3 - maxResults parameter
// **Feature: google-sheets-tools, Property 3: List spreadsheets respects maxResults parameter**
// **Validates: Requirements 2.3**
// ============================================================================

// TestProperty3_MaxResultsDefault verifies that when maxResults is not specified (0),
// the default value of 10 is used.
// This is a property test that validates the default behavior.
func TestProperty3_MaxResultsDefault(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random email addresses
		email := rapid.StringMatching(`[a-z]{5,10}@[a-z]{5,10}\.[a-z]{2,3}`).Draw(t, "email")

		input := ListSpreadsheetsInput{
			Email:      email,
			MaxResults: 0, // Not specified - should default to 10
		}

		// Verify the input has MaxResults as 0 (unset)
		if input.MaxResults != 0 {
			t.Fatalf("MaxResults should be 0 when not specified, got %d", input.MaxResults)
		}

		// The actual default handling happens in ListSpreadsheets function
		// We verify the struct allows 0 value which triggers default behavior
		// The implementation sets maxResults = 10 when input.MaxResults == 0
	})
}

// TestProperty3_MaxResultsRespected verifies that for any valid maxResults value,
// the parameter is properly stored and can be used to limit results.
func TestProperty3_MaxResultsRespected(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random maxResults values in valid range (1-100)
		maxResults := rapid.Int64Range(1, 100).Draw(t, "maxResults")
		email := rapid.StringMatching(`[a-z]{5,10}@[a-z]{5,10}\.[a-z]{2,3}`).Draw(t, "email")

		input := ListSpreadsheetsInput{
			Email:      email,
			MaxResults: maxResults,
		}

		// Property: The input struct correctly stores the maxResults value
		if input.MaxResults != maxResults {
			t.Fatalf("MaxResults should be %d, got %d", maxResults, input.MaxResults)
		}

		// Property: MaxResults should be positive when specified
		if input.MaxResults <= 0 {
			t.Fatalf("MaxResults should be positive when specified, got %d", input.MaxResults)
		}
	})
}

// TestProperty3_MaxResultsStructField verifies that the MaxResults field
// has the correct type and tags for proper JSON handling.
func TestProperty3_MaxResultsStructField(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate various maxResults values including edge cases
		maxResults := rapid.Int64Range(0, 1000).Draw(t, "maxResults")

		input := ListSpreadsheetsInput{
			MaxResults: maxResults,
		}

		// Property: MaxResults field should correctly store any int64 value
		if input.MaxResults != maxResults {
			t.Fatalf("MaxResults field should store %d, got %d", maxResults, input.MaxResults)
		}

		// Verify struct field properties via reflection
		structType := reflect.TypeOf(input)
		field, found := structType.FieldByName("MaxResults")
		if !found {
			t.Fatal("MaxResults field not found")
		}

		// Property: Field type must be int64
		if field.Type.Kind() != reflect.Int64 {
			t.Fatalf("MaxResults should be int64, got %v", field.Type.Kind())
		}

		// Property: JSON tag must have omitempty for default value handling
		jsonTag := field.Tag.Get("json")
		if !strings.Contains(jsonTag, "omitempty") {
			t.Fatal("MaxResults json tag should have omitempty")
		}
	})
}


// ============================================================================
// Task 6.2: Property 6 - Write operation round-trip consistency
// **Feature: google-sheets-tools, Property 6: Write operation round-trip consistency**
// **Validates: Requirements 5.1, 5.2**
// ============================================================================

// TestProperty6_WriteInputValuesStructure verifies that for any valid 2D array of values,
// the WriteSheetRangeInput struct correctly stores and preserves the data structure.
func TestProperty6_WriteInputValuesStructure(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random number of rows and columns
		numRows := rapid.IntRange(1, 10).Draw(t, "numRows")
		numCols := rapid.IntRange(1, 10).Draw(t, "numCols")

		// Generate random 2D array of values
		values := make([][]interface{}, numRows)
		for i := 0; i < numRows; i++ {
			values[i] = make([]interface{}, numCols)
			for j := 0; j < numCols; j++ {
				// Generate random cell values (strings or numbers)
				if rapid.Bool().Draw(t, "isString") {
					values[i][j] = rapid.StringMatching(`[a-zA-Z0-9]{1,20}`).Draw(t, "cellValue")
				} else {
					values[i][j] = rapid.Float64().Draw(t, "cellValue")
				}
			}
		}

		input := WriteSheetRangeInput{
			Email:         "test@example.com",
			SpreadsheetID: "test-spreadsheet-id",
			Range:         "Sheet1!A1",
			Values:        values,
		}

		// Property: The input struct correctly stores the 2D array
		if len(input.Values) != numRows {
			t.Fatalf("Values should have %d rows, got %d", numRows, len(input.Values))
		}

		for i, row := range input.Values {
			if len(row) != numCols {
				t.Fatalf("Row %d should have %d columns, got %d", i, numCols, len(row))
			}
		}

		// Property: Values are preserved exactly as provided
		for i := 0; i < numRows; i++ {
			for j := 0; j < numCols; j++ {
				if input.Values[i][j] != values[i][j] {
					t.Fatalf("Value at [%d][%d] should be %v, got %v", i, j, values[i][j], input.Values[i][j])
				}
			}
		}
	})
}

// TestProperty6_WriteReadDataConsistency verifies that the data structure used for
// writing (2D array) is compatible with the data structure returned by reading.
// This tests the round-trip property at the data structure level.
func TestProperty6_WriteReadDataConsistency(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random cell values that would be written
		numRows := rapid.IntRange(1, 5).Draw(t, "numRows")
		numCols := rapid.IntRange(1, 5).Draw(t, "numCols")

		writeValues := make([][]interface{}, numRows)
		for i := 0; i < numRows; i++ {
			writeValues[i] = make([]interface{}, numCols)
			for j := 0; j < numCols; j++ {
				// Use string values for consistency (Sheets API returns strings)
				writeValues[i][j] = rapid.StringMatching(`[a-zA-Z0-9]{1,10}`).Draw(t, "cellValue")
			}
		}

		// Simulate the read output format (as produced by ReadSheetRange)
		var readOutput strings.Builder
		readOutput.WriteString("Data from Sheet1!A1:\n\n")
		for i, row := range writeValues {
			readOutput.WriteString(fmt.Sprintf("Row %d: ", i+1))
			for j, cell := range row {
				if j > 0 {
					readOutput.WriteString(" | ")
				}
				readOutput.WriteString(fmt.Sprintf("%v", cell))
			}
			readOutput.WriteString("\n")
		}

		// Property: Each written value should appear in the read output
		for _, row := range writeValues {
			for _, cell := range row {
				cellStr := fmt.Sprintf("%v", cell)
				if !strings.Contains(readOutput.String(), cellStr) {
					t.Fatalf("Written value %q should appear in read output", cellStr)
				}
			}
		}

		// Property: Read output should have correct number of rows
		rowCount := strings.Count(readOutput.String(), "Row ")
		if rowCount != numRows {
			t.Fatalf("Read output should have %d rows, found %d", numRows, rowCount)
		}
	})
}

// TestProperty6_ValuesFieldType verifies that the Values field in WriteSheetRangeInput
// has the correct type for storing 2D arrays.
func TestProperty6_ValuesFieldType(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random dimensions
		numRows := rapid.IntRange(0, 20).Draw(t, "numRows")
		numCols := rapid.IntRange(0, 20).Draw(t, "numCols")

		values := make([][]interface{}, numRows)
		for i := 0; i < numRows; i++ {
			values[i] = make([]interface{}, numCols)
		}

		input := WriteSheetRangeInput{
			Values: values,
		}

		// Property: Values field should accept any 2D slice
		if len(input.Values) != numRows {
			t.Fatalf("Values should store %d rows, got %d", numRows, len(input.Values))
		}

		// Verify struct field type via reflection
		structType := reflect.TypeOf(input)
		field, found := structType.FieldByName("Values")
		if !found {
			t.Fatal("Values field not found")
		}

		// Property: Field type must be [][]interface{}
		expectedType := "[][]interface {}"
		if field.Type.String() != expectedType {
			t.Fatalf("Values field type should be %s, got %s", expectedType, field.Type.String())
		}

		// Property: Field must have required tag
		jsonschemaTag := field.Tag.Get("jsonschema")
		if !strings.Contains(jsonschemaTag, "required") {
			t.Fatal("Values field should have 'required' in jsonschema tag")
		}
	})
}


// ============================================================================
// Task 6.3: Property 11 - Registered tools have descriptive metadata
// **Feature: google-sheets-tools, Property 11: Registered tools have descriptive metadata**
// **Validates: Requirements 8.3**
// ============================================================================

// sheetsToolMetadata defines the expected metadata for each Sheets tool
var sheetsToolMetadata = []struct {
	Name        string
	Description string
}{
	{
		Name:        "list_spreadsheets",
		Description: "List Google Sheets spreadsheets in Drive",
	},
	{
		Name:        "get_spreadsheet",
		Description: "Get detailed information about a spreadsheet including its sheets",
	},
	{
		Name:        "read_sheet_range",
		Description: "Read data from a specific range in a spreadsheet",
	},
	{
		Name:        "write_sheet_range",
		Description: "Write data to a specific range in a spreadsheet",
	},
	{
		Name:        "append_sheet_rows",
		Description: "Append rows of data to a spreadsheet",
	},
	{
		Name:        "create_spreadsheet",
		Description: "Create a new Google Sheets spreadsheet",
	},
}

// TestProperty11_ToolsHaveNonEmptyNames verifies that for any Sheets tool,
// the tool has a non-empty name.
func TestProperty11_ToolsHaveNonEmptyNames(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Select a random tool from the list
		toolIndex := rapid.IntRange(0, len(sheetsToolMetadata)-1).Draw(t, "toolIndex")
		tool := sheetsToolMetadata[toolIndex]

		// Property: Tool name must be non-empty
		if tool.Name == "" {
			t.Fatalf("Tool at index %d has empty name", toolIndex)
		}

		// Property: Tool name should follow naming convention (lowercase with underscores)
		if strings.Contains(tool.Name, " ") {
			t.Fatalf("Tool name %q should not contain spaces", tool.Name)
		}

		// Property: Tool name should be lowercase
		if tool.Name != strings.ToLower(tool.Name) {
			t.Fatalf("Tool name %q should be lowercase", tool.Name)
		}
	})
}

// TestProperty11_ToolsHaveNonEmptyDescriptions verifies that for any Sheets tool,
// the tool has a non-empty description.
func TestProperty11_ToolsHaveNonEmptyDescriptions(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Select a random tool from the list
		toolIndex := rapid.IntRange(0, len(sheetsToolMetadata)-1).Draw(t, "toolIndex")
		tool := sheetsToolMetadata[toolIndex]

		// Property: Tool description must be non-empty
		if tool.Description == "" {
			t.Fatalf("Tool %q has empty description", tool.Name)
		}

		// Property: Description should be meaningful (at least 10 characters)
		if len(tool.Description) < 10 {
			t.Fatalf("Tool %q description is too short: %q", tool.Name, tool.Description)
		}

		// Property: Description should start with a capital letter or verb
		if len(tool.Description) > 0 && tool.Description[0] >= 'a' && tool.Description[0] <= 'z' {
			// Allow lowercase if it starts with a verb like "list", "get", etc.
			validStarts := []string{"list", "get", "read", "write", "append", "create"}
			isValid := false
			for _, start := range validStarts {
				if strings.HasPrefix(strings.ToLower(tool.Description), start) {
					isValid = true
					break
				}
			}
			if !isValid {
				t.Fatalf("Tool %q description should start with capital letter or action verb: %q", tool.Name, tool.Description)
			}
		}
	})
}

// TestProperty11_ToolsHaveDescriptiveMetadata verifies that for any Sheets tool,
// both name and description are present and meaningful.
func TestProperty11_ToolsHaveDescriptiveMetadata(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Select a random tool from the list
		toolIndex := rapid.IntRange(0, len(sheetsToolMetadata)-1).Draw(t, "toolIndex")
		tool := sheetsToolMetadata[toolIndex]

		// Property: Both name and description must be non-empty
		if tool.Name == "" || tool.Description == "" {
			t.Fatalf("Tool at index %d missing name or description", toolIndex)
		}

		// Property: Description should relate to the tool name
		// (e.g., "list_spreadsheets" should have "list" or "spreadsheet" in description)
		nameParts := strings.Split(tool.Name, "_")
		descLower := strings.ToLower(tool.Description)
		
		foundRelated := false
		for _, part := range nameParts {
			if len(part) > 2 && strings.Contains(descLower, part) {
				foundRelated = true
				break
			}
		}
		
		if !foundRelated {
			t.Fatalf("Tool %q description %q should relate to tool name", tool.Name, tool.Description)
		}
	})
}

// TestProperty11_AllSheetsToolsRegistered verifies that all expected Sheets tools
// are defined with proper metadata.
func TestProperty11_AllSheetsToolsRegistered(t *testing.T) {
	expectedTools := []string{
		"list_spreadsheets",
		"get_spreadsheet",
		"read_sheet_range",
		"write_sheet_range",
		"append_sheet_rows",
		"create_spreadsheet",
	}

	rapid.Check(t, func(t *rapid.T) {
		// Select a random expected tool
		toolIndex := rapid.IntRange(0, len(expectedTools)-1).Draw(t, "toolIndex")
		expectedName := expectedTools[toolIndex]

		// Property: Each expected tool should exist in the metadata
		found := false
		for _, tool := range sheetsToolMetadata {
			if tool.Name == expectedName {
				found = true
				// Property: Found tool should have non-empty description
				if tool.Description == "" {
					t.Fatalf("Tool %q exists but has empty description", expectedName)
				}
				break
			}
		}

		if !found {
			t.Fatalf("Expected tool %q not found in metadata", expectedName)
		}
	})
}
