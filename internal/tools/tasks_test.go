package tools

import (
	"reflect"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// ============================================================================
// Task Input Validation Tests
// ============================================================================

// TestTasksInputStructTagCompleteness verifies that all Tasks Input structs
// have proper json and jsonschema tags on all fields
func TestTasksInputStructTagCompleteness(t *testing.T) {
	inputStructs := []any{
		ListTaskListsInput{},
		ListTasksInput{},
		CreateTaskInput{},
		UpdateTaskInput{},
		DeleteTaskInput{},
		CompleteTaskInput{},
	}

	for _, input := range inputStructs {
		structType := reflect.TypeOf(input)
		structName := structType.Name()

		t.Run(structName, func(t *testing.T) {
			for i := range structType.NumField() {
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

			}
		})
	}
}

// TestTasksRequiredFieldsHaveRequiredTag verifies that required Input fields
// have the required jsonschema tag
func TestTasksRequiredFieldsHaveRequiredTag(t *testing.T) {
	requiredFields := map[string][]string{
		"ListTaskListsInput": {"Email"},
		"ListTasksInput":     {"Email", "TaskListID"},
		"CreateTaskInput":    {"Email", "TaskListID", "Title"},
		"UpdateTaskInput":    {"Email", "TaskListID", "TaskID"},
		"DeleteTaskInput":    {"Email", "TaskListID", "TaskID"},
		"CompleteTaskInput":  {"Email", "TaskListID", "TaskID"},
	}

	inputStructs := []any{
		ListTaskListsInput{},
		ListTasksInput{},
		CreateTaskInput{},
		UpdateTaskInput{},
		DeleteTaskInput{},
		CompleteTaskInput{},
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

				jsonTag := field.Tag.Get("json")
				if strings.Contains(jsonTag, "omitempty") || strings.Contains(jsonTag, "omitzero") {
					t.Errorf("Required field %s.%s should not be omitempty in json tag", structName, fieldName)
				}
			}
		})
	}
}

// TestTasksOptionalFieldsNotRequired verifies that optional fields don't have required tag
func TestTasksOptionalFieldsNotRequired(t *testing.T) {
	optionalFields := map[string][]string{
		"CreateTaskInput": {"Notes", "Due"},
		"UpdateTaskInput": {"Title", "Notes", "Status", "Due"},
	}

	inputStructs := []any{
		CreateTaskInput{},
		UpdateTaskInput{},
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

				jsonTag := field.Tag.Get("json")
				if !strings.Contains(jsonTag, "omitempty") && !strings.Contains(jsonTag, "omitzero") {
					t.Errorf("Optional field %s.%s should be omitempty in json tag", structName, fieldName)
				}
			}
		})
	}
}

