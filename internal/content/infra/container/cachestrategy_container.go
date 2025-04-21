package container

import (
	"github.com/IsaacDSC/search_content/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type CacheStrategies struct {
	LRUCache *cache.LRUCache
}

func NewCacheStrategies(rdc *redis.Client) CacheStrategies {
	lruCache, err := cache.NewLRUCache(rdc)
	if err != nil {
		panic("Failed to initialize cache: " + err.Error())
	}

	return CacheStrategies{
		lruCache,
	}
}
