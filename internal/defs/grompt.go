package defs

import "github.com/rafa-mori/grompt"

type Grompt = grompt.PromptEngine
type APIConfig = grompt.APIConfig
type Provider = grompt.Provider
type GromptConfig = grompt.Config
type PromptEngine = grompt.PromptEngine

func NewGromptEngine(cfg GromptConfig) Grompt {
	gmpt := grompt.NewPromptEngine(cfg)
	return gmpt
}
func NewGromptConfig(configFilePath string) GromptConfig {
	return grompt.DefaultConfig(configFilePath)
}
func NewPromptEngine(cfg GromptConfig) PromptEngine { return grompt.NewPromptEngine(cfg) }
func NewAPIConfig(configFilePath, provider string) APIConfig {
	gmptCfg := grompt.DefaultConfig(configFilePath)
	return gmptCfg.GetAPIConfig(provider)
}

// func NewProvider() Provider { return grompt.NewProvider() }
