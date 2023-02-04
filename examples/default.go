package examples

import (
	"github.com/chloyka/ginche"
	"github.com/gin-gonic/gin"
)

// Default example
// Uses the middleware with default options
func main() {
	// TTL is 1 minute, cleanup interval is also 1 minute
	store := ginche.NewCache()
	r := gin.New()
	// No need to specify the options, just pass nil
	// No exclusion rules, default key is (Request.URL.Path)
	r.Use(ginche.Middleware(store, nil))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
