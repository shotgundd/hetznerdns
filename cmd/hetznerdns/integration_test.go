package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain builds the CLI binary for testing
func TestMain(m *testing.M) {
	// Build the binary for testing
	buildCmd := exec.Command("go", "build", "-o", "hetznerdns_test")
	if err := buildCmd.Run(); err != nil {
		os.Stderr.WriteString("Failed to build test binary: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer os.Remove("hetznerdns_test")

	// Run the tests
	exitCode := m.Run()

	// Exit with the same code
	os.Exit(exitCode)
}

// runCommand runs the CLI command and returns stdout, stderr, and error
func runCommand(args ...string) (string, string, error) {
	cmd := exec.Command("./hetznerdns_test", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func TestVersionCommand(t *testing.T) {
	stdout, stderr, err := runCommand("version")
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "Hetzner DNS CLI v") {
		t.Errorf("Expected version output, got: %s", stdout)
	}
}

func TestConfigCommand(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "hetznerdns-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set environment variable to use the temp directory
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	// Test setting the API token (non-interactive mode)
	stdout, stderr, err := runCommand("config", "set", "api-token", "test-token")
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "Configuration saved") {
		t.Errorf("Expected success message, got: %s", stdout)
	}

	// Verify config file was created
	configPath := filepath.Join(tempDir, ".config", "hetznerdns", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Test showing the configuration
	stdout, stderr, err = runCommand("config", "show")
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr)
	}

	// The token is masked in the output, so we just check for "API token:" followed by something
	if !strings.Contains(stdout, "API token:") {
		t.Errorf("Expected to see API token information, got: %s", stdout)
	}

	// Make sure it's not showing "Not set"
	if strings.Contains(stdout, "Not set") {
		t.Errorf("API token should be set, but got: %s", stdout)
	}
}

// Note: The following tests require a valid API token and will make actual API calls.
// They are commented out by default and should be run manually when needed.

/*
func TestZoneCommands(t *testing.T) {
	// Skip if no API token is provided
	apiToken := os.Getenv("HETZNER_DNS_API_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping test because HETZNER_DNS_API_TOKEN is not set")
	}

	// Set up the API token for testing
	_, stderr, err := runCommand("config", "set", "api-token", apiToken)
	if err != nil {
		t.Fatalf("Failed to set API token: %v\nStderr: %s", err, stderr)
	}

	// Test listing zones
	stdout, stderr, err := runCommand("zone", "list")
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "ID") || !strings.Contains(stdout, "Name") {
		t.Errorf("Expected zone list output, got: %s", stdout)
	}
}

func TestRecordCommands(t *testing.T) {
	// Skip if no API token is provided
	apiToken := os.Getenv("HETZNER_DNS_API_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping test because HETZNER_DNS_API_TOKEN is not set")
	}

	// Set up the API token for testing
	_, stderr, err := runCommand("config", "set", "api-token", apiToken)
	if err != nil {
		t.Fatalf("Failed to set API token: %v\nStderr: %s", err, stderr)
	}

	// Get a zone ID for testing
	stdout, stderr, err := runCommand("zone", "list")
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr)
	}

	lines := strings.Split(stdout, "\n")
	if len(lines) < 2 {
		t.Skip("Skipping test because no zones are available")
	}

	// Extract zone ID from the first zone in the list
	fields := strings.Fields(lines[1])
	if len(fields) < 1 {
		t.Skip("Skipping test because zone ID could not be extracted")
	}
	zoneID := fields[0]

	// Test listing records for the zone
	stdout, stderr, err = runCommand("record", "list", "--zone-id", zoneID)
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "ID") || !strings.Contains(stdout, "Type") {
		t.Errorf("Expected record list output, got: %s", stdout)
	}
}
*/
