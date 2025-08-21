package cli

import (
	"github.com/spf13/cobra"
)

func OperationsCmdList() *cobra.Command {
	var cmds []*cobra.Command
	short := "Manage repositories, PRs, Workflows, issues, and much more."
	long := "This command provides a comprehensive set of operations for managing GitHub repositories, including analysis, automation, and productivity enhancements."

	operationsCmd := &cobra.Command{
		Use:     "operations",
		Aliases: []string{"ops", "oper", "op"},
		Short:   short,
		Long:    long,
		Annotations: GetDescriptions([]string{
			short,
			long,
		}, false),
	}

	// Define your server commands here
	cmds = append(cmds, analyzeServer())
	// Add more commands as needed
	operationsCmd.AddCommand(cmds...)

	return operationsCmd
}

func analyzeServer() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze repositories, bringing insights and recommendations.",
		Annotations: GetDescriptions([]string{
			"This command analyzes the specified repositories.",
			"This command provides insights and recommendations based on the analysis.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Analyze server logic
		},
	}
	return analyzeCmd
}
