package historical

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

var historicalEventOddsCmd = &cobra.Command{
	Use:   "event-odds <sport> <event-id>",
	Short: "Get historical odds for a single event",
	Args:  cobra.ExactArgs(2),
	Example: `  odds historical event-odds basketball_nba evt001 --date 2024-11-30T12:00:00Z --regions us --markets h2h`,
	RunE: runHistoricalEventOdds,
}

func init() {
	historicalEventOddsCmd.Flags().String("date", "", "ISO 8601 timestamp (required)")
	historicalEventOddsCmd.Flags().String("regions", "", "Comma-delimited regions (required)")
	historicalEventOddsCmd.Flags().String("markets", "", "Comma-delimited markets (required)")
	historicalEventOddsCmd.Flags().String("odds-format", "", "Odds format: decimal or american")
	historicalEventOddsCmd.MarkFlagRequired("date")
	historicalEventOddsCmd.MarkFlagRequired("regions")
	historicalEventOddsCmd.MarkFlagRequired("markets")

	HistoricalCmd.AddCommand(historicalEventOddsCmd)
}

func runHistoricalEventOdds(cmd *cobra.Command, args []string) error {
	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	params := url.Values{}
	date, _ := cmd.Flags().GetString("date")
	params.Set("date", date)
	regions, _ := cmd.Flags().GetString("regions")
	params.Set("regions", regions)
	markets, _ := cmd.Flags().GetString("markets")
	params.Set("markets", markets)
	if v, _ := cmd.Flags().GetString("odds-format"); v != "" {
		params.Set("oddsFormat", v)
	}

	path := fmt.Sprintf("/v4/historical/sports/%s/events/%s/odds", args[0], args[1])
	resp, err := c.Get(context.Background(), path, params)
	if err != nil {
		return err
	}

	data, err := client.Decode[model.HistoricalResponse[model.OddsEvent]](resp)
	if err != nil {
		return err
	}

	if isJSON(cmd) {
		return output.NewJSONWriter(os.Stdout).Write(data)
	}

	tw := output.NewTableWriter(os.Stdout, isColor(cmd))
	fmt.Fprintf(os.Stdout, "Snapshot: %s\n", data.Timestamp)
	if data.PreviousTimestamp != nil {
		fmt.Fprintf(os.Stdout, "Previous: %s\n", *data.PreviousTimestamp)
	}
	if data.NextTimestamp != nil {
		fmt.Fprintf(os.Stdout, "Next:     %s\n", *data.NextTimestamp)
	}
	fmt.Fprintln(os.Stdout)
	tw.WriteOdds([]model.OddsEvent{data.Data})
	return nil
}
