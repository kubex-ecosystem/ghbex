// Package ghclient provides a GitHub client with support for both PAT and App authentication.
package ghclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

// -------- GitHub App (sem ghinstallation) --------

type AppConfig struct {
	AppID          int64
	InstallationID int64
	PrivateKeyPath string
	BaseURL        string
	UploadURL      string
}

func NewApp(ctx context.Context, cfg AppConfig) (*github.Client, error) {
	if cfg.AppID == 0 || cfg.InstallationID == 0 || cfg.PrivateKeyPath == "" {
		return nil, errors.New("missing GitHub App config")
	}
	keyPEM, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	priv, err := parseRSAPrivateKeyFromPEM(keyPEM)
	if err != nil {
		return nil, fmt.Errorf("invalid app private key: %w", err)
	}

	baseAPI := cfg.BaseURL
	if baseAPI == "" {
		baseAPI = "https://api.github.com/"
	}
	tokenSrc := &installationTokenSource{
		appID:          cfg.AppID,
		installationID: cfg.InstallationID,
		privateKey:     priv,
		apiBase:        strings.TrimRight(baseAPI, "/"),
		client:         http.DefaultClient,
	}

	hc := oauth2.NewClient(ctx, tokenSrc)
	if cfg.BaseURL != "" {
		upload := cfg.UploadURL
		if upload == "" {
			upload = strings.Replace(cfg.BaseURL, "/api/v3/", "/api/uploads/", 1)
		}
		return github.NewEnterpriseClient(cfg.BaseURL, upload, hc)
	}
	return github.NewClient(hc), nil

}
