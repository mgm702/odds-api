package discovery

import "testing"

func TestIsPlayerPropKey(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"player_points", true},
		{"PLAYER_REBOUNDS", true},
		{"h2h", false},
		{"totals", false},
		{"anytime_player_goal_scorer", true},
	}

	for _, tc := range cases {
		got := IsPlayerPropKey(tc.key)
		if got != tc.want {
			t.Fatalf("IsPlayerPropKey(%q)=%v want %v", tc.key, got, tc.want)
		}
	}
}

func TestFlattenAgg(t *testing.T) {
	agg := map[string]*marketAgg{}
	addAgg(agg, "player_points", "evt1", "draftkings", "markets")
	addAgg(agg, "player_points", "evt2", "fanduel", "markets")
	addAgg(agg, "player_assists", "evt1", "draftkings", "markets")

	out := flattenAgg(agg)
	if len(out) != 2 {
		t.Fatalf("expected 2 markets, got %d", len(out))
	}
	if out[0].Key != "player_points" {
		t.Fatalf("expected player_points first, got %s", out[0].Key)
	}
	if out[0].EventCount != 2 {
		t.Fatalf("expected event count 2, got %d", out[0].EventCount)
	}
}
