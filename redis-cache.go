package ginche

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"strings"
	"time"
)

type RedisAdapter struct {
	conn          *redis.Client
	inMemoryCache *InMemoryCache
	pubsub        *redis.PubSub
	config        *CacheConfig
}

func NewRedisAdapter(redisConfig *redis.Options, config ...CacheConfig) (CacheAdapter, error) {
	redisClient := redis.NewClient(redisConfig)
	pubsub := redisClient.Subscribe(context.Background(), "cache_updates:*")
	var conf CacheConfig
	if config != nil {
		conf = config[0]
	} else {
		ttl := time.Minute * 5
		conf = CacheConfig{
			TTL:             &ttl,
			CleanupInterval: nil,
		}
	}
	inMemory := NewInMemoryCache(conf)
	cache := &RedisAdapter{
		conn:          redisClient,
		inMemoryCache: inMemory.(*InMemoryCache),
		pubsub:        pubsub,
		config:        &conf,
	}
	go cache.handleUpdates()
	return cache, nil
}

type item struct {
	Data interface{}
}

func (r *RedisAdapter) Set(key *string, value interface{}, config ...*ItemConfig) {
	ttl := r.config.TTL
	if config != nil && config[0].TTL != nil {
		ttl = config[0].TTL
	}

	val, _ := json.Marshal(item{Data: value})

	r.conn.Set(context.Background(), *key, string(val), *ttl)
	r.conn.Publish(context.Background(), "cache_updates:"+*key, "1")
}

func (r *RedisAdapter) Get(key string) (interface{}, bool) {
	if val, ok := r.inMemoryCache.Get(key); ok {
		return val, true
	}
	value, err := r.conn.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false
	}
	var data item
	err = json.Unmarshal([]byte(value), &data)
	if err != nil {
		return nil, false
	}
	ttl := r.conn.TTL(context.Background(), key).Val()
	r.inMemoryCache.Set(&key, data.Data, &ItemConfig{TTL: &ttl})

	return data.Data, true
}

func (r *RedisAdapter) Delete(key string) {
	r.conn.Del(context.Background(), key)
	r.inMemoryCache.Delete(key)
	r.conn.Publish(context.Background(), "cache_updates:"+key, "1")
}

func (r *RedisAdapter) Find(pattern string) []string {
	keys := make([]string, 0)
	keys, _ = r.conn.Keys(context.Background(), pattern).Result()

	return keys
}

func (r *RedisAdapter) handleUpdates() {
	for {
		msg, err := r.pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Printf("Error receiving pub/sub message: %v", err)
			continue
		}
		key := strings.TrimPrefix(msg.Channel, "cache_updates:")
		r.inMemoryCache.Delete(key)
	}
}

func (r *RedisAdapter) FlushAll() {
	// No idea for now
}
