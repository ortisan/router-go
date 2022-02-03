package route

import (
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ortisan/router-go/internal/config"
	"github.com/ortisan/router-go/internal/domain"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/ortisan/router-go/internal/integration"
	"github.com/ortisan/router-go/internal/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const (
	UrlApp               = "https://jsonplaceholder.typicode.com"
	PrefixServicesConfig = "services.prefix."
)

func HeadersDisabledInRedirection() func(string) bool {
	innerMap := map[string]int{
		"Accept-Encoding": 1, // This header transform encodings
	}
	return func(key string) bool {

		_, found := innerMap[key]
		return found
	}
}

func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case errApp.IWithMessageAndStatusCode:
					error := err.(errApp.IWithMessageAndStatusCode)
					errObj := domain.Error{Message: error.Error(), Cause: error.Cause(), StackTrace: error.StackTrace()}
					c.Data(error.Status(), c.GetHeader("Content-Type"), util.ObjectToJson(errObj)) // Data is returned
				default:
					errObj := domain.Error{Message: err.(error).Error(), StackTrace: string(debug.Stack())}
					c.Data(http.StatusInternalServerError, c.GetHeader("Content-Type"), util.ObjectToJson(errObj)) // Data is returned
				}
			}
		}()
		c.Next()
	}
}

func Redirect(c *gin.Context) {

	//traceid := uuid.New()

	resource := c.Param("resource")
	apiPaths := strings.Split(resource, "/")
	if len(apiPaths) < 2 {
		c.JSON(http.StatusBadRequest, domain.Error{Message: "Router can't process this request. Format of url must be /{prefix api}/{all_rest}"})
	}

	prefixService := apiPaths[1] // in url "http://xpto.com/api1/xpto", gets the "api1" value

	serviceApi, err := integration.GetValue(PrefixServicesConfig + prefixService)

	if err != nil {
		panic(errApp.NewIntegrationError("Error to load service config.", err))
	}

	redirectResource := serviceApi + util.GetSubstringAfter(resource, prefixService)
	method := c.Request.Method
	headers := c.Request.Header

	client := &http.Client{}
	req, err := http.NewRequest(method, redirectResource, c.Request.Body)
	if err != nil {
		panic(errApp.NewGenericError("Error to create request", err))
	}

	//Adding headers
	for name, values := range headers {
		for _, value := range values {
			log.Debug().Str(name, value).Msg("Iterating headers...")
			found := HeadersDisabledInRedirection()(name)
			if !found {
				req.Header.Set(name, value)
			}
		}
	}

	resp, err := client.Do(req) // Call API
	if err != nil {
		panic(errApp.NewIntegrationError("Error to call API.", err))
	}

	defer resp.Body.Close() // Defer will close after this function ends
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errApp.NewIntegrationError("Error read response body.", err))
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body) // Data is returned
}

func ConfigServer() {
	r := gin.Default()

	// Middlewares
	r.Use(ErrorHandle()) // Error handling
	r.Use(gin.Logger())  // Logger request/response

	// Routes
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // Prometheus metrics
	r.GET("/api/*resource", Redirect)                // By Pass
	r.POST("/api/*resource", Redirect)               // By Pass
	r.PUT("/api/*resource", Redirect)                // By Pass
	r.PATCH("/api/*resource", Redirect)              // By Pass
	r.DELETE("/api/*resource", Redirect)             // By Pass

	// Running server
	r.Run(config.ConfigObj.App.ServerAddress) // Listen server
}
