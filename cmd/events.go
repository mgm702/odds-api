package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
	"github.com/mgm702/odds-api-cli/internal/output"
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events <sport>",
	Short: "List events for a sport",
	Args:  cobra.ExactArgs(1),
	Example: `  odds events basketball_nba
  odds events basketball_nba --from 2024-12-01T00:00:00Z
  odds events basketball_nba --event-ids evt001,evt002`,
	RunE: runEvents,
}

func init() {
	eventsCmd.Flags().String("event-ids", "", "Comma-separated event IDs")
	eventsCmd.Flags().String("from", "", "Filter events starting at/after (ISO 8601)")
	eventsCmd.Flags().String("to", "", "Filter events starting at/before (ISO 8601)")
	rootCmd.AddCommand(eventsCmd)
}

func runEvents(cmd *cobra.Command, args []string) error {
	c, err := newRuntimeClient()
	if err != nil {
		return err
	}

	params := url.Values{}
	if dateFormat != "iso" {
		params.Set("dateFormat", dateFormat)
	}
	if v, _ := cmd.Flags().GetString("event-ids"); v != "" {
		params.Set("eventIds", v)
	}
	if v, _ := cmd.Flags().GetString("from"); v != "" {
		params.Set("commenceTimeFrom", v)
	}
	if v, _ := cmd.Flags().GetString("to"); v != "" {
		params.Set("commenceTimeTo", v)
	}

	path := fmt.Sprintf("/v4/sports/%s/events", args[0])
	resp, err := c.Get(context.Background(), path, params)
	if err != nil {
		return err
	}

	events, err := client.Decode[[]model.Event](resp)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(events)
	}

	output.NewTableWriter(os.Stdout, useColor()).WriteEvents(events)
	return nil
}
