package cache

import (
	"sync"
	"time"
)

type Cache struct {
	items           *sync.Map
	ttl             time.Duration
	cleanupInterval time.Duration
}

type Item struct {
	value     interface{}
	expiresAt time.Time
}

type ItemConfig struct {
	TTL *time.Duration
}

func NewCache(ttl time.Duration, cleanupInterval *time.Duration) *Cache {
	if cleanupInterval == nil {
		in := time.Minute
		cleanupInterval = &in
	}
	c := &Cache{
		items:           &sync.Map{},
		ttl:             ttl,
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

func (c *Cache) cleanup() {
	for {
		time.Sleep(time.Minute)
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
