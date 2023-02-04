package examples

import (
	"github.com/chloyka/ginche"
	"github.com/gin-gonic/gin"
)

// Example of generation of cache key
func main() {
	// TTL is 1 minute, cleanup interval is also 1 minute
	store := ginche.NewCache()
	r := gin.New()
	r.Use(ginche.Middleware(store, &ginche.Options{
		KeyFunc: func(c *gin.Context) string {
			// Use the full URL as the cache key
			if c.Request.Method == "GET" {
				return c.Request.URL.String()
			} else if c.Request.Method == "POST" {
				// Use the full URL and the name field from the body as the cache key
				var body map[string]interface{}
				err := c.BindJSON(&body)
				if err != nil {
					// If there is an error, skip caching
					return ginche.SkipCacheKeyValue
				}
				return c.Request.URL.String() + body["name"].(string)
			} else {
				// Otherwise, skip caching
				return ginche.SkipCacheKeyValue
			}
		},
	}))
	// Will be cached
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Will be cached
	r.POST("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
