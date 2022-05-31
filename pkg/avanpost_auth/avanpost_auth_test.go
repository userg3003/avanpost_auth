package avanpost_auth_test

import (
	"avanpost_auth/internal/helpers"
	"avanpost_auth/pkg/avanpost_auth"
	"avanpost_auth/pkg/avanpost_auth/constants"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/url"
	"testing"
)

type ApiTestSuite struct {
	suite.Suite
}

func TestServiceTestSuite(t *testing.T) {
	configName := "app_test.env"
	srv, osinSrv, ctx := helpers.StartServers(configName)
	suite.Run(t, new(ApiTestSuite))
	helpers.ShutdownServers(ctx, srv, osinSrv)
}

// Тестирование авторизации через OAuth2-сервер
func (suite *ApiTestSuite) TestAuthorize() {
	//if os.Getenv("APITEST") == "" && os.Getenv("ALLTESTS") != "1" {
	//	suite.T().Skip("skipping test; $APITEST not set")
	//}
	assert_ := assert.New(suite.T())
	pw, err, browser, page := helpers.InitPlaywright(assert_)
	helpers.InitConfigForTests("app_test.env")
	authUrl := url.URL{
		Scheme: viper.GetString("SERVICE_SHEMA"),
		Host:   fmt.Sprintf("%s:%d", viper.GetString("SERVICE_HOST"), viper.GetInt("SERVICE_PORT")),
		Path:   "auth",
	}
	_, err = page.Goto(authUrl.String())
	assert_.Nil(err, "could not goto "+authUrl.String())
	//if _, err = page.Screenshot(playwright.PageScreenshotOptions{
	//	Path: playwright.String("OAuth2_1.png"),
	//}); err != nil {
	//	log_.Fatalf("could not create screenshot: %v", err)
	//}
	entries, err := page.QuerySelectorAll("input")
	assert_.Nil(err, "could not get entries")
	assert_.Equal(3, len(entries), "Not found input fields")
	err = entries[0].Type("test")
	assert_.Nil(err, "could not set login entries")
	err = entries[1].Type("test")
	assert_.Nil(err, "could not set password entries")

	err = entries[2].Press("Enter")
	assert_.Nil(err, "could not press Enter")
	content, err := page.Content()
	assert_.Contains(content, ">ok<")
	assert_.Nil(err, "could not get page content")
	pageUrl := page.URL()
	contextsPage := page.Context()
	cookies, err := contextsPage.Cookies(pageUrl)
	assert_.Equal(viper.GetString("SERVICE_COOKIE_SESSION_NAME"), cookies[0].Name)
	log.Debug().Msg(cookies[0].Name)
	helpers.StopPlaywright(err, browser, assert_, pw)
}

func (suite *ApiTestSuite) TestGetConfig() {
	//if os.Getenv("APITEST") == "" && os.Getenv("ALLTESTS") != "1" {
	//	suite.T().Skip("skipping test; $SRVTEST not set")
	//}
	assert_ := assert.New(suite.T())

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/config", viper.GetInt("SERVICE_PORT")))
	assert_.Nil(err, "Error get config")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		require.Nil(suite.T(), err)
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	assert_.Nil(err, "Error get config")
	fmt.Printf("body ->:  %s", body)
	var configJson helpers.ConfigJson
	err = json.Unmarshal(body, &configJson)
	assert_.Nil(err, "Error unmarshal response body")
	assert_.Equal(configJson.Oauth2ClientSecret, viper.GetString("OAUTH2_CLIENT_SECRET"))
}

func (suite *ApiTestSuite) TestStartDuplicatedService() {
	assert_ := assert.New(suite.T())
	cmdStartService := func(in string) *cobra.Command {
		return &cobra.Command{
			Use:   "start",
			Short: "Start Avanpost auth service",
			Long:  constants.ServiceName + ` CLI сервис авторизации в Avanpost_FAM`,
			RunE:  avanpost_auth.StartService,
		}
	}

	err := cmdStartService("").Execute()
	assert_.NotNil(err)
	assert_.Equal("listen tcp 127.0.0.1:3011: bind: address already in use", err.Error())
}
