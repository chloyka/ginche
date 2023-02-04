package ginche

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MiddlewareSuite struct {
	suite.Suite
	store      *Cache
	httpServer *gin.Engine
}

func (s *MiddlewareSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.store = NewCache()
	options := &Options{
		KeyFunc: func(c *gin.Context) string {
			if c.Request.Method == "PATCH" {
				return SkipCacheKeyValue
			}
			return c.Request.URL.Path + c.Request.Method
		},
		ExcludeStatuses: []int{http.StatusCreated},
		ExcludeMethods:  []string{http.MethodPost},
		ExcludePaths:    []string{"/foo"},
	}
	s.httpServer = gin.New()
	s.httpServer.Use(Middleware(s.store, options))

	s.httpServer.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test get"})
	})
	s.httpServer.GET("/test-skip-cache", func(c *gin.Context) {
		c.Set(CTXSkipCacheKey, CTXSkipCacheValue)
		c.JSON(200, gin.H{"message": "test get"})
	})
	s.httpServer.GET("/foo", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "bar"})
	})
	s.httpServer.GET("/bar", func(c *gin.Context) {
		c.JSON(201, gin.H{"message": "baz"})
	})
	s.httpServer.POST("/test", func(c *gin.Context) {
		c.JSON(201, gin.H{"message": "test post"})
	})
}

func (s *MiddlewareSuite) AfterTest() {
	s.store.FlushAll()
}

func (s *MiddlewareSuite) TestGetFromCache() {
	// First request should add result to cache
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	s.httpServer.ServeHTTP(w, req)
	d, match := s.store.Get("/testGET")

	s.True(match)
	s.Equal(200, w.Code)
	s.Equal(`{"message":"test get"}`, w.Body.String())
	s.Equal(200, d.(*httpCacheItem).Status)
	s.Equal("application/json; charset=utf-8", d.(*httpCacheItem).Headers.Get("Content-Type"))
	s.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))
	s.Equal(`{"message":"test get"}`, d.(*httpCacheItem).Data)
	s.Equal(HeaderXCacheMiss, w.Header().Get(HeaderXCache))

	// Second request, should return from cache
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	s.httpServer.ServeHTTP(w, req)
	d, match = s.store.Get("/testGET")

	s.True(match)
	s.Equal(200, w.Code)
	s.Equal(`{"message":"test get"}`, w.Body.String())
	s.Equal(HeaderXCacheHit, w.Header().Get(HeaderXCache))
}

// Should return miss cache because of POST method excluded
func (s *MiddlewareSuite) TestSkipExcludedMethods() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	s.httpServer.ServeHTTP(w, req)
	d, match := s.store.Get("/testPOST")

	s.False(match)
	s.Nil(d)
	s.Equal(201, w.Code)
	s.Equal(`{"message":"test post"}`, w.Body.String())
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
}

// Should return miss cache because of Skip cache CTX value
func (s *MiddlewareSuite) TestMissExcludedKey() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-skip-cache", nil)

	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
}

// Should return miss cache because of Skip cache CTX value
func (s *MiddlewareSuite) TestSkipCacheWithExcludedPath() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
}

// Should return miss cache because of Skip cache CTX value
func (s *MiddlewareSuite) TestSkipCacheWithExcludedStatuses() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bar", nil)
	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
}

// Should return miss cache because of Skip cache CTX value
func (s *MiddlewareSuite) TestSkipAllCachesWithKeyGeneratorSkipValue() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/test", nil)
	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
	s.httpServer.ServeHTTP(w, req)
	s.Equal(HeaderXCacheSkip, w.Header().Get(HeaderXCache))
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
