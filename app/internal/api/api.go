package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/ortisan/router-go/docs"
	"github.com/ortisan/router-go/internal/config"
	"github.com/ortisan/router-go/internal/constant"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/ortisan/router-go/internal/loadbalancer"
	"github.com/ortisan/router-go/internal/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case errApp.IWithMessageAndStatusCode:
					error := err.(errApp.IWithMessageAndStatusCode)
					errObj := errApp.Error{Message: error.Error(), Cause: error.Cause(), StackTrace: error.StackTrace()}
					errJsonBytes, _ := util.ObjectToJsonBytes(errObj)
					c.Data(error.Status(), c.GetHeader(constant.ContentTypeHeaderName), errJsonBytes) // Data is returned
				default:
					errObj := errApp.Error{Message: err.(error).Error(), StackTrace: string(debug.Stack())}
					errJsonBytes, _ := util.ObjectToJsonBytes(errObj)
					c.Data(http.StatusInternalServerError, c.GetHeader(constant.ContentTypeHeaderName), errJsonBytes) // Data is returned
				}
			}
		}()
		c.Next()
	}
}

// Healthcheck
// @Summary Health check service
// @Description Health check service
// @Tags router healthcheck
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"status": "up",
	}
	c.JSON(http.StatusOK, res)
}

// Get the available server by prefix and redirect request
// @Summary Redirect request to healthy server
// @Description Redirect request.
// @Tags router redirect
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Success 201 {object} map[string]interface{}
// @Success 204 {object} map[string]interface{}
// @Router /api/{prefix_service}/{backend_api_service} [get]
// @Router /api/{prefix_service}/{backend_api_service} [post]
// @Router /api/{prefix_service}/{backend_api_service} [put]
// @Router /api/{prefix_service}/{backend_api_service} [patch]
// @Router /api/{prefix_service}/{backend_api_service} [delete]
func HandleRequest(c *gin.Context) {

	resource := c.Param("resource")
	apiPaths := strings.Split(resource, "/")

	r := c.Request

	if len(apiPaths) < 2 {
		panic(errApp.NewBadRequestErrorWithCause("Router can't process this request. Format of url must be /{prefix api}/{all_rest}", nil))
	}
	servicePrefix := apiPaths[1] // in url "http://xpto.com/api/api1/xpto", gets the "api1" value

	serverPool := loadbalancer.ServerPoolsObj.GetServerPoolByPrefix(servicePrefix)
	if serverPool == nil {
		panic(errApp.NewBadRequestError(fmt.Sprintf("Cannot any server that can handle the prefix \"%s\"", servicePrefix)))
	}

	retries := loadbalancer.GetRetryFromContext(r)

	if retries < loadbalancer.MaxRetries {
		select {
		case <-time.After(loadbalancer.BackoffTimeout):
			serverPool.HandleRequest(c, util.GetSubstringAfter(resource, servicePrefix), r.Method, r.Header)
		}
	}
}

func Setup() *gin.Engine {
	r := gin.Default()

	// Middlewares
	r.Use(otelgin.Middleware(config.ConfigObj.App.Name)) // Tracer
	r.Use(ErrorHandler())                                // Error handling
	r.Use(gin.Logger())                                  // Logger request/response

	// Routes
	r.GET("/", HealthCheck)                          // HealthCheck
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // Prometheus metrics
	r.GET("/api/*resource", HandleRequest)           // By Pass
	r.POST("/api/*resource", HandleRequest)          // By Pass
	r.PUT("/api/*resource", HandleRequest)           // By Pass
	r.PATCH("/api/*resource", HandleRequest)         // By Pass
	r.DELETE("/api/*resource", HandleRequest)        // By Pass

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	return r
}
