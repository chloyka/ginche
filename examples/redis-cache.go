package examples

import (
	"github.com/chloyka/ginche"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Redis example
func main() {
	r := gin.New()
	// No need to specify the options, just pass nil
	// No exclusion rules, default key is (Request.URL.Path)
	r.Use(ginche.Middleware(ginche.NewRedisAdapter(&redis.Options{
		Addr: "localhost:6379",
	}), nil))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
