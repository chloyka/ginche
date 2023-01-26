package middleware

import (
	"github.com/chloyka/ginche/cache"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MiddlewareSuite struct {
	suite.Suite
}

func (s *MiddlewareSuite) TestServe() {
	storage := cache.NewCache(time.Minute, nil)
	options := &Options{
		KeyFunc: func(c *gin.Context) string {
			return c.Request.URL.Path + c.Request.Method
		},
		ExcludeStatuses: []int{http.StatusNotFound},
		ExcludeMethods:  []string{http.MethodPost},
	}
	middleware := Serve(storage, options)

	// Test caching
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test get"})
	})

	router.POST("/test", func(c *gin.Context) {
		c.JSON(201, gin.H{"message": "test post"})
	})
	// First request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	d, match := storage.Get("/testGET")
	s.True(match)
	s.Equal(200, w.Code)
	s.Equal(`{"message":"test get"}`, w.Body.String())
	s.Equal(200, d.(*httpCacheItem).Status)
	s.Equal("application/json; charset=utf-8", d.(*httpCacheItem).Headers.Get("Content-Type"))
	s.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))
	s.Equal(`{"message":"test get"}`, d.(*httpCacheItem).Data)
	s.Equal(`MISS`, w.Header().Get("X-Cache"))

	// Second request, should return from cache
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	d, match = storage.Get("/testGET")
	s.True(match)
	s.Equal(200, w.Code)
	s.Equal(`{"message":"test get"}`, w.Body.String())
	s.Equal(`HIT`, w.Header().Get("X-Cache"))

	// Should return miss cache because of POST method excluded
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)
	d, match = storage.Get("/testPOST")
	s.False(match)
	s.Nil(d)
	s.Equal(201, w.Code)
	s.Equal(`{"message":"test post"}`, w.Body.String())
}

func TestServe(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
