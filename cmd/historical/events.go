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

var historicalEventsCmd = &cobra.Command{
	Use:   "events <sport>",
	Short: "Get historical event listings",
	Args:  cobra.ExactArgs(1),
	Example: `  odds historical events basketball_nba --date 2024-11-30T12:00:00Z`,
	RunE: runHistoricalEvents,
}

func init() {
	historicalEventsCmd.Flags().String("date", "", "ISO 8601 timestamp (required)")
	historicalEventsCmd.MarkFlagRequired("date")

	HistoricalCmd.AddCommand(historicalEventsCmd)
}

func runHistoricalEvents(cmd *cobra.Command, args []string) error {
	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	params := url.Values{}
	date, _ := cmd.Flags().GetString("date")
	params.Set("date", date)

	path := fmt.Sprintf("/v4/historical/sports/%s/events", args[0])
	resp, err := c.Get(context.Background(), path, params)
	if err != nil {
		return err
	}

	data, err := client.Decode[model.HistoricalResponse[[]model.Event]](resp)
	if err != nil {
		return err
	}

	if isJSON(cmd) {
		return output.NewJSONWriter(os.Stdout).Write(data)
	}

	output.NewTableWriter(os.Stdout, isColor(cmd)).WriteHistoricalEvents(data)
	return nil
}
