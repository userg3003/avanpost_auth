package avanpost_auth

import (
	"avanpost_auth/docs"
	"avanpost_auth/internal/utils"
	config2 "avanpost_auth/pkg/avanpost_auth/config"
	constants "avanpost_auth/pkg/avanpost_auth/constants"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/ioutil"
	"net/http"
	"net/url"
)

var RootCmdAuth = &cobra.Command{
	Use:   "start",
	Short: "Start Avanpost auth service",
	Long:  constants.ServiceName + ` CLI сервис авторизации в Avanpost_FAM`,
	RunE:  StartService,
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger().Hook(utils.CallerHook{})
	log.Debug().Msg("")
	cobra.OnInitialize(config2.InitConfig)

	RootCmdAuth.PersistentFlags().StringVarP(&config2.ConfigFile, "file", "f", constants.ServiceName+".env",
		"configuration file, default is "+constants.ServiceName+".env")
	RootCmdAuth.PersistentFlags().StringVarP(&config2.ConfigPath, "config", "c", ".",
		"configuration path, default is "+".")
	err := viper.BindPFlag("useViper", RootCmdAuth.PersistentFlags().Lookup("viper"))
	log.Error().Msg(err.Error())
}

// PingExample godoc
// @Summary  Получение информации о конфигурации сервиса
// @Schemes
// @Description  Получить конфигурацию
// @Tags     Debug
// @Accept   json
// @Produce      json
// @Success      200  {string}  Json
// @Router       /config [get]
func config(g *gin.Context) {
	log.Debug().Msg("")
	g.JSON(http.StatusOK, viper.AllSettings())
}

// AuthExample godoc
// @Summary  Авторизация на сервере Avanpost_FAM
// @Schemes
// @Description  Выполнить авторизацию
// @Description  При обращении по данному пути в браузере активируется окно для ввода
// @Description  регистрационных данных пользователя. При успешной авторизации выполняется перенаправление
// @Description  по пути указанному в разделе Internal.
// @Description
// @Tags     Auth
// @Accept   json
// @Produce  html
// @Success  200  {string}  Auth
// @Router   /auth [get]
func auth(c *gin.Context) {
	state := "xyz" //utils.RandStringRunes(10)
	u := utils.FullUrlForAuthorize(state)
	c.Redirect(http.StatusFound, u)
}

// AuthExample godoc
// @Summary  Авторизация на сервере Avanpost_FAM.
// @Schemes
// @Description  По этому пути "прилитает" ответ от сервера Avanpost_FAM с данными авторизации.
// @Description  Результат авторизации возвращается в сессионной cookie.
// @Tags         Internal
// @Accept       json
// @Produce  html
// @Success      200  {string}  Auth
// @Router       /appauth [get]
func appauth(c *gin.Context) {
	log.Debug().Msg("appauth")
	allParam := c.Request.URL.Query().Encode()
	log.Debug().Msg(allParam)
	for k, v := range c.Request.URL.Query() {
		fmt.Printf("Key: %q Values: %q\n", k, v)
	}
	code := c.DefaultQuery("code", "Guest")
	state := c.Query("state")
	var u = url.URL{
		Scheme: viper.GetString("OAUTH2_URL_AUTH_SHEMA"),
		Host:   fmt.Sprintf("%s:%d", viper.GetString("OAUTH2_URL_AUTH_HOST"), viper.GetInt("OAUTH2_URL_AUTH_PORT")),
		Path:   viper.GetString("OAUTH2_URL_TOKEN_PATH"),
	}
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {viper.GetString("OAUTH2_CLIENT_ID")},
		"code":          {code},
		"client_secret": {viper.GetString("OAUTH2_CLIENT_SECRET")},
		"state":         {state},
		"redirect_uri":  {utils.RedirectUrl(viper.GetString("SERVICE_OAUTH2_REDIRECT"))},
	}

	resp, err := http.PostForm(u.String(), data)

	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()
	var res map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	session := sessions.Default(c)

	session.Set(constants.ServiceSessionParamName, res)
	err = session.Save()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	if viper.GetString("SERVICE_REDIRECT_URL_AFTER_AUTH") != "" {
		c.Redirect(http.StatusFound, utils.RedirectUrl(viper.GetString("SERVICE_REDIRECT_URL_AFTER_AUTH")))
	} else {
		c.String(http.StatusOK, "ok")
	}
}

