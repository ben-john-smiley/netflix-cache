package monitor

import "net/http"

// Get responds with a 200 if the service is up and running, no deep health checks
func Get(w http.ResponseWriter, req *http.Request) {}
