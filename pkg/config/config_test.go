package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadAndSaveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "hetznerdns-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override config paths for testing
	origConfigDir := configDir
	origConfigFile := configFile
	defer func() {
		configDir = origConfigDir
		configFile = origConfigFile
	}()

	configDir = filepath.Join(tempDir, ".config", "hetznerdns")
	configFile = filepath.Join(configDir, "config.yaml")

	// Save original viper instance and create a new one for testing
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("HETZNER_DNS")
	v.AutomaticEnv()
	v.SetDefault("api_token", "")

	// Create a clean environment
	os.Unsetenv("HETZNER_DNS_API_TOKEN")

	// Test loading config when file doesn't exist
	config := &Config{
		APIToken: v.GetString("api_token"),
	}

	if config.APIToken != "" {
		t.Errorf("Expected empty API token, got '%s'", config.APIToken)
	}

	// Test saving config
	config.APIToken = "test-token"
	v.Set("api_token", config.APIToken)

	// Create config directory
	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		t.Fatalf("Error creating config directory: %v", err)
	}

	// Write config file
	if err := v.WriteConfigAs(configFile); err != nil {
		t.Fatalf("Error writing config file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("Config file was not created")
	}

	// Test loading config from file
	v2 := viper.New()
	v2.SetConfigFile(configFile)
	v2.SetConfigType("yaml")

	if err := v2.ReadInConfig(); err != nil {
		t.Fatalf("Error reading config file: %v", err)
	}

	loadedConfig := &Config{
		APIToken: v2.GetString("api_token"),
	}

	if loadedConfig.APIToken != "test-token" {
		t.Errorf("Expected API token 'test-token', got '%s'", loadedConfig.APIToken)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "hetznerdns-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override config paths for testing
	origConfigDir := configDir
	origConfigFile := configFile
	defer func() {
		configDir = origConfigDir
		configFile = origConfigFile
	}()

	configDir = filepath.Join(tempDir, ".config", "hetznerdns")
	configFile = filepath.Join(configDir, "config.yaml")

	// Save original environment and set test environment
	origEnv := os.Getenv("HETZNER_DNS_API_TOKEN")
	defer os.Setenv("HETZNER_DNS_API_TOKEN", origEnv)

	os.Setenv("HETZNER_DNS_API_TOKEN", "env-test-token")

	// Create a new viper instance for testing
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("HETZNER_DNS")
	v.AutomaticEnv()
	v.SetDefault("api_token", "")

	// Test loading config from environment
	config := &Config{
		APIToken: v.GetString("api_token"),
	}

	if config.APIToken != "env-test-token" {
		t.Errorf("Expected API token 'env-test-token', got '%s'", config.APIToken)
	}
}
