package utils

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"math/rand"
	"net/url"
	"runtime"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type CallerHook struct{}

func (h CallerHook) Run(event *zerolog.Event, _ zerolog.Level, _ string) {
	if pc, _, _, ok := runtime.Caller(3); ok {
		details := runtime.FuncForPC(pc)
		name := "???"
		if ok && details != nil {
			name = details.Name()
		}
		event.Str("fn", name[strings.LastIndex(name, "/")+1:])
	}
}

func FullUrlForAuthorize(state string) string {
	var v = make(url.Values)
	v.Set("response_type", "code")
	v.Set("redirect_uri", RedirectUrl(viper.GetString("SERVICE_OAUTH2_REDIRECT")))
	v.Set("client_id", viper.GetString("OAUTH2_CLIENT_ID"))
	v.Set("state", state)
	v.Set("scope", "everything")
	var u = url.URL{
		Scheme:   viper.GetString("OAUTH2_URL_AUTH_SHEMA"),
		Host:     fmt.Sprintf("%s:%d", viper.GetString("OAUTH2_URL_AUTH_HOST"), viper.GetInt("OAUTH2_URL_AUTH_PORT")),
		Path:     viper.GetString("OAUTH2_URL_AUTH_PATH"),
		RawQuery: v.Encode(),
	}
	return u.String()
}

func RedirectUrl(redirectPath string) string {
	uRedirect := url.URL{
		Scheme: viper.GetString("SERVICE_SHEMA"),
		Host:   fmt.Sprintf("%s:%d", viper.GetString("SERVICE_HOST"), viper.GetInt("SERVICE_PORT")),
		Path:   redirectPath,
	}
	return uRedirect.String()
}
