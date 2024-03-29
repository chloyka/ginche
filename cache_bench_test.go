package ginche

import (
	"fmt"
	"testing"
)

func BenchmarkCache_Set(b *testing.B) {
	// Create a new cache with a TTL of 1 minute
	c := NewInMemoryCache()

	// Run the Set method b.N times
	for i := 0; i < b.N; i++ {
		key := String(fmt.Sprintf("key%d", i))
		value := fmt.Sprintf("value%d", i)
		c.Set(key, value)
	}
}

func BenchmarkCache_Get(b *testing.B) {
	// Create a new cache with a TTL of 1 minute
	c := NewInMemoryCache()

	// Add some items to the cache
	for i := 0; i < b.N; i++ {
		key := String(fmt.Sprintf("key%d", i))
		value := fmt.Sprintf("value%d", i)
		c.Set(key, value)
	}

	// Run the Get method b.N times
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		c.Get(key)
	}
}
