package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shotgundd/hetznerdns/pkg/api"
	"github.com/shotgundd/hetznerdns/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(zoneCmd)
	zoneCmd.AddCommand(zoneListCmd)
}

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Manage DNS zones",
	Long:  `Create, list, update, and delete DNS zones.`,
}

var zoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DNS zones",
	Long:  `List all DNS zones in your Hetzner account.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if cfg.APIToken == "" {
			fmt.Println("API token not set. Please run 'hetznerdns config set' to configure your API token.")
			return
		}

		client := api.NewClient(cfg.APIToken)
		zones, err := client.GetZones()
		if err != nil {
			fmt.Printf("Error fetching zones: %v\n", err)
			return
		}

		if len(zones) == 0 {
			fmt.Println("No zones found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTTL\tRECORDS")
		for _, zone := range zones {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", zone.ID, zone.Name, zone.TTL, zone.RecordsCount)
		}
		w.Flush()
	},
}
