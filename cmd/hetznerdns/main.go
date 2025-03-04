package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hetznerdns",
	Short: "A CLI tool to manage Hetzner DNS records",
	Long: `hetznerdns is a command line tool that allows you to create, 
read, update, and delete DNS records on Hetzner DNS service.`,
}

func init() {
	// Add commands here
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
