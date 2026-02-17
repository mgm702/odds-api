package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
	"github.com/mgm702/odds-api-cli/internal/tui"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch <sport>",
	Short: "Live polling TUI for odds or scores",
	Args:  cobra.ExactArgs(1),
	Example: `  odds watch basketball_nba --regions us
  odds watch basketball_nba --regions us --markets h2h,spreads --interval 30
  odds watch basketball_nba --mode scores`,
	RunE: runWatch,
}

func init() {
	watchCmd.Flags().String("mode", "odds", "Watch mode: odds or scores")
	watchCmd.Flags().Int("interval", 60, "Poll interval in seconds (min: 30)")
	watchCmd.Flags().String("regions", "", "Comma-delimited regions (required for odds mode)")
	watchCmd.Flags().String("markets", "", "Comma-delimited markets")
	watchCmd.Flags().String("event-ids", "", "Comma-separated event IDs")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	key, err := getAPIKey()
	if err != nil {
		return err
	}

	mode, _ := cmd.Flags().GetString("mode")
	if mode != "odds" && mode != "scores" {
		return fmt.Errorf("invalid mode %q: must be odds or scores", mode)
	}

	interval, _ := cmd.Flags().GetInt("interval")
	if interval < 30 {
		interval = 30
	}

	regions, _ := cmd.Flags().GetString("regions")
	if mode == "odds" && regions == "" {
		return fmt.Errorf("--regions is required for odds mode")
	}

	c := client.New(key)
	c.Verbose = verbose

	sport := args[0]
	markets, _ := cmd.Flags().GetString("markets")
	eventIDs, _ := cmd.Flags().GetString("event-ids")

	fetch := buildFetchFunc(c, mode, sport, regions, markets, eventIDs)

	m := tui.NewWatchModel(mode, time.Duration(interval)*time.Second, fetch)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func buildFetchFunc(c *client.Client, mode, sport, regions, markets, eventIDs string) tui.FetchFunc {
	return func(ctx context.Context) (tui.WatchData, error) {
		if mode == "scores" {
			params := url.Values{}
			path := fmt.Sprintf("/v4/sports/%s/scores", sport)
			resp, err := c.Get(ctx, path, params)
			if err != nil {
				return tui.WatchData{}, err
			}
			scores, err := client.Decode[[]model.ScoreEvent](resp)
			if err != nil {
				return tui.WatchData{}, err
			}
			return tui.WatchData{ScoreEvents: scores, Quota: resp.Quota}, nil
		}

		params := url.Values{}
		params.Set("regions", regions)
		if markets != "" {
			params.Set("markets", markets)
		}
		if eventIDs != "" {
			params.Set("eventIds", eventIDs)
		}

		path := fmt.Sprintf("/v4/sports/%s/odds", sport)
		resp, err := c.Get(ctx, path, params)
		if err != nil {
			return tui.WatchData{}, err
		}
		events, err := client.Decode[[]model.OddsEvent](resp)
		if err != nil {
			return tui.WatchData{}, err
		}
		return tui.WatchData{OddsEvents: events, Quota: resp.Quota}, nil
	}
}
