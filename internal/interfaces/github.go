package interfaces

type IGitHub interface {
	// GetAuth returns the GitHub authentication configuration.
	GetAuth() IGitHubAuth
	// GetRepos returns the list of repository configurations.
	GetRepos() []IRepoCfg
}
