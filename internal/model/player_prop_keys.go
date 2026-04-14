package model

var knownPlayerPropMarketKeys = []string{
	"player_points",
	"player_rebounds",
	"player_assists",
	"player_threes",
	"player_blocks",
	"player_steals",
	"player_turnovers",
	"player_points_rebounds_assists",
	"player_points_rebounds",
	"player_points_assists",
	"player_rebounds_assists",
	"player_first_basket",
	"player_double_double",
	"player_triple_double",
	"player_shots_on_goal",
	"player_pass_tds",
	"player_pass_yds",
	"player_rush_yds",
	"player_reception_yds",
}

func KnownPlayerPropMarketKeys() []string {
	keys := make([]string, len(knownPlayerPropMarketKeys))
	copy(keys, knownPlayerPropMarketKeys)
	return keys
}
