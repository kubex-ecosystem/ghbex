package gromptz

import gromptProviders "github.com/rafa-mori/grompt/factory/providers"

type Provider = gromptProviders.Provider

func NewProviders(
	bindAddr string,
	port string,
	claudeKey string,
	openaiKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	geminiKey string,
	chatgptKey string,
	cfg GromptConfig) []Provider {
	providers := gromptProviders.Initialize(
		bindAddr,
		port,
		openaiKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatgptKey,
		nil,
	)
	providers = append(providers, gromptProviders.NewProvider("gemini", geminiKey, "v1beta", cfg))
	providers = append(providers, gromptProviders.NewProvider("chatgpt", chatgptKey, "v1", cfg))

	return providers
}
func NewProvider(
	name string,
	apiKey string,
	version string,
	cfg GromptConfig,
) Provider {
	return gromptProviders.NewProvider(name, apiKey, version, cfg)
}
