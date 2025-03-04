package main

import (
	"fmt"

	"github.com/shotgundd/hetznerdns/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)

	// Add api-token argument to set command
	configSetCmd.Flags().StringP("api-token", "t", "", "API token for Hetzner DNS")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Configure the Hetzner DNS CLI tool.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  `Set configuration values like API token.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if api-token flag is provided
		apiToken, _ := cmd.Flags().GetString("api-token")

		// If not provided via flag, prompt for it
		if apiToken == "" {
			if len(args) >= 2 && args[0] == "api-token" {
				// Support for command-line arguments: config set api-token VALUE
				apiToken = args[1]
			} else {
				// Interactive mode
				fmt.Print("Enter your Hetzner DNS API token: ")
				fmt.Scanln(&apiToken)
			}
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		cfg.APIToken = apiToken

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Println("Configuration saved successfully.")
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration values.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.APIToken == "" {
			fmt.Println("API token: Not set")
		} else {
			// Only show the first and last few characters of the token for security
			tokenLen := len(cfg.APIToken)
			if tokenLen > 8 {
				maskedToken := cfg.APIToken[0:4] + "..." + cfg.APIToken[tokenLen-4:tokenLen]
				fmt.Printf("API token: %s\n", maskedToken)
			} else {
				fmt.Println("API token: ********")
			}
		}
	},
}
