# odds-api-cli

A Go CLI for The Odds API v4 with table output, JSON mode, historical queries, and a live watch TUI.

## Features

- Sports, events, odds, event-odds, markets, participants, scores, and credits commands
- Historical snapshots for odds/events/event-odds
- `watch` mode for live polling odds or scores
- Local response caching with configurable TTL/mode
- Player prop key discovery workflow
- Verbose API request logging via `--verbose`

## Install

### Go install

```bash
go install github.com/mgm702/odds-api-cli@latest
```

### Homebrew (after release automation is enabled)

```bash
brew tap mgm702/odds
brew install odds
```

## Configuration

Set your API key:

```bash
export ODDS_API_KEY="your_key_here"
```

Global flags:

- `--api-key`
- `--json`
- `--verbose`
- `--cache`
- `--cache-mode` (`smart|off|refresh`)
- `--cache-ttl` (for example `60s`, `5m`)
- `--cache-dir`
- `--no-color`
- `--date-format`
- `--odds-format` (`decimal|american`)

## Quickstart

```bash
odds sports
odds events basketball_nba
odds odds basketball_nba --regions us --markets h2h,spreads
odds scores basketball_nba --days-from 1
odds credits
odds watch basketball_nba --regions us
odds discover player-props basketball_nba --regions us
```

## Command Overview

- `odds sports`
- `odds events <sport>`
- `odds odds <sport> --regions ...`
- `odds event-odds <sport> <event-id> --regions ... --markets ...`
- `odds markets <sport> <event-id>`
- `odds participants <sport>`
- `odds scores <sport>`
- `odds credits`
- `odds watch <sport>`
- `odds discover player-props <sport>`
- `odds historical odds <sport> --date ... --regions ...`
- `odds historical events <sport> --date ...`
- `odds historical event-odds <sport> <event-id> --date ... --regions ... --markets ...`

## Caching and Credit Safety

- Cache is enabled by default in CLI execution.
- Use `--cache-mode off` to disable cache reads/writes.
- Use `--cache-mode refresh` to force network fetches while refreshing cache entries.
- Watch mode disables cache reads by default; pass `--use-cache` to enable.

## Docs

- `docs/configuration.md`
- `docs/commands/`
