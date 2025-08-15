// Package githubx provides a GitHub client with support for both PAT and App authentication.
package githubx

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

// -------- PAT --------

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

// installationTokenSource gera e renova tokens de instalação sob demanda.
type installationTokenSource struct {
	appID          int64
	installationID int64
	privateKey     *rsa.PrivateKey
	apiBase        string
	client         *http.Client

	mu    sync.Mutex
	token *oauth2.Token
}

func (s *installationTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// ainda válido?
	if s.token != nil && s.token.Expiry.After(time.Now().Add(30*time.Second)) {
		return s.token, nil
	}

	// gera JWT do App
	jwt, err := makeAppJWT(s.privateKey, s.appID, 9*time.Minute)
	if err != nil {
		return nil, err
	}

	// troca por installation token
	tok, exp, err := createInstallationToken(s.client, s.apiBase, s.installationID, jwt)
	if err != nil {
		return nil, err
	}
	s.token = &oauth2.Token{
		AccessToken: tok,
		TokenType:   "token",
		Expiry:      exp,
	}
	return s.token, nil

}

func parseRSAPrivateKeyFromPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := k.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not RSA key")
		}
		return rsaKey, nil
	default:
		return nil, errors.New("unsupported PEM type: " + block.Type)
	}
}

type Report struct {
	Owner  string    `json:"owner"`
	Repo   string    `json:"repo"`
	When   time.Time `json:"when"`
	DryRun bool      `json:"dry_run"`

	Runs struct {
		Deleted int     `json:"deleted"`
		Kept    int     `json:"kept"`
		IDs     []int64 `json:"ids"`
	} `json:"runs"`

	Artifacts struct {
		Deleted int     `json:"deleted"`
		IDs     []int64 `json:"ids"`
	} `json:"artifacts"`

	Releases struct {
		DeletedDrafts int      `json:"deleted_drafts"`
		Tags          []string `json:"tags"`
	} `json:"releases"`

	Notes []string `json:"notes"`
}

func ToMarkdown(r *Report) string {
	return fmt.Sprintf(`# Sanitize %s/%s
- when: %s
- dry_run: %v

## runs
- deleted: %d
- kept(success last): %d

## artifacts
- deleted: %d

## releases
- deleted drafts: %d
- tags: %v

notes:
%v
`,
		r.Owner, r.Repo, r.When.Format(time.RFC3339), r.DryRun,
		r.Runs.Deleted, r.Runs.Kept,
		r.Artifacts.Deleted,
		r.Releases.DeletedDrafts, r.Releases.Tags,
		r.Notes,
	)
}