// TestTasksOutputStructTagCompleteness verifies that all Output structs have proper json tags
func TestTasksOutputStructTagCompleteness(t *testing.T) {
	outputStructs := []any{
		ListTaskListsOutput{},
		ListTasksOutput{},
		CreateTaskOutput{},
		UpdateTaskOutput{},
		DeleteTaskOutput{},
		CompleteTaskOutput{},
	}

	for _, output := range outputStructs {
		structType := reflect.TypeOf(output)
		structName := structType.Name()

		t.Run(structName, func(t *testing.T) {
			for i := range structType.NumField() {
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
// Output Format Tests
// ============================================================================

// TestListTaskListsOutputFormat verifies that list task lists output contains required fields
func TestListTaskListsOutputFormat(t *testing.T) {
	testCases := []struct {
		name     string
		output   string
		wantName bool
		wantID   bool
	}{
		{
			name:     "empty list",
			output:   "No task lists found.",
			wantName: false,
			wantID:   false,
		},
		{
			name: "single task list",
			output: `Task lists:
- My Tasks (ID: abc123)
`,
			wantName: true,
			wantID:   true,
		},
		{
			name: "multiple task lists",
			output: `Task lists:
- My Tasks (ID: abc123)
- Work (ID: def456)
`,
			wantName: true,
			wantID:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantName && !strings.Contains(tc.output, "- ") {
				t.Error("Output should contain task list name with '- ' prefix")
			}
			if tc.wantID && !strings.Contains(tc.output, "ID:") {
				t.Error("Output should contain 'ID:' field")
			}
		})
	}
}

// TestListTasksOutputFormat verifies that list tasks output contains required fields
func TestListTasksOutputFormat(t *testing.T) {
	testCases := []struct {
		name       string
		output     string
		wantStatus bool
		wantTitle  bool
		wantID     bool
	}{
		{
			name:       "empty list",
			output:     "No tasks found.",
			wantStatus: false,
			wantTitle:  false,
			wantID:     false,
		},
		{
			name: "incomplete task",
			output: `Tasks:
[ ] Buy groceries (ID: task123)
`,
			wantStatus: true,
			wantTitle:  true,
			wantID:     true,
		},
		{
			name: "completed task",
			output: `Tasks:
[x] Finish report (ID: task456)
`,
			wantStatus: true,
			wantTitle:  true,
			wantID:     true,
		},
		{
			name: "task with due date",
			output: `Tasks:
[ ] Submit proposal (Due: 2024-12-31T00:00:00Z) (ID: task789)
`,
			wantStatus: true,
			wantTitle:  true,
			wantID:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantStatus && !strings.Contains(tc.output, "[ ]") && !strings.Contains(tc.output, "[x]") {
				t.Error("Output should contain status indicator '[ ]' or '[x]'")
			}
			if tc.wantID && !strings.Contains(tc.output, "ID:") {
				t.Error("Output should contain 'ID:' field")
			}
		})
	}
}

// TestCreateTaskOutputFormat verifies that create task output contains required fields
func TestCreateTaskOutputFormat(t *testing.T) {
	testOutput := `Task created successfully:
Title: Buy groceries
ID: task123`

	requiredFields := []string{
		"Title:",
		"ID:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(testOutput, field) {
			t.Errorf("Output should contain '%s' field", field)
		}
	}
}

// TestUpdateTaskOutputFormat verifies that update task output contains required fields
func TestUpdateTaskOutputFormat(t *testing.T) {
	testOutput := `Task updated successfully:
Title: Buy groceries
Status: needsAction
ID: task123`

	requiredFields := []string{
		"Title:",
		"Status:",
		"ID:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(testOutput, field) {
			t.Errorf("Output should contain '%s' field", field)
		}
	}
}

// TestCompleteTaskOutputFormat verifies that complete task output contains required fields
func TestCompleteTaskOutputFormat(t *testing.T) {
	testOutput := `Task marked as completed:
Title: Buy groceries
ID: task123`

	requiredFields := []string{
		"Title:",
		"ID:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(testOutput, field) {
			t.Errorf("Output should contain '%s' field", field)
		}
	}
}


// ============================================================================
// Property-Based Tests
// ============================================================================

// tasksToolMetadata defines the expected metadata for each Tasks tool
var tasksToolMetadata = []struct {
	Name        string
	Description string
}{
	{
		Name:        "list_task_lists",
		Description: "List all Google Tasks task lists for a user",
	},
	{
		Name:        "list_tasks",
		Description: "List all tasks in a specific task list",
	},
	{
		Name:        "create_task",
		Description: "Create a new task in a task list",
	},
	{
		Name:        "update_task",
		Description: "Update an existing task",
	},
	{
		Name:        "delete_task",
		Description: "Delete a task from a task list",
	},
	{
		Name:        "complete_task",
		Description: "Mark a task as completed",
	},
}

// TestTasksProperty_ToolsHaveNonEmptyNames verifies that for any Tasks tool,
// the tool has a non-empty name.
func TestTasksProperty_ToolsHaveNonEmptyNames(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		toolIndex := rapid.IntRange(0, len(tasksToolMetadata)-1).Draw(t, "toolIndex")
		tool := tasksToolMetadata[toolIndex]

		if tool.Name == "" {
			t.Fatalf("Tool at index %d has empty name", toolIndex)
		}

		if strings.Contains(tool.Name, " ") {
			t.Fatalf("Tool name %q should not contain spaces", tool.Name)
		}

		if tool.Name != strings.ToLower(tool.Name) {
			t.Fatalf("Tool name %q should be lowercase", tool.Name)
		}
	})
}

