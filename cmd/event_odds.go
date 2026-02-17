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
	eventOddsCmd.Flags().String("odds-format", "", "Odds format: decimal or american")
	eventOddsCmd.MarkFlagRequired("regions")
	eventOddsCmd.MarkFlagRequired("markets")
	rootCmd.AddCommand(eventOddsCmd)
}

func runEventOdds(cmd *cobra.Command, args []string) error {
	key, err := getAPIKey()
	if err != nil {
		return err
	}

	c := client.New(key)
	c.Verbose = verbose

	params := url.Values{}
	regions, _ := cmd.Flags().GetString("regions")
	params.Set("regions", regions)
	markets, _ := cmd.Flags().GetString("markets")
	params.Set("markets", markets)
	if v, _ := cmd.Flags().GetString("odds-format"); v != "" {
		params.Set("oddsFormat", v)
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

	output.NewTableWriter(os.Stdout, useColor()).WriteOdds([]model.OddsEvent{event})
	return nil
}
