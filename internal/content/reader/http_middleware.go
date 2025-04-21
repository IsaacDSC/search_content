package reader

import (
	"encoding/json"
	"github.com/IsaacDSC/search_content/pkg/cache"
	"log"
	"net/http"
	"strings"
)

// CacheMiddleware provides caching capabilities for HTTP handlers
type CacheMiddleware struct {
	cache *cache.LRUCache
}

// NewCacheMiddleware creates a new cache middleware
func NewCacheMiddleware(cache *cache.LRUCache) *CacheMiddleware {
	return &CacheMiddleware{
		cache: cache,
	}
}

// WithCache wraps an HTTP handler with caching logic
func (m *CacheMiddleware) WithCache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip cache for non-GET requests
		if r.Method != http.MethodGet {
			next(w, r)
			return
		}

		// Generate cache key from request
		cacheKey := generateCacheKey(r)

		// Try to get from cache
		data, err := m.cache.Get(cacheKey)
		if err == nil {
			// Cache hit
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			json.NewEncoder(w).Encode(data)
			return
		}

		// Cache miss - use response writer wrapper to capture response
		crw := newCaptureResponseWriter(w)

		// Execute the original handler
		next(crw, r)

		// If successful response, store in cache
		if crw.statusCode >= 200 && crw.statusCode < 300 && len(crw.body) > 0 {
			body := make(map[string]any)
			if err := json.Unmarshal(crw.body, &body); err != nil {
				log.Println("[WARNING] Failed to unmarshal response body:", err)
			}
			m.cache.Set(cacheKey, body)
		}
	}
}

// Helper function to generate cache key from request
func generateCacheKey(r *http.Request) string {
	key := r.URL.Path
	if r.URL.RawQuery != "" {
		key += "?" + r.URL.RawQuery
	}
	return strings.TrimPrefix(key, "/")
}

// captureResponseWriter is a wrapper that captures response data
type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

// newCaptureResponseWriter creates a new response writer wrapper
func newCaptureResponseWriter(w http.ResponseWriter) *captureResponseWriter {
	return &captureResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// Write captures the response body
func (crw *captureResponseWriter) Write(b []byte) (int, error) {
	crw.body = append(crw.body, b...)
	return crw.ResponseWriter.Write(b)
}

// WriteHeader captures the status code
func (crw *captureResponseWriter) WriteHeader(statusCode int) {
	crw.statusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
