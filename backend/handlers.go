package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// WebSocketHandler handles WebSocket connections and sends updates.
func WebSocketHandler(c *gin.Context) {
	conn, err := cache.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	cache.clients[conn] = true
	defer delete(cache.clients, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
	}
}

// SetCacheHandler handles POST requests to set cache values.
func SetCacheHandler(c *gin.Context) {
	var request struct {
		Key        string `json:"key"`
		Value      string `json:"value"`
		Expiration int    `json:"expiration"` // in ms
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cache.Set(request.Key, request.Value, request.Expiration)
	c.Status(200)
}

// DeleteCacheHandler handles DELETE requests to remove cache values.
func DeleteCacheHandler(c *gin.Context) {
	key := c.Param("key")

	cache.Delete(key)
	c.Status(200)
}

// GetAllCacheHandler handles GET requests to retrieve all cache items.
func GetAllCacheHandler(c *gin.Context) {
	items := cache.GetAll()
	response := make(map[string]CacheResponse)
	for key, item := range items {
		response[key] = CacheResponse{
			Value:      item.Value,
			Expiration: item.Expiration.UnixNano(),
		}
	}
	c.JSON(200, response)
}