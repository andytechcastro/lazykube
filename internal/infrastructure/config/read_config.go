package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultNamespaces []string `json:"default_namespaces,omitempty"`
}

func ReadConfig() (*Config, error) {
	var config *Config
	home, err := os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("Error getting user home directory: %v", err)
	}

	configDir := filepath.Join(home, ".config", "lazykube")
	configFile := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return config, fmt.Errorf("Error creating config directory: %v", err)
		}

		// Create a default config
		configData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return config, fmt.Errorf("Error marshalling default config: %v", err)
		}

		if err := os.WriteFile(configFile, configData, 0644); err != nil {
			return config, fmt.Errorf("Error writing default config file: %v", err)
		}
		return config, nil
	} else {
		// Read the config file
		configData, err := os.ReadFile(configFile)
		if err != nil {
			return config, fmt.Errorf("Error reading config file: %v", err)
		}

		if err := json.Unmarshal(configData, &config); err != nil {
			return config, fmt.Errorf("Error unmarshalling config file: %v", err)
		}
	}
	return config, nil
}
