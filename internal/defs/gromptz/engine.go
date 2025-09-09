// Package gromptz provides a set of types and functions for working with the Grompt library.
package gromptz

import (
	"github.com/kubex-ecosystem/gemx/grompt"
	gromptConfig "github.com/kubex-ecosystem/gemx/grompt/factory/config"
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
	bindAddr string,
	port string,
	claudeKey string,
	openaiKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	geminiKey string,
	chatgptKey string,
) GromptConfig {
	return gromptConfig.NewConfig(bindAddr, port, openaiKey, deepSeekKey, ollamaEndpoint, claudeKey, geminiKey, chatgptKey, nil)
}
func NewPromptEngine(cfg GromptConfig) PromptEngine { return grompt.NewPromptEngine(cfg) }
func NewAPIConfig(configFilePath, provider string) APIConfig {
	gmptCfg := grompt.DefaultConfig(configFilePath)
	return gmptCfg.GetAPIConfig(provider)
}
