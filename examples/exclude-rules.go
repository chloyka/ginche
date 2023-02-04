package examples

import (
	"github.com/chloyka/ginche"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Example with exclusion rules
func main() {
	// TTL is 1 minute, cleanup interval is also 1 minute
	store := ginche.NewCache()
	r := gin.New()

	// Exclude POST and PUT requests from caching
	// Also exclude 404 responses from caching
	// Exclude requests to /foo from caching
	r.Use(ginche.Middleware(store, &ginche.Options{
		ExcludeMethods:  []string{"POST", "PUT"},
		ExcludeStatuses: []int{http.StatusNotFound},
		ExcludePaths:    []string{"/foo"},
	}))

	// Will exclude from caching because of the path rule
	r.GET("/foo", func(c *gin.Context) {
		c.String(200, "bar")
	})

	// Will exclude from caching because of setting the context skip cache value
	r.GET("/foo2", func(c *gin.Context) {
		c.Set(ginche.CTXSkipCacheKey, ginche.CTXSkipCacheValue)

		c.String(200, "bar")
	})

	// Will exclude from caching because of the statuses rule
	r.POST("/bar", func(c *gin.Context) {
		c.String(200, "baz")
	})

	// Will exclude from caching because of the methods rule
	r.GET("/baz", func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "not found"})
	})

	// Will be cached
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}
