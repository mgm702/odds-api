package client

import (
	"sync"
	"time"

	"github.com/mgm702/odds-api-cli/internal/cache"
)

type CacheMode string

const (
	CacheModeSmart   CacheMode = "smart"
	CacheModeOff     CacheMode = "off"
	CacheModeRefresh CacheMode = "refresh"
)

type CacheConfig struct {
	Enabled bool
	Mode    CacheMode
	TTL     time.Duration
	Dir     string
}

var (
	cacheDefaultsMu sync.RWMutex
	cacheDefaults   = CacheConfig{
		Enabled: false,
		Mode:    CacheModeSmart,
		TTL:     time.Minute,
	}
)

func SetDefaultCacheConfig(cfg CacheConfig) {
	cacheDefaultsMu.Lock()
	defer cacheDefaultsMu.Unlock()
	cacheDefaults = cfg
}

func DefaultCacheConfig() CacheConfig {
	cacheDefaultsMu.RLock()
	defer cacheDefaultsMu.RUnlock()
	return cacheDefaults
}

func (c *Client) SetCacheConfig(cfg CacheConfig) {
	c.CacheConfig = cfg
	c.cacheStore = nil
}

func (c *Client) cacheStoreOrInit() (*cache.Store, error) {
	if c.cacheStore != nil {
		return c.cacheStore, nil
	}
	dir := c.CacheConfig.Dir
	if dir == "" {
		resolved, err := cache.ResolveDir()
		if err != nil {
			return nil, err
		}
		dir = resolved
	}
	s := cache.NewStore(dir)
	c.cacheStore = s
	return s, nil
}
