package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var cache *LRUCache

type CacheResponse struct {
	Value      string `json:"value"`
	Expiration int64  `json:"expiration"`
}

func main() {
	cache = NewLRUCache()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "WS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/cache", SetCacheHandler)
	r.DELETE("/cache/:key", DeleteCacheHandler)
	r.GET("/cache", GetAllCacheHandler)
	r.GET("/ws", WebSocketHandler)

	log.Fatal(r.Run(":8080"))
}
