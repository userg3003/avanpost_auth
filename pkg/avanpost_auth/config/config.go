package config

import (
	"avanpost_auth/pkg/avanpost_auth/constants"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	ConfigPath string
	ConfigFile string
)

func InitConfig() {
	log.Info().Msg("initConfig")

	viper.SetDefault("SERVICE_SHEMA", "http")
	viper.SetDefault("SERVICE_HOST", "localhost")
	viper.SetDefault("SERVICE_PORT", 3011)
	viper.SetDefault("SERVICE_OAUTH2_REDIRECT", "appauth")
	viper.SetDefault("SERVICE_COOKIE_SESSION_NAME", "avanpostAuth")
	viper.SetDefault("SERVICE_COOKIE_SESSION_SECRET", "12345")
	viper.SetDefault("SERVICE_REDIRECT_URL_AFTER_AUTH", "goodauth")
	viper.SetDefault("OAUTH2_URL_AUTH_SHEMA", "http")
	viper.SetDefault("OAUTH2_URL_AUTH_HOST", "localhost")
	viper.SetDefault("OAUTH2_URL_AUTH_PORT", 14000)
	viper.SetDefault("OAUTH2_URL_AUTH_PATH", "authorize")
	viper.SetDefault("OAUTH2_URL_TOKEN_PATH", "token")
	viper.SetDefault("OAUTH2_URL_INFO_PATH", "info")
	viper.SetDefault("OAUTH2_CLIENT_ID", 1234)
	viper.SetDefault("OAUTH2_CLIENT_SECRET", "aabbccdd")

	viper.SetDefault("SWAGGER", true)
	viper.SetDefault("LOG_LEVEL", "info")
	logLevel, _ := zerolog.ParseLevel(viper.GetString("LOG_LEVEL"))
	zerolog.SetGlobalLevel(logLevel)
	viper.SetConfigType("env")
	if ConfigFile != "" {
		viper.SetConfigName(ConfigFile) // Use config file from the flag.
	} else {
		viper.SetConfigName(constants.ServiceName)
	}

	if ConfigPath != "" {
		viper.AddConfigPath(ConfigPath)
	} else {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(cwd)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if err1, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Info().Msg(fmt.Sprintf("Using config file: %s, err: %v", viper.ConfigFileUsed(), err1))
		} else {
			// Config file was found but another error was produced
			log.Info().Msg(fmt.Sprintf("Using config file: %s, err: %v", viper.ConfigFileUsed(), err1))
		}
	} else {
		log.Info().Msg(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
	logLevel, _ = zerolog.ParseLevel(viper.GetString("LOG_LEVEL"))
	zerolog.SetGlobalLevel(logLevel)
}
