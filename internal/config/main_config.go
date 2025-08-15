// Package ghconfig provides configuration structures and interfaces for the application.
package ghconfig

// var configFilePath, bindHost, port, name string
// var debug, dryRun, background bool

type MainConfigImpl struct {
	Runtime   *Runtime   `yaml:"runtime"`
	Server    *Server    `yaml:"server"`
	GitHub    *GitHub    `yaml:"github"`
	Notifiers *Notifiers `yaml:"notifiers"`
}

type MainConfig interface {
	GetRuntime() *Runtime
	GetServer() *Server
	GetGitHub() *GitHub
	GetNotifiers() *Notifiers
}

func NewMainConfigObj() MainConfig {
	return &MainConfigImpl{
		Runtime:   &Runtime{},
		Server:    &Server{},
		GitHub:    &GitHub{},
		Notifiers: &Notifiers{},
	}
}

func (c *MainConfigImpl) GetRuntime() *Runtime {
	return c.Runtime
}

func (c *MainConfigImpl) GetServer() *Server {
	return c.Server
}

func (c *MainConfigImpl) GetGitHub() *GitHub {
	return c.GitHub
}

func (c *MainConfigImpl) GetNotifiers() *Notifiers {
	return c.Notifiers
}
