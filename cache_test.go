package ginche

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CacheSuite struct {
	suite.Suite
	cache *Cache
}

func (s *CacheSuite) SetupTest() {
	s.cache = NewCache(time.Minute, nil)
}

func (s *CacheSuite) TearDownTest() {
	s.cache = nil
}

func (s *CacheSuite) TestSet() {
	key := "test_key"
	value := "test_value"
	s.cache.Set(&key, value)
	returnedValue, ok := s.cache.Get(key)
	s.True(ok)
	s.Equal(value, returnedValue)
}

func (s *CacheSuite) TestSetWithConfig() {
	key := "test_key"
	value := "test_value"
	ttl := time.Second
	s.cache.Set(&key, value, &ItemConfig{&ttl})
	returnedValue, ok := s.cache.Get(key)
	s.True(ok)
	s.Equal(value, returnedValue)
	time.Sleep(2 * time.Second)
	returnedValue, ok = s.cache.Get(key)
	s.False(ok)
	s.Nil(returnedValue)
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
