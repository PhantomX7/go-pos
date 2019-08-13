package middleware_test

import (
	"github.com/PhantomX7/go-cleanarch/app/api/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
)

func noopHandler(c *gin.Context)  {
	c.Status(http.StatusOK)
}

type MiddlewareTestSuite struct {
	suite.Suite
	middleware *middleware.Middleware
	handler    http.Handler
}

func (suite *MiddlewareTestSuite) SetupSuite() {
	router := httprouter.New()

	middlewareConfig := middleware.Config{
	}

	m := middleware.New(middlewareConfig)

	// add middleware to the route you want to test
	//e.g router.GET("/", somemiddleware(), noop)


	suite.middleware = m
	suite.handler = router
}

func (suite *MiddlewareTestSuite) Record(request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()

	suite.handler.ServeHTTP(response, request)

	return response
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
