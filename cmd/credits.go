package cmd

import (
	"context"
	"os"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/output"
	"github.com/spf13/cobra"
)

var creditsCmd = &cobra.Command{
	Use:   "credits",
	Short: "Show API quota usage report",
	Example: `  odds credits
  odds credits --json`,
	RunE: runCredits,
}

func init() {
	rootCmd.AddCommand(creditsCmd)
}

func runCredits(cmd *cobra.Command, args []string) error {
	key, err := getAPIKey()
	if err != nil {
		return err
	}

	c := client.New(key)
	c.Verbose = verbose

	quota, err := c.GetQuotaOnly(context.Background())
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(quota)
	}

	output.NewTableWriter(os.Stdout, useColor()).WriteCredits(quota)
	return nil
}
