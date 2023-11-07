# netflix-cache
A service that utilizes an in-memory cache to cache responses from Github organzations APIs
* `/orgs/Netflix`
* `/orgs/Netflix/repos`
* `/orgs/Netflix/members`

Responses from paginated endpoints (`repos` && `members`) are collected into one blob and cached/returned to callers in one call to netflix-cache service

All other requests made to the service are reverse proxied to the GitHub API

Service provides `Bottom-N` views into Netflix repos via endpoint:
* `/views/bottom/{N}/forks`
* `/views/bottom/{N}/open_issues`
* `/views/bottom/{N}/stars`
* `/views/bottom/{N}/last_updated`

Example:
```
$ curl localhost:8080/view/bottom/5/forks | jq .
[
  [
    "Netflix/apache-pyiceberg",
    0
  ],
  [
    "Netflix/eclipse-mat",
    0
  ],
  [
    "Netflix/iceberg-python",
    0
  ],
  [
    "Netflix/mantis-helm",
    0
  ],
  [
    "Netflix/octodns",
    0
  ]
]

```
# Prerequisites
### Required
* Runs on `golang` version `1.21.3`
    * Instructions to install golang [here](https://go.dev/doc/install)
### Optional
* Specify a GitHub API key in env var `GITHUB_API_TOKEN`, service will apply this token to all requests to github

# Testing
`cd /netflix-cache && go test ./...`

TODO: add comprehensive tests

Only testing happy path of a few components

# Building
`cd /netflix-cache && go get && go build netflix-cache.go`

The above command should produce a `netflix-cache` executable
# Running
Service accepts a port parameter on startup, if not provided the default port is `8080`
```azure
 $ ./netflix-cache -p 8080
 {"level":"info","msg":"Starting netflix-cache service on port 8080","time":"2023-11-06T16:59:29-08:00"}
```

Log messages will continue to trail in the terminal, and the service is ready for traffic.

Shut down the service via `ctrl+c` in the terminal
```azure
$ ./netflix-cache -p 8080
{"level":"info","msg":"Starting netflix-cache service on port 8080","time":"2023-11-06T16:59:58-08:00"}
^C[graceful] shutdown initiated
{"level":"fatal","msg":"\u003cnil\u003e","time":"2023-11-06T17:00:37-08:00"}
```

# Cache Configuration
The cache configuration is hard coded into `/netflix-cache/resources/cache/cache.go`

The cache is configured to have a maximum size of 2 GB, oldest members will be evicted if the cache attempts to grow beyond this.

The cache holds elements in the store for a configured time of one hour. Considered the option of checking etags at the github API endpoints to determine if the cached data was stale but opted to just have a set life-span of cache members for the time being.