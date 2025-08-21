package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	config "github.com/rafa-mori/ghbex/internal/config"
	"github.com/rafa-mori/ghbex/internal/interfaces"
	gl "github.com/rafa-mori/ghbex/internal/module/logger"
	ghserver "github.com/rafa-mori/ghbex/internal/server"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

func ServerCmdList() []*cobra.Command {
	var cmds []*cobra.Command

	// Define your server commands here
	cmds = append(cmds, startServer())
	cmds = append(cmds, stopServer())
	cmds = append(cmds, statusServer())
	cmds = append(cmds, configServer())
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

func configServer() *cobra.Command {
	var configFilePath, format string

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Get server configuration",
		Annotations: GetDescriptions([]string{
			"This command gets the configuration of the server.",
			"This command checks the current configuration settings of the server.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Get server configuration logic
			cfg, err := config.LoadFromFile(configFilePath)
			if err != nil {
				gl.Log("error", "Failed to load config: %v", err)
				return
			}
			// Additional logic to display or use the configuration
			if cfg != nil {
				switch format {
				case "yaml", "yml", "y":
					printYAMLConfig(cfg)
				case "json", "j":
					printJSONConfig(cfg)
				default:
					printTreeConfig(cfg)
				}
			} else {
				gl.Log("warn", "No configuration found")
			}
		},
	}

	configCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "Path to the configuration file")
	configCmd.Flags().StringVarP(&format, "format", "f", "tree", "Output format (tree/json/yaml)")

	return configCmd
}

func startServer() *cobra.Command {
	var configFilePath, bindAddr, port, name, reportDir string
	var debug, disableDryRun, background bool

	startCmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"server", "run"},
		Short:   "Start the server",
		Annotations: GetDescriptions([]string{
			"This command starts the server.",
			"This command initializes the server and starts waiting for help to build prompts.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			var cfg interfaces.IMainConfig
			if configFilePath == "" {
				if bindAddr == "" && port == "" {
					var err error
					cfg, err = config.LoadFromFile("")
					if err != nil {
						log.Fatalf("Failed to load config: %v", err)
					}
				} else {
					// Create a new config object with provided flags
					cfg = config.NewMainConfig(bindAddr, port, reportDir, debug, !disableDryRun, background)
					// Save the config to a file
					if err := config.SaveToFile(cfg, ""); err != nil {
						log.Fatalf("Failed to save config: %v", err)
					}
				}
			}

			// Start server logic
			ctx := context.Background()
			if cfg == nil {
				log.Fatal("Configuration is nil")
			}
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
	startCmd.Flags().StringVarP(&bindAddr, "host", "H", "", "Host to bind the server")
	startCmd.Flags().StringVarP(&port, "port", "p", "", "Port to run the server")
	startCmd.Flags().StringVarP(&name, "name", "n", "Grompt", "Name of the server")
	startCmd.Flags().StringVarP(&reportDir, "report-dir", "r", "", "Directory to store reports")
	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	startCmd.Flags().BoolVarP(&disableDryRun, "disable-dry-run", "D", false, "Disable dry run mode (Default: false - Dry run by default)")
	startCmd.Flags().BoolVarP(&background, "background", "b", false, "Run in background")

	return startCmd
}

func printYAMLConfig(cfg interfaces.IMainConfig) {
	gl.Log("answer", "GHbex Settings (YAML):")
	// Use a YAML library to marshal the config into YAML format
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		gl.Log("error", "Failed to marshal config to YAML: %v", err)
		return
	}
	gl.Log("answer", string(yamlData))
}

func printJSONConfig(cfg interfaces.IMainConfig) {
	gl.Log("answer", "GHbex Settings (JSON):")
	// Use a JSON library to marshal the config into JSON format
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		gl.Log("error", "Failed to marshal config to JSON: %v", err)
		return
	}
	gl.Log("answer", string(jsonData))
}

func printTreeConfig(cfg interfaces.IMainConfig) {
	firstLevelMid := "├─ "
	firstLevelEnd := "└─ "
	secLevelMid := "│    ├─ "
	secLevelEnd := "│    └─ "
	thirdLevelMid := "│       ├─ "
	thirdLevelEnd := "│       └─ "
	gl.Log("answer", "GHbex Settings:")
	gl.Log("answer", fmt.Sprintf("%sServer:", firstLevelMid))
	gl.Log("answer", fmt.Sprintf("%sAddr: %s", secLevelMid, cfg.GetServer().GetAddr()))
	gl.Log("answer", fmt.Sprintf("%sPort: %s", secLevelEnd, cfg.GetServer().GetPort()))
	gl.Log("answer", fmt.Sprintf("%sRuntime:", firstLevelMid))
	gl.Log("answer", fmt.Sprintf("%sReport Directory: %s", secLevelMid, cfg.GetRuntime().GetReportDir()))
	gl.Log("answer", fmt.Sprintf("%sDry Run: %t", secLevelMid, cfg.GetRuntime().GetDryRun()))
	gl.Log("answer", fmt.Sprintf("%sDebug: %t", secLevelEnd, cfg.GetRuntime().GetDebug()))
	gl.Log("answer", fmt.Sprintf("%sGitHub:", firstLevelMid))
	gl.Log("answer", fmt.Sprintf("%sAuth Kind: %s", secLevelMid, cfg.GetGitHub().GetAuth().GetKind()))
	gl.Log("answer", fmt.Sprintf("%sInstallation ID: %d", secLevelMid, cfg.GetGitHub().GetAuth().GetInstallationID()))
	gl.Log("answer", fmt.Sprintf("%sPrivate Key Path: %s", secLevelMid, cfg.GetGitHub().GetAuth().GetPrivateKeyPath()))
	gl.Log("answer", fmt.Sprintf("%sUpload URL: %s", secLevelMid, cfg.GetGitHub().GetAuth().GetUploadURL()))
	gl.Log("answer", fmt.Sprintf("%sBase URL: %s", secLevelMid, cfg.GetGitHub().GetAuth().GetBaseURL()))
	gl.Log("answer", fmt.Sprintf("%sRepo list (%d):", secLevelEnd, len(cfg.GetGitHub().GetRepos())))
	for i, repo := range cfg.GetGitHub().GetRepos() {
		if i >= len(cfg.GetGitHub().GetRepos())-1 {
			gl.Log("answer", fmt.Sprintf("%s[%d]: %s/%s", thirdLevelEnd, i, repo.GetOwner(), repo.GetName()))
		} else {
			gl.Log("answer", fmt.Sprintf("%s[%d]: %s/%s,", thirdLevelMid, i, repo.GetOwner(), repo.GetName()))
		}
	}
	gl.Log("answer", fmt.Sprintf("%sNotifiers (%d):", firstLevelEnd, len(cfg.GetNotifiers().GetNotifiers())))
	for i, notifier := range cfg.GetNotifiers().GetNotifiers() {
		webhook := notifier.GetWebhook()
		if webhook == "" {
			webhook = notifier.GetType()
		}
		if i < len(cfg.GetNotifiers().GetNotifiers())-1 {
			gl.Log("answer", fmt.Sprintf("    %s[%d] (%s): %s", firstLevelEnd, i, notifier.GetType(), webhook))
		} else {
			gl.Log("answer", fmt.Sprintf("    %s[%d] (%s): %s", firstLevelEnd, i, notifier.GetType(), webhook))
		}
	}
}
