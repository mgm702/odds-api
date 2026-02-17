package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
	"github.com/mgm702/odds-api-cli/internal/output"
	"github.com/spf13/cobra"
)

var participantsCmd = &cobra.Command{
	Use:   "participants <sport>",
	Short: "List teams/players for a sport",
	Args:  cobra.ExactArgs(1),
	Example: `  odds participants basketball_nba
  odds participants basketball_nba --json`,
	RunE: runParticipants,
}

func init() {
	rootCmd.AddCommand(participantsCmd)
}

func runParticipants(cmd *cobra.Command, args []string) error {
	key, err := getAPIKey()
	if err != nil {
		return err
	}

	c := client.New(key)
	c.Verbose = verbose

	path := fmt.Sprintf("/v4/sports/%s/participants", args[0])
	resp, err := c.Get(context.Background(), path, nil)
	if err != nil {
		return err
	}

	participants, err := client.Decode[[]model.Participant](resp)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.NewJSONWriter(os.Stdout).Write(participants)
	}

	output.NewTableWriter(os.Stdout, useColor()).WriteParticipants(participants)
	return nil
}
