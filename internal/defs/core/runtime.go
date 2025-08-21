// Package core provides the implementation of the Runtime interface.
package core

import (
	"github.com/rafa-mori/ghbex/internal/interfaces"
)

type Runtime struct {
	Debug      bool   `yaml:"debug" json:"debug"`
	DryRun     bool   `yaml:"dry_run" json:"dry_run"`
	ReportDir  string `yaml:"report_dir" json:"report_dir"`
	Background bool   `yaml:"background" json:"background"`
}

func NewRuntimeType(debug, dryRun bool, reportDir string, background bool) *Runtime {
	return &Runtime{
		Debug:      debug,
		DryRun:     dryRun,
		ReportDir:  reportDir,
		Background: background,
	}
}

func NewRuntime(debug, dryRun bool, reportDir string, background bool) interfaces.IRuntime {
	return NewRuntimeType(debug, dryRun, reportDir, background)
}

func (r *Runtime) GetDebug() bool                { return r.Debug }
func (r *Runtime) GetDryRun() bool               { return r.DryRun }
func (r *Runtime) GetReportDir() string          { return r.ReportDir }
func (r *Runtime) SetDebug(debug bool)           { r.Debug = debug }
func (r *Runtime) SetDryRun(dryRun bool)         { r.DryRun = dryRun }
func (r *Runtime) SetReportDir(reportDir string) { r.ReportDir = reportDir }
