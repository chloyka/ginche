package ginche

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisAdapter struct {
	conn   *redis.Client
	config *CacheConfig
}

func NewRedisAdapter(redisConfig *redis.Options, config ...CacheConfig) CacheAdapter {
	redisClient := redis.NewClient(redisConfig)
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

	return &RedisAdapter{
		conn:   redisClient,
		config: &conf,
	}
}

func (r *RedisAdapter) Set(key *string, value interface{}, config ...*ItemConfig) {
	ttl := r.config.TTL
	if config != nil && config[0].TTL != nil {
		ttl = config[0].TTL
	}
	if config != nil {
		ttl = config[0].TTL
	}

	r.conn.Set(context.Background(), *key, value, *ttl)
}

func (r *RedisAdapter) Get(key string) (interface{}, bool) {
	value, err := r.conn.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false
	}
	return value, true
}

func (r *RedisAdapter) FlushAll() {
	// No idea for now
}
