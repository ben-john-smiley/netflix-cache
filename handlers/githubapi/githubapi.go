package githubapi

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"netflix-cache/middleware/config"
	apimw "netflix-cache/middleware/githubapi"
	"os"
	"strings"
)

// NetflixGithubApiHandler handles GET requests to the github API for a few endpoints
func NetflixGithubApiHandler(w http.ResponseWriter, r *http.Request) {
	paginatedEndpoints := config.Get(r.Context()).PaginatedApis
	githubApi := apimw.Get(r.Context())
	var responseBytes []byte
	var err error

	if paginatedEndpoints[strings.ToLower(r.URL.Path)] {
		responseBytes, err = githubApi.HandlePaginatedGetRequest(r)
	} else {
		responseBytes, err = githubApi.HandleGetRequest(r)
	}
	// Just responding with the error directly from the API resources, ideally would define a well-formed error model
	// and unwrap the errors at this level and respond with whatever level of detail is necessary
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(responseBytes)
}

// ReverseProxyHandler is used to reverse proxy requests made to this service to githubs API
func ReverseProxyHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		apiUrl, err := url.Parse(config.Get(r.Context()).ApiUrl)
		apiKey := os.Getenv(config.Get(r.Context()).ApiKeyEnv)

		// Service should panic if we cannot hook up the reverse proxy
		if err != nil {
			log.Panic("Failed to parse reverse proxy URL", err)
		}

		proxy := httputil.ReverseProxy{Director: func(r *http.Request) {
			r.URL.Scheme = apiUrl.Scheme
			r.URL.Host = apiUrl.Host
			r.URL.Path = apiUrl.Path + r.URL.Path
			r.Host = apiUrl.Host
			if len(apiKey) != 0 {
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
			}
		}}
		log.Info(fmt.Sprintf("Making proxy request to host: %s at path %s", apiUrl.Host, r.URL.Path))
		proxy.ServeHTTP(w, r)
	}
}
