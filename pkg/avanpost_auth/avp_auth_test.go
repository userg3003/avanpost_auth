package avanpost_auth_test

import (
	"avanpost_auth/internal/helpers"
	"avanpost_auth/internal/utils"
	"avanpost_auth/pkg/avanpost_auth"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type EndpointsTestSuite struct {
	suite.Suite
}

func TestEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(EndpointsTestSuite))
}

func (suite *EndpointsTestSuite) TestHealthRoute() {
	assert_ := assert.New(suite.T())
	router := avanpost_auth.SetupRouter(false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert_.Equal(204, w.Code)
	assert_.Equal("", w.Body.String())
}

func (suite *EndpointsTestSuite) TestAuthRedirect() {
	assert_ := assert.New(suite.T())
	helpers.InitConfigForTests("app_test.env")
	router := avanpost_auth.SetupRouter(false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	router.ServeHTTP(w, req)

	assert_.Equal(302, w.Code)
	expected := fmt.Sprintf("<a href=\"%s\">Found</a>.\n\n", utils.FullUrlForAuthorize("xyz", "profile"))
	expected = strings.ReplaceAll(expected, "&", "&amp;")
	assert_.Equal(expected, w.Body.String())
}
