package ginche

import (
	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type RedisSuite struct {
	suite.Suite
	store CacheAdapter
	redis *miniredis.Miniredis
}

func (s *RedisSuite) SetupTest() {
	mRedis := miniredis.NewMiniRedis()
	mRedis.Start()
	s.redis = mRedis
	s.store = NewRedisAdapter(&redis.Options{
		Addr: mRedis.Addr(),
	})
}

func (s *RedisSuite) TestSet() {
	key := "test_key"
	value := "test_value"
	s.store.Set(&key, value)
	returnedValue, ok := s.store.Get(key)
	s.True(ok)
	s.Equal(value, returnedValue)
}

func (s *RedisSuite) TestSetWithConfig() {
	key := "test_key"
	value := "test_value"
	ttl := time.Minute
	s.store.Set(&key, value, &ItemConfig{&ttl})
	returnedValue, ok := s.store.Get(key)
	s.True(ok)
	s.Equal(s.redis.TTL(key), ttl)
	s.Equal(value, returnedValue)
}

func (s *RedisSuite) TearDownTest() {
	s.store = nil
	s.redis.Close()
}

func (s *RedisSuite) TestWithMiddleware() {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware(s.store, &Options{}))
	r.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"data": "test",
		})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	d, _ := s.store.Get("/test")
	s.Equal("{\"data\":\"test\"}", d.(map[string]interface{})["Data"])
	s.Equal(w.Body.String(), "{\"data\":\"test\"}")
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}
