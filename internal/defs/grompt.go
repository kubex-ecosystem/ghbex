package defs

import (
	"github.com/rafa-mori/grompt"
	gromptConfig "github.com/rafa-mori/grompt/factory/config"
	gromptProviders "github.com/rafa-mori/grompt/factory/providers"
)

type Grompt = grompt.PromptEngine
type APIConfig = grompt.APIConfig
type Provider = gromptProviders.Provider
type GromptConfig = grompt.Config
type PromptEngine = grompt.PromptEngine

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

func NewProviders(claudeKey string, openaiKey string, deepseekKey string, ollamaEndpoint string) []Provider {
	providers := gromptProviders.Initialize(
		claudeKey,
		openaiKey,
		deepseekKey,
		ollamaEndpoint,
	)

	return providers
}
