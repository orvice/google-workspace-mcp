package tools

import (
	"reflect"
	"strings"
	"testing"
)

// **Feature: go-sdk-refactor, Property 1: Input Struct Tag Completeness**
// **Validates: Requirements 7.1, 7.2, 7.3**
//
// For any Input struct used in tool handlers, the struct SHALL have both
// `json` and `jsonschema` tags on all fields, ensuring automatic schema
// generation works correctly.

// allInputStructs returns all Input structs that need to be validated
func allInputStructs() []interface{} {
	return []interface{}{
		ListUsersInput{},
		CreateUserInput{},
		ListGmailInput{},
		ListCalendarEventsInput{},
		CreateCalendarEventInput{},
	}
}

// allOutputStructs returns all Output structs that need to be validated
func allOutputStructs() []interface{} {
	return []interface{}{
		ListUsersOutput{},
		CreateUserOutput{},
		ListGmailOutput{},
		ListCalendarEventsOutput{},
		CreateCalendarEventOutput{},
	}
}

// TestInputStructTagCompleteness verifies that all Input structs have proper json and jsonschema tags
func TestInputStructTagCompleteness(t *testing.T) {
	for _, input := range allInputStructs() {
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

				// Check jsonschema tag exists. With jsonschema-go, the entire
				// tag value is the field description, so a non-empty tag is all
				// that's required.
				jsonschemaTag := field.Tag.Get("jsonschema")
				if jsonschemaTag == "" {
					t.Errorf("Field %s.%s is missing jsonschema tag", structName, fieldName)
				}
			}
		})
	}
}

// TestOutputStructTagCompleteness verifies that all Output structs have proper json tags
func TestOutputStructTagCompleteness(t *testing.T) {
	for _, output := range allOutputStructs() {
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

// TestRequiredFieldsHaveRequiredTag verifies that required Input fields are not
// marked omitempty in their json tag. With jsonschema-go, a field is required
// unless its json tag carries omitempty/omitzero.
func TestRequiredFieldsHaveRequiredTag(t *testing.T) {
	// Map of struct name to required field names based on requirements
	requiredFields := map[string][]string{
		"ListUsersInput":            {"Domain"},
		"CreateUserInput":           {"Email", "FirstName", "LastName", "Password"},
		"ListGmailInput":            {"Email"},
		"ListCalendarEventsInput":   {"Email"},
		"CreateCalendarEventInput":  {"Email", "Summary", "StartTime", "EndTime"},
	}

	for _, input := range allInputStructs() {
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

// TestOptionalFieldsNotRequired verifies that optional fields carry omitempty in
// their json tag, which is what marks them optional for jsonschema-go.
func TestOptionalFieldsNotRequired(t *testing.T) {
	// Map of struct name to optional field names
	optionalFields := map[string][]string{
		"CreateCalendarEventInput": {"Description"},
	}

	for _, input := range allInputStructs() {
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
