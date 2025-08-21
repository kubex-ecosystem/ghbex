// Package interfaces defines the main configuration interface for the application.
package interfaces

type IMainConfig interface {
	GetConfigFilePath() string
	GetRuntime() IRuntime
	GetServer() IServer
	GetGitHub() IGitHub
	GetNotifiers() INotifiers
	String() string
}
