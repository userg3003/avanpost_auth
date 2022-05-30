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
	osin_server.InitOsinServer()
	srv := &http.Server{
		Addr:    ":14000",
		Handler: nil,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
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
			log.Printf("listen: %s\n", err)
		}

	}()

	return srv

}

func StartServers(configName string) (*http.Server, *http.Server, context.Context) {
	var srv *http.Server = StartService(configName)
	osinSrv := StartOauth2Server()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	return srv, osinSrv, ctx
}

func ShutdownServers(ctx context.Context, srv *http.Server, osinSrv *http.Server) {
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf(fmt.Sprintf("Server forced to shutdown: %v", err))
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
