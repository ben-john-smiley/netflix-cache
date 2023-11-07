package config_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"netflix-cache/middleware/config"
	serviceconfig "netflix-cache/resources/config"
	"netflix-cache/test/mocks"
	"testing"
)

// ConfigSuite ...
type ConfigSuite struct {
	suite.Suite
	ctx    context.Context
	req    *http.Request
	mwMock *mocks.MockMiddleware
}

// SetupTest ...
func (s *ConfigSuite) SetupTest() {
	s.req, _ = http.NewRequest("GET", "test.netflix.com", nil)
	s.ctx = s.req.Context()
	s.mwMock = &mocks.MockMiddleware{}
}

// TestConfigLoad tests that the Middleware injects the config in the context
func (s *ConfigSuite) TestConfigLoad() {
	w := httptest.NewRecorder()
	s.mwMock.On("ServeHTTP", w, mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		s.ctx = args.Get(1).(*http.Request).Context()
	}).Once()

	config.Load(s.mwMock).ServeHTTP(w, s.req.WithContext(s.ctx))

	conf := config.Get(s.ctx)
	s.Equal(serviceconfig.LoadConfig(), conf)
}

// TestConfigSuite ...
func TestConfigSuite(t *testing.T) {
	suite.Run(t, &ConfigSuite{})
}
