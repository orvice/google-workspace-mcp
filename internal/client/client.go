package client

import (
	"context"

	"butterfly.orx.me/core/log"
	goauth "golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func NewClient(sa []byte, adminEmail string) (*admin.Service, error) {
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
