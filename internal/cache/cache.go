package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mgm702/odds-api-cli/internal/model"
)

const keyVersion = "v1"

var ErrNotFound = errors.New("cache entry not found")

type Entry struct {
	StoredAt   time.Time       `json:"stored_at"`
	StatusCode int             `json:"status_code"`
	Quota      model.QuotaInfo `json:"quota"`
	Body       []byte          `json:"body"`
}

type Store struct {
	dir string
}

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func ResolveDir() (string, error) {
	if v := os.Getenv("ODDS_CACHE_DIR"); v != "" {
		return v, nil
	}
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return filepath.Join(v, "odds-api-cli"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".cache", "odds-api-cli"), nil
}

func RequestKey(method, path string, q url.Values) string {
	clean := url.Values{}
	for k, vals := range q {
		if strings.EqualFold(k, "apiKey") {
			continue
		}
		cp := append([]string(nil), vals...)
		sort.Strings(cp)
		for _, v := range cp {
			clean.Add(k, v)
		}
	}

	keys := make([]string, 0, len(clean))
	for k := range clean {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(keyVersion)
	b.WriteString("|")
	b.WriteString(strings.ToUpper(method))
	b.WriteString("|")
	b.WriteString(path)

	for _, k := range keys {
		vals := append([]string(nil), clean[k]...)
		sort.Strings(vals)
		for _, v := range vals {
			b.WriteString("|")
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(v)
		}
	}

	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func (s *Store) Get(key string, ttl time.Duration) (*Entry, error) {
	path := s.filePath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("read cache file: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("decode cache entry: %w", err)
	}

	if ttl > 0 && time.Since(entry.StoredAt) > ttl {
		_ = os.Remove(path)
		return nil, ErrNotFound
	}

	return &entry, nil
}

func (s *Store) Put(key string, entry Entry) error {
	if entry.StoredAt.IsZero() {
		entry.StoredAt = time.Now().UTC()
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	encoded, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("encode cache entry: %w", err)
	}

	path := s.filePath(key)
	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, encoded, 0o644); err != nil {
		return fmt.Errorf("write cache temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename cache temp file: %w", err)
	}

	return nil
}

func (s *Store) filePath(key string) string {
	return filepath.Join(s.dir, key+".json")
}
