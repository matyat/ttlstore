package ttlstore

import (
	"sync"
	"time"
)

// Simple key-value store with TTL for values. Safe for concurrent access.
type Store struct {
	mu         sync.RWMutex
	items      map[string]item
	defaultTTL time.Duration
	counter    uint32
}

type item struct {
	v          interface{}
	expiryTime time.Time
}

// Create a new key-value store with a TTL for keys.
func New(ttl time.Duration) *Store {
	return &Store{
		items:      make(map[string]item),
		defaultTTL: ttl,
	}

}

// Get value of a given key. Returns nil if the key doesn't exist, or has expired.
func (c *Store) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	i, ok := c.items[key]
	if !ok {
		return nil
	}

	if time.Now().After(i.expiryTime) {
		// expired
		return nil
	}

	return i.v
}

// Set a value.
func (c *Store) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// Set a value with a custom TTL.
func (c *Store) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	i := item{value, time.Now().Add(ttl)}
	c.items[key] = i

	c.clean()
}

// Internal method, clean out expired keys. Must be Locked before calling.
func (c *Store) clean() {
	c.counter++
	if (c.counter % 10 != 0) {
		return
	}

	now := time.Now()

	for k, i := range c.items {
		if now.After(i.expiryTime) {
			delete(c.items, k)
		}
	}
}
