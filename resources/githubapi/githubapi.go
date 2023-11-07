package githubapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"netflix-cache/resources/cache"
	"strconv"
)

// New returns a pointer to an instance GithubApi struct
func New(apiUrl url.URL, apiKey string, cache *cache.BigCache, httpClient http.Client) *GithubApi {
	return &GithubApi{
		apiKey: apiKey,
		ApiUrl: apiUrl,
		cache:  cache,
		client: httpClient,
	}
}

// GithubApi struct is a representation of config etc. required to call github APIs
type GithubApi struct {
	apiKey string
	ApiUrl url.URL
	cache  *cache.BigCache
	client http.Client
}

// HandlePaginatedGetRequest handles request to github APIs that require pagination
func (g *GithubApi) HandlePaginatedGetRequest(r *http.Request) ([]byte, error) {
	u := g.ApiUrl
	u.Path = fmt.Sprintf("%s%s", u.Path, r.URL.Path)
	page := 1 // API states we start at page 1
	var responseData []map[string]interface{}
	// Check cache first
	if g.cache != nil {
		if responseBytes, ok := g.cache.Read(r.URL.Path); ok {
			log.Info(fmt.Sprintf("Cache hit on path %s", r.URL.Path))
			return responseBytes, nil
		}
	}
	// Looping forever (yikes)
	for {
		request, err := http.NewRequest("GET", u.String(), r.Body)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to craft HTTP request to path %s", u.String()))
			return nil, err
		}
		// Pass along any query params we may have received
		request.URL.RawQuery = r.URL.Query().Encode()
		// specify page number in query params
		q := request.URL.Query()
		q.Add("page", strconv.Itoa(page))
		request.URL.RawQuery = q.Encode()

		respData, err := g.doRequest(request)
		if err != nil {
			return nil, err
		}
		// Paginated APIs return an empty array when we have exhausted pages, breaking on empty array
		// Maybe could have read the docs more, seems odd there is not a method to determine number of pages
		// or to get a nextToken of some kind while we iterate. Going forward with this in the interest of time
		if bytes.Equal(respData, []byte("[]")) {
			break
		}
		// Response is an array of JSON objects, unmarshal and append to aggregate result
		var body []map[string]interface{}
		err = json.Unmarshal(respData, &body)
		if err != nil {
			return nil, err
		}
		responseData = append(responseData, body...)
		page += 1
	}
	// Marshal all the objects gathered back to bytes
	responseBytes, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}
	if g.cache != nil {
		g.cache.Write(r.URL.Path, responseBytes)
	}
	return responseBytes, nil
}

// HandleGetRequest handles a GET request to certain github APIs where pagination isn't needed
func (g *GithubApi) HandleGetRequest(r *http.Request) ([]byte, error) {
	var respData []byte
	// Check cache first
	if g.cache != nil {
		if respData, ok := g.cache.Read(r.URL.Path); ok {
			log.Info(fmt.Sprintf("Cache hit on path %s", r.URL.Path))
			return respData, nil
		}
	}
	u := g.ApiUrl
	u.Path = fmt.Sprintf("%s%s", u.Path, r.URL.Path)
	log.Info(fmt.Sprintf("Making request at URL %s", u.String()))
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to craft HTTP request to path %s", u.String()))
		return nil, err
	}
	// Pass along any query params we may have received
	request.URL.RawQuery = r.URL.Query().Encode()
	respData, err = g.doRequest(request)
	if err != nil {
		return nil, err
	}
	if g.cache != nil {
		g.cache.Write(r.URL.Path, respData)
	}
	return respData, nil
}

// doRequest helper to embed auth in request if we have it, and roughly examine response
func (g *GithubApi) doRequest(request *http.Request) ([]byte, error) {
	// if we have an API key, set the bearer token per github docs
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28
	if len(g.apiKey) != 0 {
		log.Info("API key found in ENV var, attaching bearer token")
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.apiKey))
	}

	resp, err := g.client.Do(request)
	if err != nil {
		log.Warn("Failed to issue http request with error", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Warn(fmt.Sprintf("github API request failed with status %d", resp.StatusCode))
		return nil, errors.New(fmt.Sprintf("github API request failed with status %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("Failed to read response body", err)
		return nil, err
	}
	return respData, nil
}
