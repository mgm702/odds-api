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

var eventOddsCmd = &cobra.Command{
	Use:   "event-odds <sport> <event-id>",
	Short: "Get odds for a single event",
	Args:  cobra.ExactArgs(2),
	Example: `  odds event-odds basketball_nba evt001 --regions us --markets player_points
  odds event-odds basketball_nba evt001 --regions us --markets h2h,spreads,totals`,
	RunE: runEventOdds,
}

func init() {
	eventOddsCmd.Flags().String("regions", "", "Comma-delimited regions (required)")
	eventOddsCmd.Flags().String("markets", "", "Comma-delimited market keys (required)")
	_ = eventOddsCmd.MarkFlagRequired("regions")
	_ = eventOddsCmd.MarkFlagRequired("markets")
	rootCmd.AddCommand(eventOddsCmd)
}

func runEventOdds(cmd *cobra.Command, args []string) error {
	c, err := newRuntimeClient()
	if err != nil {
		return err
	}

	params := url.Values{}
	regions, _ := cmd.Flags().GetString("regions")
	params.Set("regions", regions)
	markets, _ := cmd.Flags().GetString("markets")
	params.Set("markets", markets)
	if oddsFormat != "" {
		params.Set("oddsFormat", oddsFormat)
	}
	if dateFormat != "iso" {
		params.Set("dateFormat", dateFormat)
	}

	path := fmt.Sprintf("/v4/sports/%s/events/%s/odds", args[0], args[1])
	resp, err := c.Get(context.Background(), path, params)
	if err != nil {
		return err
	}

	event, err := client.Decode[model.OddsEvent](resp)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(event)
	}

	tw := output.NewTableWriter(os.Stdout, useColor())
	tw.OddsFormat = oddsFormat
	tw.WriteOdds([]model.OddsEvent{event})
	return nil
}
