package config

const (
	ApiURL    = "https://api.github.com"
	ApiKeyEnv = "GITHUB_API_TOKEN"
)

var (
	PaginatedApis = map[string]bool{"/orgs/netflix/members": true, "/orgs/netflix/repos": true}
)

type Config struct {
	ApiUrl        string
	ApiKeyEnv     string
	PaginatedApis map[string]bool
}

var staticConfig Config

// LoadConfig returns a static config struct to be embedded in request context
func LoadConfig() Config {
	staticConfig = Config{
		ApiUrl:        ApiURL,
		ApiKeyEnv:     ApiKeyEnv,
		PaginatedApis: PaginatedApis,
	}
	return staticConfig
}
