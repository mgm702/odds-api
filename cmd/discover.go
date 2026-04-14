package cmd

import "github.com/spf13/cobra"

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discovery utilities for market keys and metadata",
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
