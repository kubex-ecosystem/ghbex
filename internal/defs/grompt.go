package defs

import (
	"github.com/rafa-mori/grompt"
	gromptConfig "github.com/rafa-mori/grompt/factory/config"
	gromptProviders "github.com/rafa-mori/grompt/factory/providers"
)

type IGrompt = grompt.Grompt

type Grompt = grompt.PromptEngine
type APIConfig = grompt.APIConfig
type Capabilities = gromptProviders.Capabilities
type Provider = gromptProviders.Provider
type GromptConfig = grompt.Config
type PromptEngine = grompt.PromptEngine

type GromptResult struct {
	Score      float64 `json:"score"`
	Assessment string  `json:"assessment"`
}

func NewGromptEngine(cfg GromptConfig) Grompt {
	gmpt := grompt.NewPromptEngine(cfg)
	return gmpt
}
func NewGromptConfigFromFile(filePath string) GromptConfig {
	return gromptConfig.NewConfigFromFile(filePath)
}
func NewGromptConfig(
	port string,
	openAIKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	claudeKey string,
	geminiKey string,
) GromptConfig {
	return gromptConfig.NewConfig(port, openAIKey, deepSeekKey, ollamaEndpoint, claudeKey, geminiKey)
}
func NewPromptEngine(cfg GromptConfig) PromptEngine { return grompt.NewPromptEngine(cfg) }
func NewAPIConfig(configFilePath, provider string) APIConfig {
	gmptCfg := grompt.DefaultConfig(configFilePath)
	return gmptCfg.GetAPIConfig(provider)
}

func NewProviders(claudeKey string, openaiKey string, deepseekKey string, ollamaEndpoint string, geminiKey string, chatgptKey string, cfg GromptConfig) []Provider {
	providers := gromptProviders.Initialize(
		claudeKey,
		openaiKey,
		deepseekKey,
		ollamaEndpoint,
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
