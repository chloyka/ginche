package ginche

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CacheSuite struct {
	suite.Suite
	cache CacheAdapter
}

func (s *CacheSuite) SetupTest() {
	minute := time.Minute
	s.cache = NewInMemoryCache(CacheConfig{
		TTL:             &minute,
		CleanupInterval: &minute,
	})
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

func (s *CacheSuite) TestFlushAll() {
	key := "test_key"
	value := "test_value"
	ttl := time.Minute
	s.cache.Set(&key, value, &ItemConfig{&ttl})
	returnedValue, ok := s.cache.Get(key)
	s.True(ok)
	s.Equal(value, returnedValue)
	s.cache.FlushAll()
	returnedValue, ok = s.cache.Get(key)
	s.False(ok)
	s.Nil(returnedValue)
}

func (s *CacheSuite) TestString() {
	str := "test"
	s.Equal(str, *String(str))
}

func (s *CacheSuite) TestCleanup() {
	testValue := "test"
	s.cache.(*InMemoryCache).cleanupInterval = time.Millisecond * 100
	s.cache.Set(&testValue, testValue, &ItemConfig{
		TTL: &s.cache.(*InMemoryCache).cleanupInterval,
	})
	d, ok := s.cache.Get(testValue)
	s.True(ok)
	s.Equal(testValue, d)
	go s.cache.(*InMemoryCache).cleanup()
	time.Sleep(time.Second)
	d, ok = s.cache.Get(testValue)
	s.False(ok)
	s.Nil(d)
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
