package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandOption struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Examples    []Example       `json:"examples,omitempty"`
	Subcommands []CommandOption `json:"subcommands,omitempty"`
	Flags       []FlagOption    `json:"flags,omitempty"`
}

type Example struct {
	Description string `json:"description"`
	Command     string `json:"command"`
}

type FlagOption struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
}

func init() {
	rootCmd.AddCommand(llmHelpCmd)
}

var llmHelpCmd = &cobra.Command{
	Use:    "llmhelp",
	Short:  "Print all commands and options in a machine-readable format for language models",
	Long:   "Print detailed information about all commands, flags, and examples in a JSON format that's easy for language models to parse and understand.",
	Hidden: false,
	Run: func(cmd *cobra.Command, args []string) {
		options := getCommandOptions(rootCmd)
		output, err := json.MarshalIndent(options, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling options: %v\n", err)
			return
		}
		fmt.Println(string(output))
	},
}

func getCommandOptions(cmd *cobra.Command) CommandOption {
	option := CommandOption{
		Name:        cmd.Name(),
		Description: cmd.Long,
	}

	if option.Description == "" {
		option.Description = cmd.Short
	}

	// Always add examples for known commands
	option.Examples = parseExamples(cmd.Name())

	// Get flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		required := false
		if cmd.Flags().Lookup(flag.Name) != nil {
			required = cmd.MarkFlagRequired(flag.Name) == nil
		}
		option.Flags = append(option.Flags, FlagOption{
			Name:        flag.Name,
			Shorthand:   flag.Shorthand,
			Description: flag.Usage,
			Type:        flag.Value.Type(),
			Required:    required,
		})
	})

	// Get subcommands
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			option.Subcommands = append(option.Subcommands, getCommandOptions(subCmd))
		}
	}

	return option
}

// parseExamples returns structured examples for a given command
func parseExamples(cmdName string) []Example {
	var examples []Example

	// Add default examples based on command name
	switch cmdName {
	case "config":
		examples = append(examples, []Example{
			{
				Description: "Set API token interactively",
				Command:     "hetznerdns config set",
			},
			{
				Description: "Set API token directly",
				Command:     "hetznerdns config set api-token YOUR_API_TOKEN",
			},
			{
				Description: "Show current configuration",
				Command:     "hetznerdns config show",
			},
		}...)
	case "zone":
		examples = append(examples, Example{
			Description: "List all DNS zones",
			Command:     "hetznerdns zone list",
		})
	case "record":
		examples = append(examples, []Example{
			{
				Description: "List records for a zone",
				Command:     "hetznerdns record list --zone example.com",
			},
			{
				Description: "Create an A record",
				Command:     "hetznerdns record create --zone example.com --name www --type A --value 192.168.1.1 --ttl 3600",
			},
			{
				Description: "Create a CNAME record",
				Command:     "hetznerdns record create --zone example.com --name blog --type CNAME --value example.com",
			},
			{
				Description: "Create an MX record",
				Command:     "hetznerdns record create --zone example.com --name @ --type MX --value \"10 mail.example.com\"",
			},
			{
				Description: "Update a record",
				Command:     "hetznerdns record update --id RECORD_ID --zone example.com --value 192.168.1.2",
			},
			{
				Description: "Delete a record",
				Command:     "hetznerdns record delete --id RECORD_ID",
			},
		}...)
	case "version":
		examples = append(examples, Example{
			Description: "Show version information",
			Command:     "hetznerdns version",
		})
	case "llmhelp":
		examples = append(examples, Example{
			Description: "Get machine-readable help information",
			Command:     "hetznerdns llmhelp",
		})
	}

	return examples
}
