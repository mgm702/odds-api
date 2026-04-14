package discovery

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/mgm702/odds-api-cli/internal/client"
	"github.com/mgm702/odds-api-cli/internal/model"
)

type PlayerPropOptions struct {
	Sport      string
	Regions    string
	Bookmakers string
	EventIDs   []string
	SampleSize int
	MaxCredits int
	DeepProbe  bool
}

type PlayerPropMarket struct {
	Key            string  `json:"key"`
	EventCount     int     `json:"event_count"`
	BookmakerCount int     `json:"bookmaker_count"`
	Occurrences    int     `json:"occurrences"`
	Confidence     float64 `json:"confidence"`
	Source         string  `json:"source"`
}

type PlayerPropResult struct {
	Sport         string             `json:"sport"`
	Regions       string             `json:"regions"`
	SampledEvents int                `json:"sampled_events"`
	RequestsUsed  int                `json:"requests_used"`
	GeneratedAt   string             `json:"generated_at"`
	Markets       []PlayerPropMarket `json:"markets"`
}

type marketAgg struct {
	key        string
	events     map[string]struct{}
	bookmakers map[string]struct{}
	count      int
	source     string
}

func DiscoverPlayerProps(ctx context.Context, c *client.Client, opts PlayerPropOptions) (PlayerPropResult, error) {
	if strings.TrimSpace(opts.Sport) == "" {
		return PlayerPropResult{}, fmt.Errorf("sport is required")
	}
	if opts.SampleSize <= 0 {
		opts.SampleSize = 5
	}
	if opts.MaxCredits <= 0 {
		opts.MaxCredits = 25
	}

	eventIDs := opts.EventIDs
	requestsUsed := 0
	if len(eventIDs) == 0 {
		if requestsUsed >= opts.MaxCredits {
			return PlayerPropResult{}, fmt.Errorf("max credits budget reached before discovery")
		}
		ids, err := fetchEventIDs(ctx, c, opts.Sport)
		if err != nil {
			return PlayerPropResult{}, err
		}
		requestsUsed++
		if len(ids) > opts.SampleSize {
			ids = ids[:opts.SampleSize]
		}
		eventIDs = ids
	}

	agg := map[string]*marketAgg{}
	for _, eventID := range eventIDs {
		if requestsUsed >= opts.MaxCredits {
			break
		}
		if err := fetchEventMarkets(ctx, c, opts, eventID, agg); err != nil {
			continue
		}
		requestsUsed++

		if opts.DeepProbe && requestsUsed < opts.MaxCredits {
			_ = deepProbeKnownKeys(ctx, c, opts, eventID, agg)
			requestsUsed++
		}
	}

	markets := flattenAgg(agg)
	return PlayerPropResult{
		Sport:         opts.Sport,
		Regions:       opts.Regions,
		SampledEvents: len(eventIDs),
		RequestsUsed:  requestsUsed,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Markets:       markets,
	}, nil
}

func fetchEventIDs(ctx context.Context, c *client.Client, sport string) ([]string, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("/v4/sports/%s/events", sport), nil)
	if err != nil {
		return nil, err
	}
	events, err := client.Decode[[]model.Event](resp)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(events))
	for _, e := range events {
		if strings.TrimSpace(e.ID) != "" {
			ids = append(ids, e.ID)
		}
	}
	return ids, nil
}

func fetchEventMarkets(ctx context.Context, c *client.Client, opts PlayerPropOptions, eventID string, agg map[string]*marketAgg) error {
	params := url.Values{}
	if strings.TrimSpace(opts.Bookmakers) != "" {
		params.Set("bookmakers", opts.Bookmakers)
	}
	resp, err := c.Get(ctx, fmt.Sprintf("/v4/sports/%s/events/%s/markets", opts.Sport, eventID), params)
	if err != nil {
		return err
	}
	data, err := client.Decode[model.EventMarkets](resp)
	if err != nil {
		return err
	}

	for _, bm := range data.Bookmakers {
		for _, m := range bm.Markets {
			if !IsPlayerPropKey(m.Key) {
				continue
			}
			addAgg(agg, m.Key, eventID, bm.Key, "markets")
		}
	}
	return nil
}

func deepProbeKnownKeys(ctx context.Context, c *client.Client, opts PlayerPropOptions, eventID string, agg map[string]*marketAgg) error {
	keys := model.KnownPlayerPropMarketKeys()
	if len(keys) == 0 {
		return nil
	}
	params := url.Values{}
	params.Set("regions", defaultRegions(opts.Regions))
	params.Set("markets", strings.Join(keys, ","))
	if strings.TrimSpace(opts.Bookmakers) != "" {
		params.Set("bookmakers", opts.Bookmakers)
	}
	resp, err := c.Get(ctx, fmt.Sprintf("/v4/sports/%s/events/%s/odds", opts.Sport, eventID), params)
	if err != nil {
		return err
	}
	odds, err := client.Decode[model.OddsEvent](resp)
	if err != nil {
		return err
	}

	for _, bm := range odds.Bookmakers {
		for _, m := range bm.Markets {
			if !IsPlayerPropKey(m.Key) {
				continue
			}
			addAgg(agg, m.Key, eventID, bm.Key, "event-odds")
		}
	}
	return nil
}

func defaultRegions(v string) string {
	if strings.TrimSpace(v) == "" {
		return "us"
	}
	return v
}

func addAgg(agg map[string]*marketAgg, key, eventID, bookmaker, source string) {
	item, ok := agg[key]
	if !ok {
		item = &marketAgg{
			key:        key,
			events:     map[string]struct{}{},
			bookmakers: map[string]struct{}{},
			source:     source,
		}
		agg[key] = item
	}
	item.count++
	if eventID != "" {
		item.events[eventID] = struct{}{}
	}
	if bookmaker != "" {
		item.bookmakers[bookmaker] = struct{}{}
	}
}

func flattenAgg(agg map[string]*marketAgg) []PlayerPropMarket {
	if len(agg) == 0 {
		return nil
	}
	maxCount := 0
	for _, item := range agg {
		if item.count > maxCount {
			maxCount = item.count
		}
	}
	out := make([]PlayerPropMarket, 0, len(agg))
	for _, item := range agg {
		confidence := 1.0
		if maxCount > 0 {
			confidence = float64(item.count) / float64(maxCount)
		}
		out = append(out, PlayerPropMarket{
			Key:            item.key,
			EventCount:     len(item.events),
			BookmakerCount: len(item.bookmakers),
			Occurrences:    item.count,
			Confidence:     confidence,
			Source:         item.source,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Confidence == out[j].Confidence {
			return out[i].Key < out[j].Key
		}
		return out[i].Confidence > out[j].Confidence
	})
	return out
}

func IsPlayerPropKey(key string) bool {
	k := strings.ToLower(strings.TrimSpace(key))
	if k == "" {
		return false
	}
	for _, known := range model.KnownPlayerPropMarketKeys() {
		if k == known {
			return true
		}
	}
	if strings.Contains(k, "player_") || strings.Contains(k, "_player") || strings.HasPrefix(k, "player") {
		return true
	}
	return false
}
