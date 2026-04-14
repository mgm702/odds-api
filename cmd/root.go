package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mgm702/odds-api-cli/cmd/historical"
	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	apiKey       string
	jsonOutput   bool
	quiet        bool
	verbose      bool
	noColor      bool
	dateFormat   string
	oddsFormat   string
	cacheEnabled bool
	cacheMode    string
	cacheTTL     time.Duration
	cacheDir     string
	version      = "dev"
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
	Version:           version,
	PersistentPreRunE: configureClientDefaults,
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
	rootCmd.PersistentFlags().StringVar(&oddsFormat, "odds-format", "", "Odds format: decimal or american")
	rootCmd.PersistentFlags().BoolVar(&cacheEnabled, "cache", true, "Enable local response caching")
	rootCmd.PersistentFlags().StringVar(&cacheMode, "cache-mode", string(client.CacheModeSmart), "Cache mode: smart, off, refresh")
	rootCmd.PersistentFlags().DurationVar(&cacheTTL, "cache-ttl", time.Minute, "Cache TTL duration (e.g. 60s, 5m)")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "Cache directory (overrides ODDS_CACHE_DIR)")

	rootCmd.AddCommand(historical.HistoricalCmd)
}

func configureClientDefaults(cmd *cobra.Command, args []string) error {
	cfg, err := currentCacheConfig()
	if err != nil {
		return err
	}
	client.SetDefaultCacheConfig(cfg)
	return nil
}

func currentCacheConfig() (client.CacheConfig, error) {
	mode := client.CacheMode(strings.ToLower(strings.TrimSpace(cacheMode)))
	switch mode {
	case client.CacheModeSmart, client.CacheModeOff, client.CacheModeRefresh:
	default:
		return client.CacheConfig{}, fmt.Errorf("invalid --cache-mode %q: must be smart, off, or refresh", cacheMode)
	}
	return client.CacheConfig{
		Enabled: cacheEnabled,
		Mode:    mode,
		TTL:     cacheTTL,
		Dir:     strings.TrimSpace(cacheDir),
	}, nil
}

func newRuntimeClient() (*client.Client, error) {
	key, err := getAPIKey()
	if err != nil {
		return nil, err
	}
	c := client.New(key)
	c.Verbose = verbose
	cfg, cfgErr := currentCacheConfig()
	if cfgErr != nil {
		return nil, cfgErr
	}
	c.SetCacheConfig(cfg)
	return c, nil
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
