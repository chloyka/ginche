package ginche

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)
	c := NewCache()
	r := gin.New()
	r.Use(Middleware(c, nil))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	// Run the Set method b.N times
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ping", nil)
		r.ServeHTTP(w, req)
	}
}
