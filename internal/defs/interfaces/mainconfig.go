// Package interfaces defines the main configuration interface for the application.
package interfaces

import (
	"github.com/rafa-mori/grompt"
)

type IMainConfig interface {
	GetConfigFilePath() string
	GetRuntime() IRuntime
	GetGitHub() IGitHub
	GetNotifiers() INotifiers
	GetGrompt() grompt.PromptEngine
	GetConfigObject() any
	String() string
}
