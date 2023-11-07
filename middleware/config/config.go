package config

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"netflix-cache/resources/config"
)

type contextKey string

const ContextKey contextKey = "middleware/config"

// Load inserts the static service config in the context
func Load(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKey, config.LoadConfig())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Get retrieves the config from context using key
func Get(ctx context.Context) config.Config {
	cfg, ok := ctx.Value(ContextKey).(config.Config)
	if !ok {
		log.Panic("Config not in context")
	}
	return cfg
}
