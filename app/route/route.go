package route

import (
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ortisan/router-go/config"
	domain "github.com/ortisan/router-go/domain"
	"github.com/ortisan/router-go/integration"
	"github.com/ortisan/router-go/util"
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

func Redirect(c *gin.Context) {
	resource := c.Param("resource")
	apiPaths := strings.Split(resource, "/")
	if len(apiPaths) < 2 {
		c.JSON(http.StatusBadRequest, domain.Error{Message: "Router can't process this request. Format of url must be /{prefix api}/{all_rest}"})
	}

	prefixService := apiPaths[1] // in url "http://xpto.com/api1/xpto", gets the "api1" value

	serviceApi, err := integration.GetValue(PrefixServicesConfig + prefixService)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to load service config.")
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.(error).Error(), StackTrace: string(debug.Stack())})
		return
	}

	redirectResource := serviceApi + util.GetSubstringAfter(resource, prefixService)
	method := c.Request.Method
	headers := c.Request.Header

	client := &http.Client{}
	req, err := http.NewRequest(method, redirectResource, c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to create request.")
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.(error).Error(), StackTrace: string(debug.Stack())})
		return
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
		log.Error().Stack().Err(err).Msg("Error to call API")
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.(error).Error(), StackTrace: string(debug.Stack())})
		return
	}

	defer resp.Body.Close() // Defer will close after this function ends
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to read response body")
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.(error).Error(), StackTrace: string(debug.Stack())})
		return
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body) // Data is returned
}

func ConfigServer(config config.Config) {
	r := gin.Default()

	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // Prometheus metrics
	r.GET("/api/*resource", Redirect)                // By Pass
	r.POST("/api/*resource", Redirect)               // By Pass
	r.PUT("/api/*resource", Redirect)                // By Pass
	r.PATCH("/api/*resource", Redirect)              // By Pass
	r.DELETE("/api/*resource", Redirect)             // By Pass

	r.Run(config.ServerAddress) // Listen server
}
