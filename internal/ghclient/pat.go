package ghclient

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

type PATConfig struct {
	Token     string
	BaseURL   string
	UploadURL string
}

func NewPAT(ctx context.Context, cfg PATConfig) (*github.Client, error) {
	if strings.TrimSpace(cfg.Token) == "" {
		return nil, errors.New("missing PAT token")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	hc := oauth2.NewClient(ctx, ts)
	if cfg.BaseURL != "" {
		upload := cfg.UploadURL
		if upload == "" {
			upload = strings.Replace(cfg.BaseURL, "/api/v3/", "/api/uploads/", 1)
		}
		return github.NewEnterpriseClient(cfg.BaseURL, upload, hc)
	}
	return github.NewClient(hc), nil
}
