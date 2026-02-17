package model

type Sport struct {
	Key          string `json:"key"`
	Group        string `json:"group"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
	HasOutrights bool   `json:"has_outrights"`
}

type Event struct {
	ID            string `json:"id"`
	SportKey      string `json:"sport_key"`
	SportTitle    string `json:"sport_title"`
	CommenceTime  string `json:"commence_time"`
	HomeTeam      string `json:"home_team"`
	AwayTeam      string `json:"away_team"`
}

type OddsEvent struct {
	ID            string      `json:"id"`
	SportKey      string      `json:"sport_key"`
	SportTitle    string      `json:"sport_title"`
	CommenceTime  string      `json:"commence_time"`
	HomeTeam      string      `json:"home_team"`
	AwayTeam      string      `json:"away_team"`
	Bookmakers    []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	LastUpdate string   `json:"last_update"`
	Markets    []Market `json:"markets"`
}

type Market struct {
	Key        string    `json:"key"`
	LastUpdate string    `json:"last_update"`
	Outcomes   []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string   `json:"name"`
	Price float64  `json:"price"`
	Point *float64 `json:"point,omitempty"`
}

type Score struct {
	Name  string `json:"name"`
	Score string `json:"score"`
}

type ScoreEvent struct {
	ID           string  `json:"id"`
	SportKey     string  `json:"sport_key"`
	SportTitle   string  `json:"sport_title"`
	CommenceTime string  `json:"commence_time"`
	HomeTeam     string  `json:"home_team"`
	AwayTeam     string  `json:"away_team"`
	Completed    bool    `json:"completed"`
	Scores       []Score `json:"scores"`
	LastUpdate   *string `json:"last_update"`
}

type Participant struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type MarketInfo struct {
	Key string `json:"key"`
}

type BookmakerMarkets struct {
	Key     string       `json:"key"`
	Title   string       `json:"title"`
	Markets []MarketInfo `json:"markets"`
}

type EventMarkets struct {
	ID         string             `json:"id"`
	SportKey   string             `json:"sport_key"`
	Bookmakers []BookmakerMarkets `json:"bookmakers"`
}

type HistoricalResponse[T any] struct {
	Timestamp         string `json:"timestamp"`
	PreviousTimestamp *string `json:"previous_timestamp"`
	NextTimestamp     *string `json:"next_timestamp"`
	Data              T      `json:"data"`
}

type QuotaInfo struct {
	Remaining int `json:"remaining"`
	Used      int `json:"used"`
	LastCost  int `json:"last_cost"`
}
