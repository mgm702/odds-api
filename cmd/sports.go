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

var sportsCmd = &cobra.Command{
	Use:   "sports",
	Short: "List available sports",
	Example: `  odds sports
  odds sports --all
  odds sports --interactive
  odds sports --json`,
	RunE: runSports,
}

func init() {
	sportsCmd.Flags().Bool("all", false, "Include out-of-season sports")
	sportsCmd.Flags().Bool("interactive", false, "Launch interactive sports browser")
	rootCmd.AddCommand(sportsCmd)
}

func runSports(cmd *cobra.Command, args []string) error {
	c, err := newRuntimeClient()
	if err != nil {
		return err
	}

	params := url.Values{}
	if all, _ := cmd.Flags().GetBool("all"); all {
		params.Set("all", "true")
	}

	resp, err := c.Get(context.Background(), "/v4/sports", params)
	if err != nil {
		return err
	}

	sports, err := client.Decode[[]model.Sport](resp)
	if err != nil {
		return err
	}

	interactive, _ := cmd.Flags().GetBool("interactive")
	if interactive && !jsonOutput {
		m := tui.NewSportsModel(sports)
		p := tea.NewProgram(m, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		if sm, ok := finalModel.(tui.SportsModel); ok {
			sm.WriteSelected(os.Stdout)
		}
		return nil
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(sports)
	}

	output.NewTableWriter(os.Stdout, useColor()).WriteSports(sports)
	return nil
}
