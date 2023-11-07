package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
	"netflix-cache/resources/config"
)

// MockMiddleware ...
type MockMiddleware struct {
	mock.Mock
}

// ServeHTTP ...
func (mw *MockMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw.Called(w, r)
}

// ServiceConfig ...
func (mw *MockMiddleware) ServiceConfig(c context.Context) config.Config {
	args := mw.Called(c)
	return args.Get(0).(config.Config)
}
