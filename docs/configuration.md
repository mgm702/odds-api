# Configuration

## API Key

- `ODDS_API_KEY` environment variable, or
- `--api-key` flag

## Verbose Logging

Use `--verbose` to print request URL (with redacted key), status, elapsed time, and quota metadata.

## Caching

- `--cache` (default: true)
- `--cache-mode smart|off|refresh` (default: smart)
- `--cache-ttl` (default: 60s)
- `--cache-dir` (overrides `ODDS_CACHE_DIR`)

Directory resolution when `--cache-dir` is unset:

- `ODDS_CACHE_DIR` if provided
- `${XDG_CACHE_HOME}/odds-api-cli` if provided
- `~/.cache/odds-api-cli` fallback

## Output

- `--json` for machine-readable output
- table output by default
- watch and selected commands support interactive TUI paths

## Date Format

- `--date-format iso` (default)
- `--date-format unix`
