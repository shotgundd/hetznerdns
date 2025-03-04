package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the configuration for the application
type Config struct {
	APIToken string
}

// Default config paths
var (
	configDir  string
	configFile string
)

// LoadConfig loads the configuration from the config file or environment variables
func LoadConfig() (*Config, error) {
	// Set default config file paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	configDir = filepath.Join(homeDir, ".config", "hetznerdns")
	configFile = filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("error creating config directory: %w", err)
		}
	}

	// Set up viper
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set environment variable prefix
	viper.SetEnvPrefix("HETZNER_DNS")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("api_token", "")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return an error if it's not a ConfigFileNotFoundError
			// This makes it friendlier when the config file doesn't exist yet
			fmt.Printf("Notice: Config file not found, will create a new one when you save settings.\n")
		}
		// Config file not found, will use defaults and env vars
	}

	// Create config struct
	config := &Config{
		APIToken: viper.GetString("api_token"),
	}

	return config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	viper.Set("api_token", config.APIToken)

	// Check if the config file exists
	fileExists := false
	if _, err := os.Stat(configFile); err == nil {
		fileExists = true
	}

	if fileExists {
		// Config file exists, use WriteConfig
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("error writing config file: %w", err)
		}
	} else {
		// Config file doesn't exist, write directly to the file
		if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}

		// Create the config file with the API token
		content := fmt.Sprintf("api_token: %s\n", config.APIToken)
		if err := os.WriteFile(configFile, []byte(content), 0600); err != nil {
			return fmt.Errorf("error writing config file: %w", err)
		}
		fmt.Printf("Created config file at: %s\n", configFile)
	}

	return nil
}
