package avanpost_auth_test

import (
	"avanpost_auth/internal/helpers"
	"avanpost_auth/internal/utils"
	"avanpost_auth/pkg/avanpost_auth"
	"fmt"
	"github.com/rs/zerolog/log"
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
	assert_ := assert.New(suite.T())
	configName := "app_test_bad_config.env"
	srv, osinSrv, ctx := helpers.StartServers(configName)
	defer helpers.ShutdownServers(ctx, srv, osinSrv)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/auth", viper.GetInt("SERVICE_PORT")))
	assert_.Nil(err, "Error authorize")
	assert_.Equal(404, resp.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Debug().Msg(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	assert_.Contains(string(body), "404 page not found")
}

func (suite *BadConfigTestSuite) TestBadAuthRedirect() {
	if os.Getenv("APITEST") == "" {
		suite.T().Skip("skipping test; $APITEST not set")
	}
	assert_ := assert.New(suite.T())
	helpers.InitConfigForTests("app_test_bad_config.env")
	router := avanpost_auth.SetupRouter(false)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/auth", nil)
	assert_.Nil(err, "Error get config")
	router.ServeHTTP(w, req)

	assert_.Equal(302, w.Code)
	expected := fmt.Sprintf("<a href=\"%s\">Found</a>.\n\n", utils.FullUrlForAuthorize("xyz"))
	expected = strings.ReplaceAll(expected, "&", "&amp;")
	assert_.Equal(expected, w.Body.String())
}
