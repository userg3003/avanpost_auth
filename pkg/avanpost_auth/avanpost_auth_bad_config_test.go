package avanpost_auth_test

import (
	"avanpost_auth/internal/helpers"
	"avanpost_auth/internal/utils"
	"avanpost_auth/pkg/avanpost_auth"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type BadConfigTestSuite struct {
	suite.Suite
}

func TestBadConfigTestSuite(t *testing.T) {
	suite.Run(t, new(BadConfigTestSuite))
}

func (suite *BadConfigTestSuite) TestBadOAut2AuthorizePath() {
	if os.Getenv("APITEST") == "" {
		suite.T().Skip("skipping test; $APITEST not set")
	}
	assert := assert.New(suite.T())
	configName := "app_test_bad_config.env"
	srv, osinSrv, ctx := helpers.StartServers(configName)
	defer helpers.ShutdownServers(ctx, srv, osinSrv)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/auth", viper.GetInt("SERVICE_PORT")))
	assert.Nil(err, "Error authorize")
	assert.Equal(404, resp.StatusCode)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.Contains(string(body), "404 page not found")
}

func (suite *BadConfigTestSuite) TestAuthRedirect1() {
	//if os.Getenv("APITEST") == "" && os.Getenv("ALLTESTS") != "1" {
	//	suite.T().Skip("skipping test; $APITEST not set")
	//}
	assert := assert.New(suite.T())
	helpers.InitConfigForTests("app_test_bad_config.env")
	router := avanpost_auth.SetupRouter(false)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/auth", nil)
	assert.Nil(err, "Error get config")
	router.ServeHTTP(w, req)

	assert.Equal(302, w.Code)
	expected := fmt.Sprintf("<a href=\"%s\">Found</a>.\n\n", utils.FullUrlForAuthorize("xyz"))
	expected = strings.ReplaceAll(expected, "&", "&amp;")
	assert.Equal(expected, w.Body.String())
}
