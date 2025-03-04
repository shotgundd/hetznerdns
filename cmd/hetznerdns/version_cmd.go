package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set during build using ldflags
var Version = "dev"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the version of the Hetzner DNS CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hetzner DNS CLI v%s\n", Version)
	},
}
