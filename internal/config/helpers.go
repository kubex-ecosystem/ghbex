package config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"github.com/rafa-mori/ghbex/internal/interfaces"

	gl "github.com/rafa-mori/ghbex/internal/module/logger"

	"gopkg.in/yaml.v3"
)

func GetEnvOrDefault[T any](key string, defaultValue T) T {
	// Get environment variable and expand its
	value := os.Getenv(key)
	value = os.ExpandEnv(value)
	if value == "" {
		gl.Log("debug", "Environment variable %s not set, using default value: %v", key, defaultValue)
		return defaultValue
	}
	if _, ok := any(value).(T); ok {
		return any(value).(T)
	}
	return defaultValue
}

func GetBaseFilesPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		gl.Log("error", "Failed to get home directory: %v", err)
		return ""
	}
	basePath := filepath.Join(homeDir, ".kubex", "ghbex")
	return basePath
}

func EnsureDirs() error {
	configDir := filepath.Join(GetBaseFilesPath(), "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		gl.Log("error", "Failed to create config directory: %v", err)
		return err
	}
	reportDir := filepath.Join(GetBaseFilesPath(), "reports")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		gl.Log("error", "Failed to create report directory: %v", err)
		return err
	}
	return nil
}

func GetConfigFilePath(filePath string) string {
	if filePath == "" {
		filePath = filepath.Join(GetBaseFilesPath(), "config", "ghbex.yaml")
	}
	if _, err := os.Stat(filepath.Dir(filePath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			gl.Log("error", "Failed to create config directory: %v", err)
			return ""
		}
	}
	return filePath
}

func LoadEnvFromCurrentDir() {

	runtimePath, err := os.Executable()
	if err != nil {
		gl.Log("fatal", "Failed to get executable path: ", err)
	}
	if runtimePath == "" {
		runtimePath, err = filepath.Abs(".")
		if err != nil {
			gl.Log("fatal", "Failed to get current directory: ", err)
		}
	}
	gl.Log(
		"debug",
		"Loading environment variables from current directory (",
		runtimePath,
		")",
	)
	envFilePath := filepath.Join(filepath.Dir(runtimePath), ".env")
	if _, err := os.Stat(envFilePath); err == nil {
		gl.Log("debug", fmt.Sprintf("Loading environment variables from %s", envFilePath))
		envs, err := godotenv.Read(envFilePath)
		if err != nil {
			gl.Log("error", fmt.Sprintf("Failed to load environment variables: %v", err))
			return
		}
		gl.Log("debug", fmt.Sprintf("Loaded environment variables: %v", len(envs)))
		for key, value := range envs {
			gl.Log("debug", fmt.Sprintf("Setting environment variable %s=%s", key, value))
			err = os.Setenv(key, GetEnvOrDefault(key, value))
			if err != nil {
				gl.Log("error", fmt.Sprintf("Failed to set environment variable %s: %v", key, err))
			} else {
				gl.Log("debug", fmt.Sprintf("Set environment variable %s=%s", key, value))
			}
		}
	} else if !os.IsNotExist(err) {
		gl.Log("debug", fmt.Sprintf("Continuing without loading environment variables: %v", err))
	}
}

func LoadFromFile(filePath string) (interfaces.IMainConfig, error) {
	LoadEnvFromCurrentDir()
	filePath = GetConfigFilePath(filePath)
	err := EnsureDirs()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		gl.Log("warn", "Configuration file does not exist: %s", filePath)
		gl.Log("debug", "Creating a new configuration file.")
		cfg, err := NewMainConfig(
			GetEnvOrDefault("GHBEX_BIND_ADDR", "0.0.0.0"),
			GetEnvOrDefault("GHBEX_PORT", "8088"),
			GetEnvOrDefault("GHBEX_REPORT_DIR", "reports"),
			GetEnvOrDefault("GHBEX_OWNER", ""),
			GetEnvOrDefault("GHBEX_REPOSITORIES", []string{}),
			GetEnvOrDefault("GHBEX_DEBUG", false),
			GetEnvOrDefault("GHBEX_DISABLE_DRY_RUN", false),
			GetEnvOrDefault("GHBEX_BACKGROUND", true),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create new configuration: %v", err)
		}
		if err := SaveToFile(cfg, filePath); err != nil {
			return nil, fmt.Errorf("failed to create new configuration file: %v", err)
		}
		gl.Log("debug", "New configuration file created at: %s", filePath)
		return cfg, nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg MainConfig
	switch filepath.Ext(filePath) {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	case ".xml":
		if err := xml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", filepath.Ext(filePath))
	}
	return &cfg, nil
}

func SaveToFile(cfg interfaces.IMainConfig, filePath string) error {
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}
	filePath = GetConfigFilePath(filePath)
	var err error
	var data []byte

	switch filepath.Ext(filePath) {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return err
		}
	case ".json":
		data, err = json.Marshal(cfg)
		if err != nil {
			return err
		}
	case ".xml":
		data, err = xml.Marshal(cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", filepath.Ext(filePath))
	}

	// Ensure the directory exists before writing the file
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return err
		}
	} else {
		// File already exists, handle accordingly.
		// Ask the user if they want to overwrite it with 10s timeout
		if os.Getenv("FORCE") == "true" || os.Getenv("FORCE") == "y" {
			gl.Log("debug", "Overwriting existing file: %s", filePath)
			if err := os.WriteFile(filePath, data, 0644); err != nil {
				return err
			}
		} else {
			timeout := time.After(10 * time.Second)
			var overwrite string
			gl.Log("question", "File %s already exists. Overwrite? (y/n): ", filePath)
			select {
			case <-timeout:
				return fmt.Errorf("timeout")
			default:
				fmt.Scanln(&overwrite)
				if overwrite != "y" {
					return fmt.Errorf("aborted")
				} else {
					// Overwrite the file
					if err := os.WriteFile(filePath, data, 0644); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
