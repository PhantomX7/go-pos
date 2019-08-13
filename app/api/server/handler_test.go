package server_test

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhantomX7/go-cleanarch/app/api/middleware"
	"github.com/PhantomX7/go-cleanarch/app/api/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type noopHandler struct{}

func (nh *noopHandler) Register(r *gin.Engine, m *middleware.Middleware) {}

// suite can implement setupsuite, teardownsuite, etc.
// check the Interface for more info
type ServerTestSuite struct {
	suite.Suite
	handler http.Handler
}

// setup the test
func (suite *ServerTestSuite) SetupSuite() {
	mConfig := middleware.Config{
	}
	suite.handler = server.BuildHandler(
		middleware.New(mConfig),
		&noopHandler{},
	)
}

// executed when the test finish the suite
func (suite *ServerTestSuite) TearDownSuite() {

}

func (suite *ServerTestSuite) TestRootReturns404() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nothing", nil)

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *ServerTestSuite) TestHealthzReturnsOk() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "ok", w.Body.String())
}

func (suite *ServerTestSuite) TestMetricsReturnsOk() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)

	suite.handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
