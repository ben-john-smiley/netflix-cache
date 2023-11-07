package view

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	apimw "netflix-cache/middleware/githubapi"
	"netflix-cache/resources/view"
	"strconv"
)

// View handles calls made to the service at the `/view/bottom/{N}/{filter}` endpoints
func View(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(chi.URLParam(r, "nAmount"))
	if err != nil {
		// Just going to write bad request and a generic error if the url param isn't an integer. Might be better to
		// re-route to github reverse proxy if it is not within route definitions
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to parse url param '%s' as integer", chi.URLParam(r, "nAmount"))))
		return
	}
	githubApi := apimw.Get(r.Context())
	responseData, err := view.HandleBottomN(n, githubApi, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // again, need well-formed error, should be careful what we respond with
		return
	}

	respJson, err := json.Marshal(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // again, need well-formed error, should be careful what we respond with
		return
	}
	w.Write(respJson)
}
