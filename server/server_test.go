package server

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"testing"
	"time"
)

// ServerTestSuite ...
type ServerTestSuite struct {
	suite.Suite
	httpServ *graceful.Server
}

// SetupTest inits server, passes a nil cache
func (s *ServerTestSuite) SetupTest() {
	var err error
	mux := http.NewServeMux()
	mux.HandleFunc("/", New(nil).ServeHTTP)
	s.httpServ = &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
	go func() {
		err = s.httpServ.ListenAndServe()
		s.Error(err)
	}()
	time.Sleep(300 * time.Millisecond) // give service some time to initialize
}

// TestStart just checks if the server can start and we can hit healthcheck
func (s *ServerTestSuite) TestStart() {
	_, err := http.Get("http://localhost:8080/healthcheck")
	s.NoError(err)
}

// TearDownTest ensures we close the server
func (s *ServerTestSuite) TearDownTest() {
	err := s.httpServ.Close()
	s.NoError(err)
}

// TestServer runs the suite
func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
