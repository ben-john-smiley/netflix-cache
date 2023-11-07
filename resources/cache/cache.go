package cache

import (
	"github.com/allegro/bigcache"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	MaxEntrySize = 20971520      // 20 MiB entry size (used in initial allocation, can grow)
	MaxCacheSize = 2000          // ~ 2 GB max cache size, shouldn't need to be larger given constraints
	EntryLife    = 1 * time.Hour // mostly arbitrary value, should suit bottom-N implementation
)

type BigCache struct {
	NetflixOrgs *bigcache.BigCache
}

// New initializes a cache and returns it
func New() (*BigCache, error) {
	cache, err := bigcache.NewBigCache(bigcache.Config{
		Shards:           1024,
		LifeWindow:       EntryLife,
		MaxEntrySize:     MaxEntrySize,
		HardMaxCacheSize: MaxCacheSize,
	})

	// Not panicking here, if the cache fails to initialize the service can still operate
	// ideally would emit a metric here to be monitored external to the service to alert operators that the
	// service is running at reduced capacity
	if err != nil {
		log.Error("Failed to initialize cache", err)
		return nil, err
	}
	return &BigCache{
		NetflixOrgs: cache,
	}, nil
}

// Read attempts to read a given key from the cache
func (c *BigCache) Read(key string) (item []byte, ok bool) {
	item, err := c.NetflixOrgs.Get(key)
	if err != nil {
		return nil, false
	}
	return item, true
}

// Write writes a value into the cache, no collision checks, just clobbering what is already there
func (c *BigCache) Write(key string, value []byte) error {
	return c.NetflixOrgs.Set(key, value)
}
