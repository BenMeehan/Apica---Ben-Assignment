package main

import (
	"log"
	"time"
)

func (c *LRUCache) notifyClients(event, key string, item *CacheItem) {
	var value string
	var expiration int64

	if item != nil {
		value = item.Value
		expiration = item.Expiration.UnixNano()
	}

	data := map[string]interface{}{
		"event":      event,
		"key":        key,
		"value":      value,
		"expiration": expiration,
	}

	for client := range c.clients {
		err := client.WriteJSON(data)
		if err != nil {
			log.Println("Error sending update to client:", err)
			client.Close()
			delete(c.clients, client)
		}
	}
}

// expireChecker periodically checks for expired items and notifies clients.
func (c *LRUCache) expireChecker() {
	for {
		time.Sleep(1 * time.Second) // Check for expired items every second
		now := time.Now()

		for key, item := range c.store.Items() {
			if !item.Expiration.IsZero() && item.Expiration.Before(now) {
				c.Delete(key)
			}
		}
	}
}