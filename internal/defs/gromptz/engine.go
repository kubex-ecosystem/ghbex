// Package grompt provides a set of types and functions for working with the Grompt library.
package gromptz

import (
	"github.com/rafa-mori/grompt"
	gromptConfig "github.com/rafa-mori/grompt/factory/config"
)

type IGrompt = grompt.Grompt
type Grompt = grompt.PromptEngine

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