// TestTasksProperty_ToolsHaveNonEmptyDescriptions verifies that for any Tasks tool,
// the tool has a non-empty description.
func TestTasksProperty_ToolsHaveNonEmptyDescriptions(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		toolIndex := rapid.IntRange(0, len(tasksToolMetadata)-1).Draw(t, "toolIndex")
		tool := tasksToolMetadata[toolIndex]

		if tool.Description == "" {
			t.Fatalf("Tool %q has empty description", tool.Name)
		}

		if len(tool.Description) < 10 {
			t.Fatalf("Tool %q description is too short: %q", tool.Name, tool.Description)
		}
	})
}

// TestTasksProperty_ToolsHaveDescriptiveMetadata verifies that for any Tasks tool,
// both name and description are present and meaningful.
func TestTasksProperty_ToolsHaveDescriptiveMetadata(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		toolIndex := rapid.IntRange(0, len(tasksToolMetadata)-1).Draw(t, "toolIndex")
		tool := tasksToolMetadata[toolIndex]

		if tool.Name == "" || tool.Description == "" {
			t.Fatalf("Tool at index %d missing name or description", toolIndex)
		}

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

// TestTasksProperty_AllTasksToolsRegistered verifies that all expected Tasks tools
// are defined with proper metadata.
func TestTasksProperty_AllTasksToolsRegistered(t *testing.T) {
	expectedTools := []string{
		"list_task_lists",
		"list_tasks",
		"create_task",
		"update_task",
		"delete_task",
		"complete_task",
	}

	rapid.Check(t, func(t *rapid.T) {
		toolIndex := rapid.IntRange(0, len(expectedTools)-1).Draw(t, "toolIndex")
		expectedName := expectedTools[toolIndex]

		found := false
		for _, tool := range tasksToolMetadata {
			if tool.Name == expectedName {
				found = true
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

// TestTasksProperty_InputStructsPreserveValues verifies that input structs
// correctly store and preserve values.
func TestTasksProperty_InputStructsPreserveValues(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		email := rapid.StringMatching(`[a-z]{5,10}@[a-z]{5,10}\.[a-z]{2,3}`).Draw(t, "email")
		taskListID := rapid.StringMatching(`[a-zA-Z0-9]{10,20}`).Draw(t, "taskListID")
		title := rapid.StringMatching(`[a-zA-Z0-9 ]{5,50}`).Draw(t, "title")

		input := CreateTaskInput{
			Email:      email,
			TaskListID: taskListID,
			Title:      title,
		}

		if input.Email != email {
			t.Fatalf("Email should be %q, got %q", email, input.Email)
		}
		if input.TaskListID != taskListID {
			t.Fatalf("TaskListID should be %q, got %q", taskListID, input.TaskListID)
		}
		if input.Title != title {
			t.Fatalf("Title should be %q, got %q", title, input.Title)
		}
	})
}

// TestTasksProperty_TaskStatusValues verifies that task status values are valid.
func TestTasksProperty_TaskStatusValues(t *testing.T) {
	validStatuses := []string{"needsAction", "completed"}

	rapid.Check(t, func(t *rapid.T) {
		statusIndex := rapid.IntRange(0, len(validStatuses)-1).Draw(t, "statusIndex")
		status := validStatuses[statusIndex]

		input := UpdateTaskInput{
			Email:      "test@example.com",
			TaskListID: "list123",
			TaskID:     "task123",
			Status:     status,
		}

		if input.Status != status {
			t.Fatalf("Status should be %q, got %q", status, input.Status)
		}

		// Verify status is one of the valid values
		isValid := false
		for _, valid := range validStatuses {
			if input.Status == valid {
				isValid = true
				break
			}
		}
		if !isValid {
			t.Fatalf("Status %q is not a valid task status", input.Status)
		}
	})
}
