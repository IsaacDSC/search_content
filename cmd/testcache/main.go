package testcache

import (
	"github.com/IsaacDSC/search_content/pkg/cache"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {

	// In cmd/api/main.go
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	// Create LRU cache with 500MB limit
	lruCache, err := cache.NewLRUCache(client)
	if err != nil {
		log.Fatal("Failed to configure LRU cache:", err)
	}

	// Store content
	content := map[string]any{
		"Video": map[string]any{
			"VideoUrl":    "https://video1.com.br",
			"TambnailUrl": "https://thumbnail1.com.br",
		},
	}
	lruCache.Set("/home/camisa/masculina", content)

	// Get content - automatically updates LRU status
	content, err = lruCache.Get("/home/camisa/masculina")
}
