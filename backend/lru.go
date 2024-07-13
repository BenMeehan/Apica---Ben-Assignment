package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type CacheItem struct {
	Value      string
	Expiration time.Time
}

// A concurrent map-based cache.
type LRUCache struct {
	store    cmap.ConcurrentMap[string, *CacheItem]
	clients  map[*websocket.Conn]bool
	upgrader websocket.Upgrader
}

func NewLRUCache() *LRUCache {
	cache := &LRUCache{
		store:    cmap.New[*CacheItem](),
		clients:  make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}

	// Start a background goroutine to check for expired keys
	go cache.expireChecker()

	return cache
}

// Set stores a value in the cache with an optional expiration time.
func (c *LRUCache) Set(key, value string, expiration int) {
	item, exists := c.store.Get(key)

	if exists {
		item.Value = value
	} else {
		expTime := time.Now()
		if expiration > 0 {
			expTime = time.Now().Add(time.Millisecond * time.Duration(expiration))
		}

		item = &CacheItem{
			Value:      value,
			Expiration: expTime,
		}
	}

	c.store.Set(key, item)
	c.notifyClients("set", key, item)
}

// Delete removes a value from the cache.
func (c *LRUCache) Delete(key string) {
	c.store.Remove(key)
	c.notifyClients("delete", key, nil)
}

// GetAll retrieves all items from the cache and removes expired ones.
func (c *LRUCache) GetAll() map[string]*CacheItem {
	result := make(map[string]*CacheItem)
	for key, item := range c.store.Items() {
		if item.Expiration.IsZero() || item.Expiration.After(time.Now()) {
			result[key] = item
		} else {
			c.Delete(key)
		}
	}
	return result
}
