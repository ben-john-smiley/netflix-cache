package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	apihandler "netflix-cache/handlers/githubapi"
	"netflix-cache/handlers/monitor"
	"netflix-cache/handlers/view"
	cachemw "netflix-cache/middleware/cache"
	"netflix-cache/middleware/config"
	apimw "netflix-cache/middleware/githubapi"
	"netflix-cache/resources/cache"
)

// New initializes chi infrastructure
func New(cache *cache.BigCache) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	// simple health check route
	router.Get("/healthcheck", monitor.Get)
	// handle routes doc suggests should not be reverse proxied to github
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.Recoverer)      // panic recovery
		r.Use(middleware.StripSlashes)   // no trailing slashes
		r.Use(config.Load)               // service configuration
		r.Use(cachemw.Load(cache))       // Load cache into context
		r.Use(cachemw.CheckCache(cache)) // check if this routes response is cached
		r.Use(apimw.Load)                // github API helper
		r.Get("/", apihandler.NetflixGithubApiHandler)
	})
	router.Route("/orgs/Netflix", func(r chi.Router) {
		r.Use(middleware.Recoverer)      // panic recovery
		r.Use(middleware.StripSlashes)   // no trailing slashes
		r.Use(config.Load)               // service configuration
		r.Use(cachemw.Load(cache))       // Load cache into context
		r.Use(cachemw.CheckCache(cache)) // check if this routes response is cached
		r.Use(apimw.Load)                // github API helper
		r.Get("/", apihandler.NetflixGithubApiHandler)
		r.Get("/members", apihandler.NetflixGithubApiHandler)
		r.Get("/repos", apihandler.NetflixGithubApiHandler)
	})

	// Define view routes that will read from the cache to provide bottom-N views
	router.Route("/view/bottom/{nAmount}", func(r chi.Router) {
		r.Use(middleware.Recoverer)    // panic recovery
		r.Use(middleware.StripSlashes) // no trailing slashes
		r.Use(config.Load)             // service configuration
		r.Use(cachemw.Load(cache))     // Load cache into context
		r.Use(apimw.Load)              // github API helper
		r.Get("/forks", view.View)
		r.Get("/last_updated", view.View)
		r.Get("/open_issues", view.View)
		r.Get("/stars", view.View)
	})

	// Little bit hacky here, if the route is not found in the above route definitions
	// sending it through a reverse proxy to githubapi
	router.NotFound(apihandler.ReverseProxyHandler())

	return router
}
