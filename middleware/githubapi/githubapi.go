package githubapi

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	cachemw "netflix-cache/middleware/cache"
	"netflix-cache/middleware/config"
	"netflix-cache/resources/githubapi"
	"os"
)

type contextKey string

const (
	ContextKey contextKey = "middleware/githubapi"
)

// Load inserts a github API helper into the context
func Load(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cache := cachemw.Get(ctx)
		serviceCfg := config.Get(ctx)
		apiUrl, _ := url.Parse(serviceCfg.ApiUrl)
		// Some fragility here, would be better to use almost any other config store and not ENV vars
		apiKey := os.Getenv(serviceCfg.ApiKeyEnv)
		ctx = context.WithValue(ctx, ContextKey, githubapi.New(*apiUrl, apiKey, cache, http.Client{}))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Get retrieves the githubapi from context using key
func Get(ctx context.Context) *githubapi.GithubApi {
	api, ok := ctx.Value(ContextKey).(*githubapi.GithubApi)
	if !ok {
		log.Panic("Github API not in context")
	}
	return api
}
