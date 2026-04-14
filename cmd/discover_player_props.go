package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mgm702/odds-api-cli/internal/discovery"
	"github.com/mgm702/odds-api-cli/internal/output"
	"github.com/spf13/cobra"
)

var discoverPlayerPropsCmd = &cobra.Command{
	Use:   "player-props <sport>",
	Short: "Discover available player prop market keys for a sport",
	Args:  cobra.ExactArgs(1),
	Example: `  odds discover player-props basketball_nba
  odds discover player-props basketball_nba --regions us --sample-size 3
  odds discover player-props basketball_nba --deep-probe --json`,
	RunE: runDiscoverPlayerProps,
}

func init() {
	discoverPlayerPropsCmd.Flags().String("regions", "us", "Comma-delimited regions")
	discoverPlayerPropsCmd.Flags().String("bookmakers", "", "Comma-separated bookmaker keys")
	discoverPlayerPropsCmd.Flags().StringSlice("event-id", nil, "Specific event IDs to inspect (repeatable)")
	discoverPlayerPropsCmd.Flags().Int("sample-size", 5, "Max number of events to sample when event IDs are not provided")
	discoverPlayerPropsCmd.Flags().Int("max-credits", 25, "Approximate max request budget for discovery")
	discoverPlayerPropsCmd.Flags().Bool("deep-probe", false, "Probe event-odds endpoint with known player prop key candidates")
	discoverCmd.AddCommand(discoverPlayerPropsCmd)
}

func runDiscoverPlayerProps(cmd *cobra.Command, args []string) error {
	c, err := newRuntimeClient()
	if err != nil {
		return err
	}

	regions, _ := cmd.Flags().GetString("regions")
	bookmakers, _ := cmd.Flags().GetString("bookmakers")
	eventIDs, _ := cmd.Flags().GetStringSlice("event-id")
	sampleSize, _ := cmd.Flags().GetInt("sample-size")
	maxCredits, _ := cmd.Flags().GetInt("max-credits")
	deepProbe, _ := cmd.Flags().GetBool("deep-probe")

	result, err := discovery.DiscoverPlayerProps(context.Background(), c, discovery.PlayerPropOptions{
		Sport:      args[0],
		Regions:    strings.TrimSpace(regions),
		Bookmakers: strings.TrimSpace(bookmakers),
		EventIDs:   eventIDs,
		SampleSize: sampleSize,
		MaxCredits: maxCredits,
		DeepProbe:  deepProbe,
	})
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(result)
	}

	fmt.Fprintf(os.Stdout, "Sport: %s\n", result.Sport)
	fmt.Fprintf(os.Stdout, "Regions: %s\n", result.Regions)
	fmt.Fprintf(os.Stdout, "Sampled events: %d\n", result.SampledEvents)
	fmt.Fprintf(os.Stdout, "Requests used: %d\n\n", result.RequestsUsed)

	if len(result.Markets) == 0 {
		fmt.Fprintln(os.Stdout, "No player prop market keys discovered.")
		return nil
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "MARKET KEY\tCONFIDENCE\tEVENTS\tBOOKMAKERS\tOBSERVED\tSOURCE")
	for _, m := range result.Markets {
		fmt.Fprintf(tw, "%s\t%.2f\t%d\t%d\t%d\t%s\n", m.Key, m.Confidence, m.EventCount, m.BookmakerCount, m.Occurrences, m.Source)
	}
	return tw.Flush()
}
