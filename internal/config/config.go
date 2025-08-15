// Package config provides configuration structures and interfaces for the application.
package config

import "github.com/rafa-mori/ghbex/internal/defs"

// var configFilePath, bindHost, port, name string
// var debug, dryRun, background bool

type MainConfigImpl struct {
	Runtime   *defs.Runtime   `yaml:"runtime"`
	Server    *defs.Server    `yaml:"server"`
	GitHub    *defs.GitHub    `yaml:"github"`
	Notifiers *defs.Notifiers `yaml:"notifiers"`
}

type MainConfig interface {
	GetRuntime() *defs.Runtime
	GetServer() *defs.Server
	GetGitHub() *defs.GitHub
	GetNotifiers() *defs.Notifiers
}

func NewMainConfigObj() MainConfig {
	return &MainConfigImpl{
		Runtime:   &defs.Runtime{},
		Server:    &defs.Server{},
		GitHub:    &defs.GitHub{},
		Notifiers: &defs.Notifiers{},
	}
}

func (c *MainConfigImpl) GetRuntime() *defs.Runtime {
	return c.Runtime
}

func (c *MainConfigImpl) GetServer() *defs.Server {
	return c.Server
}

func (c *MainConfigImpl) GetGitHub() *defs.GitHub {
	return c.GitHub
}

func (c *MainConfigImpl) GetNotifiers() *defs.Notifiers {
	return c.Notifiers
}
