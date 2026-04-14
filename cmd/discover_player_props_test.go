package cmd

import "testing"

func TestDiscoverPlayerProps_Args(t *testing.T) {
	if err := discoverPlayerPropsCmd.Args(discoverPlayerPropsCmd, []string{}); err == nil {
		t.Fatal("expected arg error when sport missing")
	}
	if err := discoverPlayerPropsCmd.Args(discoverPlayerPropsCmd, []string{"basketball_nba"}); err != nil {
		t.Fatalf("unexpected arg error: %v", err)
	}
}

func TestDiscoverPlayerProps_DefaultFlags(t *testing.T) {
	regions, err := discoverPlayerPropsCmd.Flags().GetString("regions")
	if err != nil {
		t.Fatalf("failed to get regions flag: %v", err)
	}
	if regions != "us" {
		t.Fatalf("expected default regions us, got %s", regions)
	}

	sample, err := discoverPlayerPropsCmd.Flags().GetInt("sample-size")
	if err != nil {
		t.Fatalf("failed to get sample-size flag: %v", err)
	}
	if sample != 5 {
		t.Fatalf("expected sample-size 5, got %d", sample)
	}
}
