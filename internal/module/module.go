// Package module provides internal types and functions for the Ghbex application.
package module

import (
	cc "github.com/kubex-ecosystem/ghbex/cmd/cli"
	gl "github.com/kubex-ecosystem/ghbex/internal/module/logger"
	vs "github.com/kubex-ecosystem/ghbex/internal/module/version"
	"github.com/spf13/cobra"

	"os"
	"strings"
)

type Ghbex struct {
	parentCmdName string
	PrintBanner   bool
}

func (m *Ghbex) Alias() string {
	return ""
}
func (m *Ghbex) ShortDescription() string {
	return "GHbex is a command-line tool for managing GitHub repositories and file markers."
}
func (m *Ghbex) LongDescription() string {
	return `GHbex is a command-line tool for managing GitHub repositories and file markers.
`
}
func (m *Ghbex) Usage() string {
	return "ghbex [command] [args]"
}
func (m *Ghbex) Examples() []string {
	return []string{
		"ghbex operations analyze <owner> <repo>",
		"ghbex oper <command>",
	}
}
func (m *Ghbex) Active() bool {
	return true
}
func (m *Ghbex) Module() string {
	return "ghbex"
}
func (m *Ghbex) Execute() error {
	return m.Command().Execute()
}
func (m *Ghbex) Command() *cobra.Command {
	gl.Log("debug", "Starting Ghbex CLI...")

	var rtCmd = &cobra.Command{
		Use: m.Module(),
		// Aliases: []string{m.Alias()},
		Example: m.concatenateExamples(),
		Version: vs.GetVersion(),
		Annotations: cc.GetDescriptions([]string{
			m.LongDescription(),
			m.ShortDescription(),
		}, m.PrintBanner),
	}

	rtCmd.AddCommand(cc.OperationsCmdList())
	rtCmd.AddCommand(cc.ScoreCardRootCmd())
	rtCmd.AddCommand(vs.CliCommand())

	// Set usage definitions for the command and its subcommands
	setUsageDefinition(rtCmd)
	for _, c := range rtCmd.Commands() {
		setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}

	return rtCmd
}
func (m *Ghbex) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *Ghbex) concatenateExamples() string {
	examples := ""
	rtCmd := m.parentCmdName
	if rtCmd != "" {
		rtCmd = rtCmd + " "
	}
	for _, example := range m.Examples() {
		examples += rtCmd + example + "\n  "
	}
	return examples
}
