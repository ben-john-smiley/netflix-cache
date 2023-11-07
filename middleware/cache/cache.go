package cache

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"netflix-cache/resources/cache"
)

type contextKey string

const ContextKey contextKey = "middleware/cache"

// CheckCache attempts to read the request path from cache and responds with that payload if a cache hit
// otherwise middleware allows request to continue to handler
func CheckCache(cache *cache.BigCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cache != nil {
				cachedValue, ok := cache.Read(r.URL.Path)
				if ok {
					log.Info(fmt.Sprintf("Cache hit on route %s", r.URL.Path))
					w.Header().Set("Content-Type", "application/json")
					w.Write(cachedValue)
					return
				}
			}
			log.Info(fmt.Sprintf("Cache miss on route %s", r.URL.Path))
			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}

// Load inserts the cache in the request context
func Load(cache *cache.BigCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKey, cache)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Get retrieves the cache from context using key
func Get(ctx context.Context) *cache.BigCache {
	cache, ok := ctx.Value(ContextKey).(*cache.BigCache)
	if !ok {
		log.Panic("Cache not in context")
	}
	return cache
}
