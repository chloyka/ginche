package ginche

import (
	"sync"
	"time"
)

// Cache is a thread-safe in-memory cache.
// It is safe to use concurrently.
// It will automatically cleanup expired items.
type Cache struct {
	items           *sync.Map
	ttl             time.Duration
	cleanupInterval time.Duration
}

// CacheConfig is used to configure a cache.
// If CleanupInterval or TTL is nil, it will default to 1 minute.
type CacheConfig struct {
	TTL             *time.Duration
	CleanupInterval *time.Duration
}

// Item is an item in the cache.
// It contains the value and the time it expires.
type Item struct {
	value     interface{}
	expiresAt time.Time
}

// ItemConfig is used to configure an item.
// If TTL is nil, it will use the cache's default TTL.
type ItemConfig struct {
	TTL *time.Duration
}

// NewCache creates a new cache with the given ttl and cleanupInterval.
// If cleanupInterval or ttl is nil, it will default to 1 minute.
func NewCache(config ...CacheConfig) *Cache {
	in := time.Minute
	cleanupInterval := &in
	ttl := &in

	if config != nil && config[0].CleanupInterval != nil {
		cleanupInterval = config[0].CleanupInterval
	}
	if config != nil && ttl != config[0].TTL {
		ttl = config[0].TTL
	}
	c := &Cache{
		items:           &sync.Map{},
		ttl:             *ttl,
		cleanupInterval: *cleanupInterval,
	}

	go c.cleanup()
	return c
}

// Set adds an item to the cache with the given key and value.
// If config is not nil, it will use the TTL from the config.
// Otherwise, it will use the cache's default TTL.
func (c *Cache) Set(key *string, value interface{}, config ...*ItemConfig) {
	var expiresAt time.Time
	if config != nil {
		if config[0].TTL != nil {
			expiresAt = time.Now().Add(*config[0].TTL)
		}
	} else {
		expiresAt = time.Now().Add(c.ttl)
	}

	c.items.Store(*key, &Item{value: value, expiresAt: expiresAt})
}

// Get returns the value of the item with the given key.
// If the item does not exist or has expired, it will return nil and false.
func (c *Cache) Get(key string) (interface{}, bool) {
	itemInterface, ok := c.items.Load(key)
	if !ok {
		return nil, ok
	}
	item := itemInterface.(*Item)
	if time.Now().After(item.expiresAt) {
		c.items.Delete(key)
		return nil, false
	}
	return item.value, true
}

// FlushAll deletes all items from the cache.
func (c *Cache) FlushAll() {
	c.items = &sync.Map{}
}

// cleanup deletes all expired items from the cache.
// It is called every cleanupInterval.
func (c *Cache) cleanup() {
	for {
		time.Sleep(c.cleanupInterval)
		c.items.Range(func(k, v interface{}) bool {
			if time.Now().After(v.(*Item).expiresAt) {
				c.items.Delete(k)
			}
			return true
		})
	}
}

// String Converts a string to a string pointer
func String(str string) *string {
	return &str
}

// TODO: Implement adapter interface for external storages
// TODO: Implement Redis storage
// TODO: Implement Memcached storage
