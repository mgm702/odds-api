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

var HistoricalCmd = &cobra.Command{
	Use:   "historical",
	Short: "Query historical odds and events",
}

var historicalOddsCmd = &cobra.Command{
	Use:   "odds <sport>",
	Short: "Get historical odds snapshot",
	Args:  cobra.ExactArgs(1),
	Example: `  odds historical odds basketball_nba --date 2024-11-30T12:00:00Z --regions us
  odds historical odds basketball_nba --date 2024-11-30T12:00:00Z --regions us --markets h2h,spreads`,
	RunE: runHistoricalOdds,
}

func init() {
	historicalOddsCmd.Flags().String("date", "", "ISO 8601 timestamp (required)")
	historicalOddsCmd.Flags().String("regions", "", "Comma-delimited regions (required)")
	historicalOddsCmd.Flags().String("markets", "", "Comma-delimited markets")
	historicalOddsCmd.Flags().String("odds-format", "", "Odds format: decimal or american")
	historicalOddsCmd.MarkFlagRequired("date")
	historicalOddsCmd.MarkFlagRequired("regions")

	HistoricalCmd.AddCommand(historicalOddsCmd)
}

func runHistoricalOdds(cmd *cobra.Command, args []string) error {
	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	params := url.Values{}
	date, _ := cmd.Flags().GetString("date")
	params.Set("date", date)
	regions, _ := cmd.Flags().GetString("regions")
	params.Set("regions", regions)
	if v, _ := cmd.Flags().GetString("markets"); v != "" {
		params.Set("markets", v)
	}
	if v, _ := cmd.Flags().GetString("odds-format"); v != "" {
		params.Set("oddsFormat", v)
	}

	path := fmt.Sprintf("/v4/historical/sports/%s/odds", args[0])
	resp, err := c.Get(context.Background(), path, params)
	if err != nil {
		return err
	}

	data, err := client.Decode[model.HistoricalResponse[[]model.OddsEvent]](resp)
	if err != nil {
		return err
	}

	if isJSON(cmd) {
		return output.NewJSONWriter(os.Stdout).Write(data)
	}

	output.NewTableWriter(os.Stdout, isColor(cmd)).WriteHistoricalOdds(data)
	return nil
}

func newClient(cmd *cobra.Command) (*client.Client, error) {
	key, _ := cmd.Root().PersistentFlags().GetString("api-key")
	if key == "" {
		key = os.Getenv("ODDS_API_KEY")
	}
	if key == "" {
		return nil, fmt.Errorf("API key required: set ODDS_API_KEY or use --api-key")
	}
	c := client.New(key)
	v, _ := cmd.Root().PersistentFlags().GetBool("verbose")
	c.Verbose = v
	return c, nil
}

func isJSON(cmd *cobra.Command) bool {
	v, _ := cmd.Root().PersistentFlags().GetBool("json")
	return v
}

func isColor(cmd *cobra.Command) bool {
	nc, _ := cmd.Root().PersistentFlags().GetBool("no-color")
	if nc {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return true
}
