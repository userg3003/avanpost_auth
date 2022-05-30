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

func (suite *ApiTestSuite) TestAuthorize() {
	//if os.Getenv("APITEST") == "" && os.Getenv("ALLTESTS") != "1" {
	//	suite.T().Skip("skipping test; $APITEST not set")
	//}
	assert := assert.New(suite.T())
	pw, err, browser, page := helpers.InitPlaywright(assert)
	helpers.InitConfigForTests("app_test.env")
	authUrl := url.URL{
		Scheme: viper.GetString("SERVICE_SHEMA"),
		Host:   fmt.Sprintf("%s:%d", viper.GetString("SERVICE_HOST"), viper.GetInt("SERVICE_PORT")),
		Path:   "auth",
	}
	_, err = page.Goto(authUrl.String())
	assert.Nil(err, "could not goto "+authUrl.String())
	//if _, err = page.Screenshot(playwright.PageScreenshotOptions{
	//	Path: playwright.String("OAuth2_1.png"),
	//}); err != nil {
	//	log_.Fatalf("could not create screenshot: %v", err)
	//}
	entries, err := page.QuerySelectorAll("input")
	assert.Nil(err, "could not get entries")
	assert.Equal(3, len(entries), "Not found input fields")
	err = entries[0].Type("test")
	assert.Nil(err, "could not set login entries")
	entries[1].Type("test")
	assert.Nil(err, "could not set password entries")

	err = entries[2].Press("Enter")
	assert.Nil(err, "could not press Enter")
	content, err := page.Content()
	assert.Contains(content, ">ok<")
	assert.Nil(err, "could not get page content")
	pageUrl := page.URL()
	contextsPage := page.Context()
	cookies, err := contextsPage.Cookies(pageUrl)
	assert.Equal(viper.GetString("SERVICE_COOKIE_SESSION_NAME"), cookies[0].Name)
	log.Debug().Msg(cookies[0].Name)
	helpers.StopPlaywright(err, browser, assert, pw)
}

type ConfigJson struct {
	LogLevel                    string `json:"log_level"`
	Oauth2ClientID              string `json:"oauth2_client_id"`
	Oauth2ClientSecret          string `json:"oauth2_client_secret"`
	Oauth2URLAuthHost           string `json:"oauth2_url_auth_host"`
	Oauth2URLAuthPath           string `json:"oauth2_url_auth_path"`
	Oauth2URLAuthPort           string `json:"oauth2_url_auth_port"`
	Oauth2URLAuthShema          string `json:"oauth2_url_auth_shema"`
	Oauth2URLInfoPath           string `json:"oauth2_url_info_path"`
	Oauth2URLTokenPath          string `json:"oauth2_url_token_path"`
	ServiceCookieSessionName    string `json:"service_cookie_session_name"`
	ServiceCookieSessionSecret  string `json:"service_cookie_session_secret"`
	ServiceHost                 string `json:"service_host"`
	ServiceOauth2Redirect       string `json:"service_oauth2_redirect"`
	ServicePort                 string `json:"service_port"`
	ServiceRedirectURLAfterAuth string `json:"service_redirect_url_after_auth"`
	ServiceShema                string `json:"service_shema"`
	Swagger                     string `json:"swagger"`
}

func (suite *ApiTestSuite) TestGetConfig() {
	//if os.Getenv("APITEST") == "" && os.Getenv("ALLTESTS") != "1" {
	//	suite.T().Skip("skipping test; $SRVTEST not set")
	//}
	assert := assert.New(suite.T())

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/config", viper.GetInt("SERVICE_PORT")))
	assert.Nil(err, "Error get config")
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.Nil(err, "Error get config")
	fmt.Printf("body ->:  %s", body)
	var configJson ConfigJson
	json.Unmarshal(body, &configJson)
	assert.Equal(configJson.Oauth2ClientSecret, viper.GetString("OAUTH2_CLIENT_SECRET"))
}

func (suite *ApiTestSuite) TestStartDuplicatedService() {
	assert := assert.New(suite.T())
	cmdStartService := func(in string) *cobra.Command {
		return &cobra.Command{
			Use:   "start",
			Short: "Start Avanpost auth service",
			Long:  constants.ServiceName + ` CLI сервис авторизации в Avanpost_FAM`,
			RunE:  avanpost_auth.StartService,
		}
	}

	err := cmdStartService("").Execute()
	assert.NotNil(err)
	assert.Equal("listen tcp 127.0.0.1:3011: bind: address already in use", err.Error())
}
