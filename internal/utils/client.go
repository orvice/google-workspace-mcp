package utils

import (
	"context"
	"fmt"
	"os"

	"butterfly.orx.me/core/log"
	goauth "golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func DefaultClient() (*admin.Service, error) {
	path := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	if path == "" {
		return nil, fmt.Errorf("GOOGLE_SERVICE_ACCOUNT environment variable not set")
	}

	// Read service account JSON from file
	sa, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}

	// Get admin email from environment variable
	adminEmail := os.Getenv("GOOGLE_ADMIN_EMAIL")
	if adminEmail == "" {
		return nil, fmt.Errorf("GOOGLE_ADMIN_EMAIL environment variable not set")
	}

	// Create client using service account and admin email
	return newClient(sa, adminEmail)
}

func newClient(sa []byte, adminEmail string) (*admin.Service, error) {
	ctx := context.Background()
	logger := log.FromContext(ctx)

	scopes := []string{
		admin.AdminDirectoryUserScope,
	}

	cfg, err := goauth.JWTConfigFromJSON(sa, scopes...)
	if err != nil {
		logger.Error("failed to parse service account JSON", "error", err)
		return nil, err
	}
	cfg.Subject = adminEmail

	ts := cfg.TokenSource(ctx)

	srv, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		logger.Error("failed to create admin service", "error", err)
		return nil, err
	}
	return srv, nil
}
