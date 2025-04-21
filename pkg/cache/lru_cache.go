package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// LRUCache manages content with Redis LRU eviction policy
type LRUCache struct {
	client *redis.Client
	ctx    context.Context
	prefix string
}

// NewLRUCache creates a Redis cache with 500MB LRU eviction
func NewLRUCache(client *redis.Client) (*LRUCache, error) {
	ctx := context.Background()

	// Configure Redis for LRU with 500MB limit
	if err := client.ConfigSet(ctx, "maxmemory", "100mb").Err(); err != nil {
		return nil, fmt.Errorf("failed to set maxmemory: %w", err)
	}

	if err := client.ConfigSet(ctx, "maxmemory-policy", "allkeys-lru").Err(); err != nil {
		return nil, fmt.Errorf("failed to set LRU policy: %w", err)
	}

	return &LRUCache{
		client: client,
		ctx:    ctx,
		prefix: "content.render.",
	}, nil
}

// Get retrieves content and automatically updates LRU status
func (c *LRUCache) Get(path string) (map[string]any, error) {
	key := c.prefix + path

	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	// Redis automatically updates the LRU on access
	// Track access frequency for popular content
	c.client.ZIncrBy(c.ctx, "content:popular", 1, path)

	var content map[string]any
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, err
	}

	return content, nil
}

// Set stores content with TTL based on popularity
func (c *LRUCache) Set(path string, content map[string]any) error {
	key := c.prefix + path

	data, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// Check popularity score
	score, _ := c.client.ZScore(c.ctx, "content:popular", path).Result()

	ttl := 30 * time.Minute // Default TTL -> TODO: passar para utilizar em uma env
	if score > 5 {
		ttl = 24 * time.Hour // Popular content stays longer -> TODO: passar para utilizar em uma env
	}

	return c.client.Set(c.ctx, key, data, ttl).Err()
}

// PrewarmCache loads top popular content
func (c *LRUCache) PrewarmCache(loader func(string) (map[string]any, error)) error {
	// Get top 20 popular paths
	paths, err := c.client.ZRevRange(c.ctx, "content:popular", 0, 19).Result()
	if err != nil {
		return err
	}

	pipe := c.client.Pipeline()
	for _, path := range paths {
		content, err := loader(path)
		if err != nil {
			continue
		}

		data, _ := json.Marshal(content)
		key := c.prefix + path
		pipe.Set(c.ctx, key, data, 24*time.Hour)
	}

	_, err = pipe.Exec(c.ctx)
	return err
}

// GetMemoryUsage returns current cache memory usage
func (c *LRUCache) GetMemoryUsage() (string, error) {
	info, err := c.client.Info(c.ctx, "memory").Result()
	if err != nil {
		return "", err
	}
	return info, nil
}
