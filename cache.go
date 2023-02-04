package ginche

import (
	"sync"
	"time"
)

type Cache struct {
	items           *sync.Map
	ttl             time.Duration
	cleanupInterval time.Duration
}

type CacheConfig struct {
	TTL             *time.Duration
	CleanupInterval *time.Duration
}

type Item struct {
	value     interface{}
	expiresAt time.Time
}

type ItemConfig struct {
	TTL *time.Duration
}

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

func (c *Cache) FlushAll() {
	c.items = &sync.Map{}
}

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
