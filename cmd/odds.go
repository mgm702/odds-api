package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
	"github.com/mgm702/odds-api-cli/internal/output"
	"github.com/mgm702/odds-api-cli/internal/tui"
	"github.com/spf13/cobra"
)

var linesCmd = &cobra.Command{
	Use:   "lines <sport>",
	Short: "Get lines for a sport",
	Args:  cobra.ExactArgs(1),
	Example: `  odds lines basketball_nba --regions us
  odds lines basketball_nba --regions us,uk --markets h2h,spreads
  odds lines upcoming --regions us`,
	RunE: runOdds,
}

func init() {
	linesCmd.Flags().String("regions", "", "Comma-delimited regions: us, us2, uk, au, eu (required)")
	linesCmd.Flags().String("markets", "", "Comma-delimited markets: h2h, spreads, totals, outrights")
	linesCmd.Flags().String("event-ids", "", "Comma-separated event IDs")
	linesCmd.Flags().String("bookmakers", "", "Comma-separated bookmaker keys")
	linesCmd.Flags().String("from", "", "Filter events starting at/after (ISO 8601)")
	linesCmd.Flags().String("to", "", "Filter events starting at/before (ISO 8601)")
	_ = linesCmd.MarkFlagRequired("regions")
	rootCmd.AddCommand(linesCmd)
}

func runOdds(cmd *cobra.Command, args []string) error {
	c, err := newRuntimeClient()
	if err != nil {
		return err
	}

	params := url.Values{}
	regions, _ := cmd.Flags().GetString("regions")
	params.Set("regions", regions)

	if v, _ := cmd.Flags().GetString("markets"); v != "" {
		params.Set("markets", v)
	}
	if oddsFormat != "" {
		params.Set("oddsFormat", oddsFormat)
	}
	if dateFormat != "iso" {
		params.Set("dateFormat", dateFormat)
	}
	if v, _ := cmd.Flags().GetString("event-ids"); v != "" {
		params.Set("eventIds", v)
	}
	if v, _ := cmd.Flags().GetString("bookmakers"); v != "" {
		params.Set("bookmakers", v)
	}
	if v, _ := cmd.Flags().GetString("from"); v != "" {
		params.Set("commenceTimeFrom", v)
	}
	if v, _ := cmd.Flags().GetString("to"); v != "" {
		params.Set("commenceTimeTo", v)
	}

	path := fmt.Sprintf("/v4/sports/%s/odds", args[0])
	resp, err := c.Get(context.Background(), path, params)
	if err != nil {
		return err
	}

	events, err := client.Decode[[]model.OddsEvent](resp)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(events)
	}

	if !quiet && len(events) > 0 {
		m := tui.NewOddsModel(events, oddsFormat)
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err := p.Run()
		return err
	}

	tw := output.NewTableWriter(os.Stdout, useColor())
	tw.OddsFormat = oddsFormat
	tw.WriteOdds(events)
	return nil
}
