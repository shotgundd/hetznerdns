package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/shotgundd/hetznerdns/pkg/api"
	"github.com/shotgundd/hetznerdns/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.AddCommand(recordListCmd)
	recordCmd.AddCommand(recordCreateCmd)
	recordCmd.AddCommand(recordUpdateCmd)
	recordCmd.AddCommand(recordDeleteCmd)

	// Flags for record list command
	recordListCmd.Flags().StringP("zone", "z", "", "Zone ID or name (required)")
	recordListCmd.MarkFlagRequired("zone")

	// Flags for record create command
	recordCreateCmd.Flags().StringP("zone", "z", "", "Zone ID or name (required)")
	recordCreateCmd.Flags().StringP("name", "n", "", "Record name (required)")
	recordCreateCmd.Flags().StringP("type", "t", "", "Record type (A, AAAA, CNAME, MX, TXT, etc.) (required)")
	recordCreateCmd.Flags().StringP("value", "v", "", "Record value (required)")
	recordCreateCmd.Flags().IntP("ttl", "", 0, "Time to live in seconds (optional)")
	recordCreateCmd.MarkFlagRequired("zone")
	recordCreateCmd.MarkFlagRequired("name")
	recordCreateCmd.MarkFlagRequired("type")
	recordCreateCmd.MarkFlagRequired("value")

	// Flags for record update command
	recordUpdateCmd.Flags().StringP("id", "i", "", "Record ID (required)")
	recordUpdateCmd.Flags().StringP("zone", "z", "", "Zone ID or name (required)")
	recordUpdateCmd.Flags().StringP("name", "n", "", "Record name")
	recordUpdateCmd.Flags().StringP("type", "t", "", "Record type (A, AAAA, CNAME, MX, TXT, etc.)")
	recordUpdateCmd.Flags().StringP("value", "v", "", "Record value")
	recordUpdateCmd.Flags().IntP("ttl", "", 0, "Time to live in seconds")
	recordUpdateCmd.MarkFlagRequired("id")
	recordUpdateCmd.MarkFlagRequired("zone")

	// Flags for record delete command
	recordDeleteCmd.Flags().StringP("id", "i", "", "Record ID (required)")
	recordDeleteCmd.MarkFlagRequired("id")
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage DNS records",
	Long:  `Create, list, update, and delete DNS records.`,
}

// resolveZoneID tries to resolve a zone ID from either an ID or a name
func resolveZoneID(client *api.Client, zoneIDOrName string) (string, error) {
	// Try to resolve it as a name first
	fmt.Printf("Attempting to resolve '%s' as a zone name...\n", zoneIDOrName)

	// Get all zones to check against
	zones, err := client.GetZones()
	if err != nil {
		fmt.Printf("Error fetching zones: %v\n", err)
		return "", err
	}

	// First check if it's an exact match for a zone ID
	for _, zone := range zones {
		if zone.ID == zoneIDOrName {
			fmt.Printf("Found exact match for zone ID: %s (Name: %s)\n", zone.ID, zone.Name)
			return zone.ID, nil
		}
	}

	// Then check if it matches a zone name
	normalizedInput := strings.ToLower(strings.TrimSuffix(zoneIDOrName, "."))
	for _, zone := range zones {
		zoneName := strings.ToLower(strings.TrimSuffix(zone.Name, "."))
		if zoneName == normalizedInput {
			fmt.Printf("Found zone with name '%s', ID: %s\n", zone.Name, zone.ID)
			return zone.ID, nil
		}
	}

	// No match found
	fmt.Printf("Could not find any zone matching '%s'\n", zoneIDOrName)

	// Print available zones to help debugging
	fmt.Println("\nAvailable zones:")
	for _, zone := range zones {
		fmt.Printf("- %s (ID: %s)\n", zone.Name, zone.ID)
	}

	return "", fmt.Errorf("could not find zone with ID or name '%s'", zoneIDOrName)
}

var recordListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DNS records",
	Long:  `List all DNS records for a specific zone.`,
	Run: func(cmd *cobra.Command, args []string) {
		zoneIDOrName, _ := cmd.Flags().GetString("zone")

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

		// Resolve zone ID from name if needed
		zoneID, err := resolveZoneID(client, zoneIDOrName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		records, err := client.GetRecords(zoneID)
		if err != nil {
			fmt.Printf("Error fetching records: %v\n", err)
			return
		}

		if len(records) == 0 {
			fmt.Println("No records found for this zone.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tVALUE\tTTL")
		for _, record := range records {
			ttl := strconv.Itoa(record.TTL)
			if record.TTL == 0 {
				ttl = "default"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", record.ID, record.Name, record.Type, record.Value, ttl)
		}
		w.Flush()
	},
}

var recordCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS record",
	Long:  `Create a new DNS record in a specific zone.`,
	Run: func(cmd *cobra.Command, args []string) {
		zoneIDOrName, _ := cmd.Flags().GetString("zone")
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		value, _ := cmd.Flags().GetString("value")
		ttl, _ := cmd.Flags().GetInt("ttl")

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

		// Resolve zone ID from name if needed
		zoneID, err := resolveZoneID(client, zoneIDOrName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		record := api.Record{
			ZoneID: zoneID,
			Name:   name,
			Type:   recordType,
			Value:  value,
			TTL:    ttl,
		}

		createdRecord, err := client.CreateRecord(record)
		if err != nil {
			fmt.Printf("Error creating record: %v\n", err)
			return
		}

		fmt.Printf("Record created successfully with ID: %s\n", createdRecord.ID)
	},
}

var recordUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a DNS record",
	Long:  `Update an existing DNS record.`,
	Run: func(cmd *cobra.Command, args []string) {
		recordID, _ := cmd.Flags().GetString("id")
		zoneIDOrName, _ := cmd.Flags().GetString("zone")
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		value, _ := cmd.Flags().GetString("value")
		ttl, _ := cmd.Flags().GetInt("ttl")

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

		// Resolve zone ID from name if needed
		zoneID, err := resolveZoneID(client, zoneIDOrName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		record := api.Record{
			ID:     recordID,
			ZoneID: zoneID,
		}

		// Only update fields that were provided
		if name != "" {
			record.Name = name
		}
		if recordType != "" {
			record.Type = recordType
		}
		if value != "" {
			record.Value = value
		}
		if cmd.Flags().Changed("ttl") {
			record.TTL = ttl
		}

		updatedRecord, err := client.UpdateRecord(record)
		if err != nil {
			fmt.Printf("Error updating record: %v\n", err)
			return
		}

		fmt.Printf("Record updated successfully: %s\n", updatedRecord.ID)
	},
}

var recordDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DNS record",
	Long:  `Delete an existing DNS record.`,
	Run: func(cmd *cobra.Command, args []string) {
		recordID, _ := cmd.Flags().GetString("id")

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
		err = client.DeleteRecord(recordID)
		if err != nil {
			fmt.Printf("Error deleting record: %v\n", err)
			return
		}

		fmt.Println("Record deleted successfully.")
	},
}
