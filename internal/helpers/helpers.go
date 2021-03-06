package helpers

import (
	"avanpost_auth/pkg/avanpost_auth"
	"avanpost_auth/pkg/avanpost_auth/config"
	"avanpost_auth/pkg/osin_server"
	"context"
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"time"
)

func StartOauth2Server() *http.Server {
	mux := osin_server.InitOsinServer()
	srv := &http.Server{
		Addr:    ":14000",
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen OAuth2 server: %s\n", err)
		}
	}()
	return srv

}

func StartService(configName string) *http.Server {
	InitConfigForTests(configName)
	router := avanpost_auth.SetupRouter(false)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("SERVICE_PORT")),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen Avanpost-service: %s\n", err)
		}

	}()
	return srv

}

func StartServers(configName string) (*http.Server, *http.Server, context.Context) {
	srv := StartService(configName)
	osinSrv := StartOauth2Server()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return srv, osinSrv, ctx
}

func ShutdownServers(ctx context.Context, srv *http.Server, osinSrv *http.Server) {
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf(fmt.Sprintf("Service forced to shutdown: %v", err))
	}
	if err := osinSrv.Shutdown(ctx); err != nil {
		log.Printf(fmt.Sprintf("Server forced to shutdown: %v", err))
	}
}

func InitConfigForTests(fileName string) {
	config.ConfigFile = fileName
	config.ConfigPath = "../../config"
	config.InitConfig()
}

func StopPlaywright(err error, browser playwright.Browser, assert *assert.Assertions, pw *playwright.Playwright) {
	err = browser.Close()
	assert.Nil(err, fmt.Sprintf("could not close browser: %v", err))

	err = pw.Stop()
	assert.Nil(err, fmt.Sprintf("could not stop Playwright: %v", err))
}

func InitPlaywright(assert *assert.Assertions) (*playwright.Playwright, error, playwright.Browser, playwright.Page) {
	pw, err := playwright.Run()
	assert.Nil(err, fmt.Sprintf("could not start playwright: %v", err))
	browser, err := pw.Chromium.Launch()
	assert.Nil(err, fmt.Sprintf("could not launch browser: %v", err))
	page, err := browser.NewPage()
	assert.Nil(err, fmt.Sprintf("could not create page: %v", err))
	return pw, err, browser, page
}

type ErrAuthorize struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	State            string `json:"state"`
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
