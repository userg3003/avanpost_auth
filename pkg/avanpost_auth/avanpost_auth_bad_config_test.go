package avanpost_auth_test

import (
	"avanpost_auth/internal/helpers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"os"
	"testing"
)

type BadConfigTestSuite struct {
	suite.Suite
}

func TestBadConfigTestSuite(t *testing.T) {
	helpers.InitConfigForTests("app_test.env")
	suite.Run(t, new(BadConfigTestSuite))
}

// Проверить перенаправление на автоизацию по несуществующему пути на OAuth2-сервере
func (suite *BadConfigTestSuite) TestBadOAut2AuthorizePath() {
	//if os.Getenv("APITEST") == "" {
	//	suite.T().Skip("skipping test; $APITEST not set")
	//}
	assert_ := assert.New(suite.T())
	configName := "app_test.env"
	configParam := os.Getenv("OAUTH2_URL_AUTH_PATH")
	tmpConfig, srv, osinSrv, ctx := startServers(configParam, configName)
	defer func() {
		os.Setenv(configParam, tmpConfig)
		helpers.ShutdownServers(ctx, srv, osinSrv)
	}()

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

// Проверить несовпадение  id-client сервиса и  OAuth2-сервера
func (suite *BadConfigTestSuite) TestBadOAut2ClientId() {
	assert_ := assert.New(suite.T())
	configName := "app_test.env"
	configParam := "OAUTH2_URL_AUTH_PATH"
	tmpConfig, srv, osinSrv, ctx := startServers(configParam, configName)
	defer func() {
		os.Setenv(configParam, tmpConfig)
		helpers.ShutdownServers(ctx, srv, osinSrv)
	}()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/auth", viper.GetInt("SERVICE_PORT")))
	assert_.Nil(err, "Error authorize")
	assert_.Equal(200, resp.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Debug().Msg(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	var errAuthorize helpers.ErrAuthorize
	err = json.Unmarshal(body, &errAuthorize)
	assert_.Nil(err, "Error unmarshal response body")
	assert_.Equal(errAuthorize.Error, "unauthorized_client")
}

func startServers(configParam string, configName string) (string, *http.Server, *http.Server, context.Context) {
	tmpConfig := os.Getenv(configParam)
	os.Setenv(configParam, os.Getenv(configParam)+"123")
	srv, osinSrv, ctx := helpers.StartServers(configName)
	return tmpConfig, srv, osinSrv, ctx
}

// Проверить несовпадение  секрета сервиса и  OAuth2-сервера
// При несовпадении OAuth2-сервер пренаправляет на страницу ввода логина и пароля.
func (suite *BadConfigTestSuite) TestBadOAut2ClientSecret() {
	assert_ := assert.New(suite.T())
	configName := "app_test.env"
	os.Setenv("OAUTH2_CLIENT_SECRET", "123")
	srv, osinSrv, ctx := helpers.StartServers(configName)

	defer func() {
		helpers.InitConfigForTests(configName)
		helpers.ShutdownServers(ctx, srv, osinSrv)
	}()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/auth", viper.GetInt("SERVICE_PORT")))
	assert_.Nil(err, "Error authorize")
	assert_.Equal(200, resp.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Debug().Msg(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	assert_.Contains(string(body), "<body>LOGIN 1234 (use test/test)<br/><form")
}

// Проверить несовпадение  пути редиректа авторизации сервиса   и  OAuth2-сервера
func (suite *BadConfigTestSuite) TestBadServiceOAuth2Redirect() {
	assert_ := assert.New(suite.T())
	configName := "app_test.env"
	os.Setenv("SERVICE_OAUTH2_REDIRECT", "badredirect")
	srv, osinSrv, ctx := helpers.StartServers(configName)
	defer func() {
		helpers.InitConfigForTests(configName)
		helpers.ShutdownServers(ctx, srv, osinSrv)
	}()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/auth", viper.GetInt("SERVICE_PORT")))
	assert_.Nil(err, "Error authorize")
	assert_.Equal(200, resp.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Debug().Msg(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	var errAuthorize helpers.ErrAuthorize
	err = json.Unmarshal(body, &errAuthorize)
	assert_.Nil(err, "Error unmarshal response body")
	assert_.Equal(errAuthorize.Error, "invalid_request")
	assert_.Contains(errAuthorize.ErrorDescription, "The request is missing a required paramete")
}
