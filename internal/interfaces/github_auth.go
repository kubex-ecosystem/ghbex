package interfaces

type IGitHubAuth interface {
	GetKind() string
	GetToken() string
	GetAppID() int64
	GetInstallationID() int64
	GetPrivateKeyPath() string
	GetBaseURL() string
	GetUploadURL() string
	SetKind(kind string)
	SetToken(token string)
	SetAppID(appID int64)
	SetInstallationID(installationID int64)
	SetPrivateKeyPath(privateKeyPath string)
	SetBaseURL(baseURL string)
	SetUploadURL(uploadURL string)
}
