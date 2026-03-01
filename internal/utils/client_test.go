package utils

import (
	"os"
	"strings"
	"testing"
)

// **Feature: go-sdk-refactor, Property 3: Environment Variable Error Messages**
// **Validates: Requirements 8.3, 8.4, 8.5**
//
// For any client creation function in the utils package, when a required
// environment variable is not set, the function SHALL return an error
// containing the name of the missing variable.

// TestDefaultClientMissingServiceAccount verifies error message when GOOGLE_SERVICE_ACCOUNT is not set
func TestDefaultClientMissingServiceAccount(t *testing.T) {
	// Save and clear environment variables
	origServiceAccount := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	origAdminEmail := os.Getenv("GOOGLE_ADMIN_EMAIL")
	os.Unsetenv("GOOGLE_SERVICE_ACCOUNT")
	os.Unsetenv("GOOGLE_ADMIN_EMAIL")
	defer func() {
		if origServiceAccount != "" {
			os.Setenv("GOOGLE_SERVICE_ACCOUNT", origServiceAccount)
		}
		if origAdminEmail != "" {
			os.Setenv("GOOGLE_ADMIN_EMAIL", origAdminEmail)
		}
	}()

	_, err := DefaultClient()
	if err == nil {
		t.Fatal("Expected error when GOOGLE_SERVICE_ACCOUNT is not set, got nil")
	}

	if !strings.Contains(err.Error(), "GOOGLE_SERVICE_ACCOUNT") {
		t.Errorf("Error message should contain 'GOOGLE_SERVICE_ACCOUNT', got: %s", err.Error())
	}
}

// TestDefaultClientMissingAdminEmail verifies error message when GOOGLE_ADMIN_EMAIL is not set
func TestDefaultClientMissingAdminEmail(t *testing.T) {
	// Save and clear environment variables
	origServiceAccount := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	origAdminEmail := os.Getenv("GOOGLE_ADMIN_EMAIL")
	
	// Create a temporary file to simulate service account
	tmpFile, err := os.CreateTemp("", "test-sa-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	// Write minimal valid JSON (will fail later but we want to test env var check)
	tmpFile.WriteString(`{"type": "service_account"}`)
	tmpFile.Close()
	
	os.Setenv("GOOGLE_SERVICE_ACCOUNT", tmpFile.Name())
	os.Unsetenv("GOOGLE_ADMIN_EMAIL")
	defer func() {
		if origServiceAccount != "" {
			os.Setenv("GOOGLE_SERVICE_ACCOUNT", origServiceAccount)
		} else {
			os.Unsetenv("GOOGLE_SERVICE_ACCOUNT")
		}
		if origAdminEmail != "" {
			os.Setenv("GOOGLE_ADMIN_EMAIL", origAdminEmail)
		}
	}()

	_, err = DefaultClient()
	if err == nil {
		t.Fatal("Expected error when GOOGLE_ADMIN_EMAIL is not set, got nil")
	}

	if !strings.Contains(err.Error(), "GOOGLE_ADMIN_EMAIL") {
		t.Errorf("Error message should contain 'GOOGLE_ADMIN_EMAIL', got: %s", err.Error())
	}
}

// TestNewGmailClientMissingServiceAccount verifies error message when GOOGLE_SERVICE_ACCOUNT is not set
func TestNewGmailClientMissingServiceAccount(t *testing.T) {
	// Save and clear environment variable
	origServiceAccount := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	os.Unsetenv("GOOGLE_SERVICE_ACCOUNT")
	defer func() {
		if origServiceAccount != "" {
			os.Setenv("GOOGLE_SERVICE_ACCOUNT", origServiceAccount)
		}
	}()

	_, err := NewGmailClient("test@example.com")
	if err == nil {
		t.Fatal("Expected error when GOOGLE_SERVICE_ACCOUNT is not set, got nil")
	}

	if !strings.Contains(err.Error(), "GOOGLE_SERVICE_ACCOUNT") {
		t.Errorf("Error message should contain 'GOOGLE_SERVICE_ACCOUNT', got: %s", err.Error())
	}
}

// TestNewCalendarClientMissingServiceAccount verifies error message when GOOGLE_SERVICE_ACCOUNT is not set
func TestNewCalendarClientMissingServiceAccount(t *testing.T) {
	// Save and clear environment variable
	origServiceAccount := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	os.Unsetenv("GOOGLE_SERVICE_ACCOUNT")
	defer func() {
		if origServiceAccount != "" {
			os.Setenv("GOOGLE_SERVICE_ACCOUNT", origServiceAccount)
		}
	}()

	_, err := NewCalendarClient("test@example.com")
	if err == nil {
		t.Fatal("Expected error when GOOGLE_SERVICE_ACCOUNT is not set, got nil")
	}

	if !strings.Contains(err.Error(), "GOOGLE_SERVICE_ACCOUNT") {
		t.Errorf("Error message should contain 'GOOGLE_SERVICE_ACCOUNT', got: %s", err.Error())
	}
}

// TestServiceAccountFileNotFound verifies error when service account file doesn't exist
func TestServiceAccountFileNotFound(t *testing.T) {
	// Save and clear environment variable
	origServiceAccount := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	os.Setenv("GOOGLE_SERVICE_ACCOUNT", "/nonexistent/path/to/service-account.json")
	defer func() {
		if origServiceAccount != "" {
			os.Setenv("GOOGLE_SERVICE_ACCOUNT", origServiceAccount)
		} else {
			os.Unsetenv("GOOGLE_SERVICE_ACCOUNT")
		}
	}()

	_, err := DefaultClient()
	if err == nil {
		t.Fatal("Expected error when service account file doesn't exist, got nil")
	}

	// Error should indicate file reading failure
	if !strings.Contains(err.Error(), "failed to read") {
		t.Errorf("Error message should indicate file reading failure, got: %s", err.Error())
	}
}
