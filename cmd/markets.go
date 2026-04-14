package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
	"github.com/mgm702/odds-api-cli/internal/output"
	"github.com/spf13/cobra"
)

var marketsCmd = &cobra.Command{
	Use:   "markets <sport> <event-id>",
	Short: "List available markets for an event",
	Args:  cobra.ExactArgs(2),
	Example: `  odds markets basketball_nba evt001
  odds markets basketball_nba evt001 --json`,
	RunE: runMarkets,
}

func init() {
	rootCmd.AddCommand(marketsCmd)
}

func runMarkets(cmd *cobra.Command, args []string) error {
	c, err := newRuntimeClient()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v4/sports/%s/events/%s/markets", args[0], args[1])
	resp, err := c.Get(context.Background(), path, nil)
	if err != nil {
		return err
	}

	em, err := client.Decode[model.EventMarkets](resp)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(em)
	}

	output.NewTableWriter(os.Stdout, useColor()).WriteEventMarkets(em)
	return nil
}
