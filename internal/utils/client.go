package utils

import (
	"context"
	"fmt"
	"os"

	"butterfly.orx.me/core/log"
	"golang.org/x/oauth2"
	goauth "golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func defaultServiceAccount() ([]byte, error) {
	path := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	if path == "" {
		return nil, fmt.Errorf("GOOGLE_SERVICE_ACCOUNT environment variable not set")
	}

	// Read service account JSON from file
	sa, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}
	return sa, nil
}

func DefaultClient() (*admin.Service, error) {
	sa, err := defaultServiceAccount()
	if err != nil {
		return nil, err
	}

	// Get admin email from environment variable
	adminEmail := os.Getenv("GOOGLE_ADMIN_EMAIL")
	if adminEmail == "" {
		return nil, fmt.Errorf("GOOGLE_ADMIN_EMAIL environment variable not set")
	}

	// Create client using service account and admin email
	return newClient(sa, adminEmail)
}

func NewGmailClient(email string) (*gmail.Service, error) {
	sa, err := defaultServiceAccount()
	if err != nil {
		return nil, err
	}
	// Create client using service account and admin email
	return newGmailClient(sa, email)
}

func newGmailClient(sa []byte, email string) (*gmail.Service, error) {
	ts, err := tokenSource(sa, email, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func tokenSource(sa []byte, adminEmail string, scopes ...string) (oauth2.TokenSource, error) {
	ctx := context.Background()
	logger := log.FromContext(ctx)

	cfg, err := goauth.JWTConfigFromJSON(sa, scopes...)
	if err != nil {
		logger.Error("failed to parse service account JSON", "error", err)
		return nil, err
	}
	cfg.Subject = adminEmail

	return cfg.TokenSource(ctx), nil
}

func newClient(sa []byte, adminEmail string) (*admin.Service, error) {
	ts, err := tokenSource(sa, adminEmail, admin.AdminDirectoryUserScope)
	if err != nil {
		return nil, err
	}

	srv, err := admin.NewService(context.Background(), option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func NewCalendarClient(email string) (*calendar.Service, error) {
	sa, err := defaultServiceAccount()
	if err != nil {
		return nil, err
	}
	// Create client using service account and email
	return newCalendarClient(sa, email)
}

func newCalendarClient(sa []byte, email string) (*calendar.Service, error) {
	ts, err := tokenSource(sa, email, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(context.Background(), option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}
	return srv, nil
}