// AuthExample godoc
// @Summary  Получение авторизацонных данных от сервера Avanpost_FAM
// @Schemes
// @Description  Получить данные сервера Avanpost_FAM по токену.
// @Description  Токен передаётся в сессионной cookie.
// @Description  Результат возвращается в ... (см. ReadMe).
// @Tags         Auth
// @Accept       json
// @Produce      html
// @Success      200  {string}  Auth
// @Router       /info [get]
func info(c *gin.Context) {
	session := sessions.Default(c)
	res := session.Get(constants.ServiceSessionParamName)
	accessToken := res.(map[string]interface{})["access_token"]
	fmt.Printf("Result: %+v", accessToken)

	var u = url.URL{
		Scheme: viper.GetString("OAUTH2_URL_AUTH_SHEMA"),
		Host:   fmt.Sprintf("%s:%d", viper.GetString("OAUTH2_URL_AUTH_HOST"), viper.GetInt("OAUTH2_URL_AUTH_PORT")),
		Path:   viper.GetString("OAUTH2_URL_INFO_PATH"),
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+accessToken.(string))
	response, errGet := client.Do(req)
	if errGet != nil {
		log.Error().Msg(errGet.Error())
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()
	body, errRead := ioutil.ReadAll(response.Body)

	if errRead != nil {
		log.Error().Msg(errRead.Error())
	}

	log.Debug().Msg(string(body))

	//c.JSON(http.StatusOK, body)
	c.String(http.StatusOK, string(body))
}

// AuthExample godoc
// @Summary  Страница после успешной авторизации
// @Schemes
// @Description  Страница после успешной авторизации
// @Description  Реальный путь до конечной точеи задаётся в конфигурвции сервиса.
// @Description
// @Tags         Debug
// @Accept       json
// @Produce      html
// @Success  200  {string}  goodauth
// @Router   /goodauth [get]
func goodauth(c *gin.Context) {
	session := sessions.Default(c)
	jsonStr := session.Get("currentJson")
	res := session.Get("currentRes")

	fmt.Println(jsonStr)
	fmt.Println(res)

	//c.JSON(http.StatusOK, res)
	c.String(http.StatusOK, "ok")
}

// AuthExample godoc
// @Summary  Показать состояние сервера.
// @Schemes
// @Description  Получить статус сервиса.
// @Tags         Info
// @Accept       json
// @Produce      html
// @Success      204  {string}  health
// @Router       /health [get]
func health(c *gin.Context) {
	c.String(http.StatusNoContent, "")
}

func SetupRouter(swagger bool) *gin.Engine {
	//r := gin.Default()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	store := cookie.NewStore([]byte(viper.GetString("SERVICE_COOKIE_SESSION_SECRET")))
	r.Use(sessions.Sessions(viper.GetString("SERVICE_COOKIE_SESSION_NAME"), store))

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/auth", auth)
	r.GET(viper.GetString("SERVICE_OAUTH2_REDIRECT"), appauth)
	r.GET("/info", info)
	r.GET("/health", health)
	r.GET("/config", config)
	if viper.GetString("SERVICE_REDIRECT_URL_AFTER_AUTH") != "" {
		r.GET(viper.GetString("SERVICE_REDIRECT_URL_AFTER_AUTH"), goodauth)
	}
	if swagger {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return r
}

// @title           Avanpost auth Swagger API
// @version         1.0
// @description     Swagger API for Golang Project Avanpost_auth.
// @termsOfService  http://swagger.io/terms/

// @host                  localhost:3011
// @contact.name          API Support
// @contact.email         s.urvanov@gmail.com

// @BasePath  /

func StartService(_ *cobra.Command, _ []string) error {
	r := SetupRouter(viper.GetBool("SWAGGER"))
	err := r.Run(fmt.Sprintf("%s:%d", viper.GetString("SERVICE_HOST"), viper.GetInt("SERVICE_PORT")))
	log.Error().Msg(err.Error())
	return err
}
