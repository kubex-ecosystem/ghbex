package cli

import (
	"context"
	"log"

	config "github.com/rafa-mori/ghbex/internal/config"
	ghserver "github.com/rafa-mori/ghbex/internal/server"
	"github.com/spf13/cobra"
)

func ServerCmdList() []*cobra.Command {
	var cmds []*cobra.Command

	// Define your server commands here
	cmds = append(cmds, startServer())
	cmds = append(cmds, stopServer())
	cmds = append(cmds, statusServer())
	return cmds
}

func stopServer() *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the server",
		Annotations: GetDescriptions([]string{
			"This command stops server.",
			"This command stops the Grompt server and releases any resources held by it.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Stop server logic
		},
	}
	return stopCmd
}

func statusServer() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get server status",
		Annotations: GetDescriptions([]string{
			"This command gets the status of the server.",
			"This command checks if the server is running and returns its status.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Get server status logic
		},
	}
	return statusCmd
}

func startServer() *cobra.Command {
	var configFilePath, bindAddr, port, name string
	var debug, dryRun, background bool

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Annotations: GetDescriptions([]string{
			"This command starts the server.",
			"This command initializes the server and starts waiting for help to build prompts.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Check if config file path was provided
			configPath := "config/sanitize.yaml"
			if configFilePath != "" {
				configPath = configFilePath
			}

			cfg, err := config.LoadFromFile(configPath)
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			// Start server logic
			ctx := context.Background()
			// Initialize the server
			srv := ghserver.NewGHServerEngine(cfg)
			// Start the server
			if err := srv.Start(ctx); err != nil {
				log.Fatal(err)
			}
		},
	}

	// Define flags for the command
	startCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "Path to the configuration file")
	startCmd.Flags().StringVarP(&bindAddr, "host", "H", "localhost", "Host to bind the server")
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server")
	startCmd.Flags().StringVarP(&name, "name", "n", "Grompt", "Name of the server")
	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	startCmd.Flags().BoolVarP(&dryRun, "dry-run", "D", false, "Enable dry run mode")
	startCmd.Flags().BoolVarP(&background, "background", "b", false, "Run in background")

	return startCmd
}
