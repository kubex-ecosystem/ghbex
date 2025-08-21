package defs

import "github.com/rafa-mori/ghbex/internal/interfaces"

type GitHubAuth struct {
	Kind           string `yaml:"kind" json:"kind"` // pat|app
	Token          string `yaml:"token" json:"token"`
	AppID          int64  `yaml:"app_id" json:"app_id"`
	InstallationID int64  `yaml:"installation_id" json:"installation_id"`
	PrivateKeyPath string `yaml:"private_key_path" json:"private_key_path"`
	BaseURL        string `yaml:"base_url" json:"base_url"`
	UploadURL      string `yaml:"upload_url" json:"upload_url"`
}

func NewGitHubAuthType(kind, token string, appID, installationID int64, privateKeyPath, baseURL, uploadURL string) *GitHubAuth {
	return &GitHubAuth{
		Kind:           kind,
		Token:          token,
		AppID:          appID,
		InstallationID: installationID,
		PrivateKeyPath: privateKeyPath,
		BaseURL:        baseURL,
		UploadURL:      uploadURL,
	}
}

func NewGitHubAuth(kind, token string, appID, installationID int64, privateKeyPath, baseURL, uploadURL string) interfaces.IGitHubAuth {
	return NewGitHubAuthType(kind, token, appID, installationID, privateKeyPath, baseURL, uploadURL)
}

func (a *GitHubAuth) GetKind() string {
	return a.Kind
}

func (a *GitHubAuth) SetKind(kind string) {
	a.Kind = kind
}

func (a *GitHubAuth) GetToken() string {
	return a.Token
}

func (a *GitHubAuth) SetToken(token string) {
	a.Token = token
}

func (a *GitHubAuth) GetAppID() int64 {
	return a.AppID
}

func (a *GitHubAuth) SetAppID(appID int64) {
	a.AppID = appID
}

func (a *GitHubAuth) GetInstallationID() int64 {
	return a.InstallationID
}

func (a *GitHubAuth) SetInstallationID(installationID int64) {
	a.InstallationID = installationID
}

func (a *GitHubAuth) GetPrivateKeyPath() string {
	return a.PrivateKeyPath
}

func (a *GitHubAuth) SetPrivateKeyPath(privateKeyPath string) {
	a.PrivateKeyPath = privateKeyPath
}

func (a *GitHubAuth) GetBaseURL() string {
	return a.BaseURL
}

func (a *GitHubAuth) SetBaseURL(baseURL string) {
	a.BaseURL = baseURL
}

func (a *GitHubAuth) GetUploadURL() string {
	return a.UploadURL
}

func (a *GitHubAuth) SetUploadURL(uploadURL string) {
	a.UploadURL = uploadURL
}
