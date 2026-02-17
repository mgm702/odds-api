package cmd

import (
	"context"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestParticipantsCommand_Decode(t *testing.T) {
	srv, _ := setupTestServer(t, fixtureHandler("../testdata/participants.json"))

	c := client.New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports/basketball_nba/participants", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	participants, err := client.Decode[[]model.Participant](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(participants) != 3 {
		t.Fatalf("expected 3 participants, got %d", len(participants))
	}
	if !participants[0].IsActive {
		t.Error("expected first participant to be active")
	}
	if participants[2].IsActive {
		t.Error("expected third participant to be inactive")
	}
}
