// Package ghbex provides a set of utilities for working with GitHub repositories.
package ghbex

import (
	"github.com/rafa-mori/ghbex/internal/config"
	"github.com/rafa-mori/ghbex/internal/defs/gitz"
	"github.com/rafa-mori/ghbex/internal/interfaces"
	ghserver "github.com/rafa-mori/ghbex/internal/server"
)

type MainConfig = interfaces.IMainConfig
type GHServerEngine = ghserver.GHServerEngine
type GitHub = gitz.GitHub

func NewMainConfigObj() (MainConfig, error) {
	return config.NewMainConfig(
		"",
		"",
		"",
		"",
		[]string{},
		false,
		true,
		false,
	)
}
