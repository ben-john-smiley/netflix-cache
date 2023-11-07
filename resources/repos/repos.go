package repos

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// Repos is a minimal representation of the repo definition in the github api response from /orgs/{ORG}/repos
// only defining fields needed to achieve the bottom-N implementation
type Repos []Repo
type Repo struct {
	Forks      int       `json:"forks"`
	OpenIssues int       `json:"open_issues_count"`
	Stars      int       `json:"stargazers_count"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       string    `json:"full_name"`
}

// sortByField sorts an array of repos by a given field name, little bit of code smell here but in the
// interests of time just using a switch to determine how to sort
func (r Repos) sortByField(field string) Repos {
	// Sorting by desired field first, using name of the repo as a fallback
	switch field {
	case "forks":
		sort.Slice(r, func(i, j int) bool {
			if r[i].Forks != r[j].Forks {
				return r[i].Forks < r[j].Forks
			}
			return r[i].Name < r[j].Name
		})
	case "stars":
		sort.Slice(r, func(i, j int) bool {
			if r[i].Stars != r[j].Stars {
				return r[i].Stars < r[j].Stars
			}
			return r[i].Name < r[j].Name
		})
	case "open_issues":
		sort.Slice(r, func(i, j int) bool {
			if r[i].OpenIssues != r[j].OpenIssues {
				return r[i].OpenIssues < r[j].OpenIssues
			}
			return r[i].Name < r[j].Name
		})
	case "last_updated":
		sort.Slice(r, func(i, j int) bool {
			if r[i].UpdatedAt != r[j].UpdatedAt {
				return r[i].UpdatedAt.Before(r[j].UpdatedAt)
			}
			return r[i].Name < r[j].Name
		})
	}
	return r
}

// getResponseInterface returns an interface to supply view with requested repo information
func (r Repo) getResponseInterface(field string) []interface{} {
	switch field {
	case "forks":
		return []interface{}{r.Name, r.Forks}
	case "stars":
		return []interface{}{r.Name, r.Stars}
	case "open_issues":
		return []interface{}{r.Name, r.OpenIssues}
	case "last_updated":
		return []interface{}{r.Name, r.UpdatedAt}
	}
	return nil
}

// RetrieveBottomN sorts repos by the given field and then returns the bottom n elements
func (r Repos) RetrieveBottomN(n int, field string) []interface{} {
	r = r.sortByField(field)
	var response []interface{}
	// If the user wants more bottom-N values than we have, just give them all of them
	if n > len(r) {
		log.Info(fmt.Sprintf("Bottom %d requested, only %d records available", n, len(r)))
		n = len(r)
	}
	for i := 0; i < n; i++ {
		response = append(response, r[i].getResponseInterface(field))
	}
	return response
}
