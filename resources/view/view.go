package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"netflix-cache/resources/githubapi"
	"netflix-cache/resources/repos"
	"strings"
)

// HandleBottomN handles a request to report the bottom-N repos given a filter
func HandleBottomN(n int, apiClient *githubapi.GithubApi, urlPath string) ([]interface{}, error) {
	// less than pretty method of getting the slug from a urlPath to determine sort priority/response type
	sortField := urlPath[strings.LastIndex(urlPath, "/")+1:]
	repos, err := parseRepos(apiClient, "/orgs/Netflix/repos") // Hard coding in cache key here, should live in config
	if err != nil {
		return nil, err
	}
	return repos.RetrieveBottomN(n, sortField), nil
}

// parseRepos takes the response from the githubapi/cache and marshals it into an array of the slimmed down Repo struct
func parseRepos(apiClient *githubapi.GithubApi, urlPath string) (repos.Repos, error) {
	var repos repos.Repos
	// need a dummy request here so we can continue to preserve URL params in githubapi
	dummyRequest, err := http.NewRequest("GET", fmt.Sprintf("%s%s", apiClient.ApiUrl.String(), urlPath), nil)
	if err != nil {
		return nil, err
	}
	repoBytes, err := apiClient.HandlePaginatedGetRequest(dummyRequest)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(repoBytes, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}
