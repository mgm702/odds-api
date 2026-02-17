package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/mgm702/odds-api-cli/internal/model"
)

type TableWriter struct {
	Out   io.Writer
	Color bool
}

func NewTableWriter(out io.Writer, color bool) *TableWriter {
	return &TableWriter{Out: out, Color: color}
}

func (t *TableWriter) WriteSports(sports []model.Sport) {
	w := tabwriter.NewWriter(t.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tTITLE\tGROUP\tACTIVE")
	for _, s := range sports {
		active := "no"
		if s.Active {
			active = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Key, s.Title, s.Group, active)
	}
	w.Flush()
}

func (t *TableWriter) WriteEvents(events []model.Event) {
	w := tabwriter.NewWriter(t.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tHOME\tAWAY\tSTART TIME")
	for _, e := range events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.ID, e.HomeTeam, e.AwayTeam, e.CommenceTime)
	}
	w.Flush()
}

func (t *TableWriter) WriteOdds(events []model.OddsEvent) {
	for i, e := range events {
		if i > 0 {
			fmt.Fprintln(t.Out)
		}
		fmt.Fprintf(t.Out, "%s vs %s (%s)\n", e.HomeTeam, e.AwayTeam, e.CommenceTime)
		for _, b := range e.Bookmakers {
			fmt.Fprintf(t.Out, "  %s\n", b.Title)
			for _, m := range b.Markets {
				fmt.Fprintf(t.Out, "    [%s]\n", m.Key)
				w := tabwriter.NewWriter(t.Out, 0, 0, 2, ' ', 0)
				for _, o := range m.Outcomes {
					if o.Point != nil {
						fmt.Fprintf(w, "      %s\t%+.1f\t%.2f\n", o.Name, *o.Point, o.Price)
					} else {
						fmt.Fprintf(w, "      %s\t\t%.2f\n", o.Name, o.Price)
					}
				}
				w.Flush()
			}
		}
	}
}

func (t *TableWriter) WriteScores(events []model.ScoreEvent) {
	w := tabwriter.NewWriter(t.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOME\tAWAY\tSCORE\tSTATUS\tUPDATED")
	for _, e := range events {
		score := "-"
		if e.Scores != nil && len(e.Scores) >= 2 {
			score = fmt.Sprintf("%s - %s", e.Scores[0].Score, e.Scores[1].Score)
		}
		status := "upcoming"
		if e.Completed {
			status = "completed"
		} else if e.Scores != nil {
			status = "live"
		}
		updated := "-"
		if e.LastUpdate != nil {
			updated = *e.LastUpdate
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.HomeTeam, e.AwayTeam, score, status, updated)
	}
	w.Flush()
}

func (t *TableWriter) WriteCredits(q model.QuotaInfo) {
	total := q.Used + q.Remaining
	pct := 0.0
	if total > 0 {
		pct = float64(q.Used) / float64(total) * 100
	}

	fmt.Fprintln(t.Out, "Credits Report")
	fmt.Fprintln(t.Out, strings.Repeat("\u2500", 30))
	w := tabwriter.NewWriter(t.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Used:\t%d\n", q.Used)
	fmt.Fprintf(w, "Remaining:\t%d\n", q.Remaining)
	fmt.Fprintf(w, "Last Cost:\t%d\n", q.LastCost)
	w.Flush()

	barLen := 20
	filled := int(pct / 100 * float64(barLen))
	if filled > barLen {
		filled = barLen
	}
	bar := strings.Repeat("\u2588", filled) + strings.Repeat("\u2591", barLen-filled)
	fmt.Fprintf(t.Out, "            %s  %.1f%% used\n", bar, pct)
}

func (t *TableWriter) WriteParticipants(participants []model.Participant) {
	w := tabwriter.NewWriter(t.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tACTIVE")
	for _, p := range participants {
		active := "no"
		if p.IsActive {
			active = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.ID, p.Name, active)
	}
	w.Flush()
}

func (t *TableWriter) WriteEventMarkets(em model.EventMarkets) {
	fmt.Fprintf(t.Out, "Event: %s\n\n", em.ID)
	for _, b := range em.Bookmakers {
		fmt.Fprintf(t.Out, "  %s (%s)\n", b.Title, b.Key)
		for _, m := range b.Markets {
			fmt.Fprintf(t.Out, "    - %s\n", m.Key)
		}
	}
}

func (t *TableWriter) WriteHistoricalOdds(h model.HistoricalResponse[[]model.OddsEvent]) {
	fmt.Fprintf(t.Out, "Snapshot: %s\n", h.Timestamp)
	if h.PreviousTimestamp != nil {
		fmt.Fprintf(t.Out, "Previous: %s\n", *h.PreviousTimestamp)
	}
	if h.NextTimestamp != nil {
		fmt.Fprintf(t.Out, "Next:     %s\n", *h.NextTimestamp)
	}
	fmt.Fprintln(t.Out)
	t.WriteOdds(h.Data)
}

func (t *TableWriter) WriteHistoricalEvents(h model.HistoricalResponse[[]model.Event]) {
	fmt.Fprintf(t.Out, "Snapshot: %s\n", h.Timestamp)
	if h.PreviousTimestamp != nil {
		fmt.Fprintf(t.Out, "Previous: %s\n", *h.PreviousTimestamp)
	}
	if h.NextTimestamp != nil {
		fmt.Fprintf(t.Out, "Next:     %s\n", *h.NextTimestamp)
	}
	fmt.Fprintln(t.Out)
	t.WriteEvents(h.Data)
}
