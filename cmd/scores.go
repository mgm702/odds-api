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

var scoresCmd = &cobra.Command{
	Use:   "scores <sport>",
	Short: "Get scores for a sport",
	Args:  cobra.ExactArgs(1),
	Example: `  odds scores basketball_nba
  odds scores basketball_nba --days-from 1
  odds scores basketball_nba --json`,
	RunE: runScores,
}

func init() {
	scoresCmd.Flags().Int("days-from", 0, "Include completed games from past N days (1-3)")
	rootCmd.AddCommand(scoresCmd)
}

func runScores(cmd *cobra.Command, args []string) error {
	c, err := newRuntimeClient()
	if err != nil {
		return err
	}

	params := url.Values{}
	if dateFormat != "iso" {
		params.Set("dateFormat", dateFormat)
	}
	if v, _ := cmd.Flags().GetInt("days-from"); v > 0 {
		params.Set("daysFrom", fmt.Sprintf("%d", v))
	}

	path := fmt.Sprintf("/v4/sports/%s/scores", args[0])
	resp, err := c.Get(context.Background(), path, params)
	if err != nil {
		return err
	}

	scores, err := client.Decode[[]model.ScoreEvent](resp)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(scores)
	}

	output.NewTableWriter(os.Stdout, useColor()).WriteScores(scores)
	return nil
}
