package ginche

import (
	"context"
	"encoding/json"
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

type item struct {
	Data interface{}
}

func (r *RedisAdapter) Set(key *string, value interface{}, config ...*ItemConfig) {
	ttl := r.config.TTL
	if config != nil && config[0].TTL != nil {
		ttl = config[0].TTL
	}
	if config != nil {
		ttl = config[0].TTL
	}
	val, _ := json.Marshal(item{Data: value})

	r.conn.Set(context.Background(), *key, string(val), *ttl)
}

func (r *RedisAdapter) Get(key string) (interface{}, bool) {
	value, err := r.conn.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false
	}
	var data item
	err = json.Unmarshal([]byte(value), &data)
	if err != nil {
		return nil, false
	}

	return data.Data, true
}

func (r *RedisAdapter) FlushAll() {
	// No idea for now
}
