package examples

import (
	"github.com/chloyka/ginche"
	"github.com/redis/go-redis/v9"
)

// Redis example
func main() {
	store := ginche.NewRedisAdapter(&redis.Options{
		Addr: "localhost:6379",
	})
	key := "test_key"
	store.Set(&key, "test_value")
	store.Get("test_key")
}
