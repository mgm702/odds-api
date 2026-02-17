package cmd

import (
	"fmt"
	"os"

	"github.com/mgm702/odds-api-cli/cmd/historical"
	"github.com/spf13/cobra"
)

var (
	apiKey     string
	jsonOutput bool
	quiet      bool
	verbose    bool
	noColor    bool
	dateFormat string
	version    = "dev"
)

func SetVersion(v string) {
	version = v
}

var rootCmd = &cobra.Command{
	Use:   "odds",
	Short: "CLI for The Odds API v4",
	Long:  "Query sports betting odds, scores, and historical data from The Odds API.",
	Example: `  odds sports
  odds events basketball_nba
  odds odds basketball_nba --regions us --markets h2h,spreads
  odds scores basketball_nba --days-from 1
  odds credits
  odds watch basketball_nba --regions us`,
	Version: version,
}

func Execute() error {
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (overrides ODDS_API_KEY env var)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show request URL, quota usage, timing")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().StringVar(&dateFormat, "date-format", "iso", "Date format: iso or unix")

	rootCmd.AddCommand(historical.HistoricalCmd)
}

func getAPIKey() (string, error) {
	if apiKey != "" {
		return apiKey, nil
	}
	key := os.Getenv("ODDS_API_KEY")
	if key == "" {
		return "", fmt.Errorf("API key required: set ODDS_API_KEY or use --api-key")
	}
	return key, nil
}

func useColor() bool {
	if noColor {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return true
}
