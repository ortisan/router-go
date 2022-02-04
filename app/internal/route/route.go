package route

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ortisan/router-go/internal/config"
	"github.com/ortisan/router-go/internal/constant"
	"github.com/ortisan/router-go/internal/domain"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/ortisan/router-go/internal/loadbalancer"
	"github.com/ortisan/router-go/internal/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ConfigTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			c.Set(constant.TraceIdHeaderName, uuid.New().String())
		}()
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case errApp.IWithMessageAndStatusCode:
					error := err.(errApp.IWithMessageAndStatusCode)
					errObj := domain.Error{Message: error.Error(), Cause: error.Cause(), StackTrace: error.StackTrace()}
					c.Data(error.Status(), c.GetHeader(constant.ContentTypeHeaderName), util.ObjectToJson(errObj)) // Data is returned
				default:
					errObj := domain.Error{Message: err.(error).Error(), StackTrace: string(debug.Stack())}
					c.Data(http.StatusInternalServerError, c.GetHeader(constant.ContentTypeHeaderName), util.ObjectToJson(errObj)) // Data is returned
				}
			}
		}()
		c.Next()
	}
}

func ConfigServer() {
	r := gin.Default()

	// Middlewares
	r.Use(ConfigTraceId()) // Config TraceId
	r.Use(ErrorHandler())  // Error handling
	r.Use(gin.Logger())    // Logger request/response

	// Routes
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))       // Prometheus metrics
	r.GET("/api/*resource", loadbalancer.HandleRequest)    // By Pass
	r.POST("/api/*resource", loadbalancer.HandleRequest)   // By Pass
	r.PUT("/api/*resource", loadbalancer.HandleRequest)    // By Pass
	r.PATCH("/api/*resource", loadbalancer.HandleRequest)  // By Pass
	r.DELETE("/api/*resource", loadbalancer.HandleRequest) // By Pass

	// Running server
	r.Run(config.ConfigObj.App.ServerAddress) // Listen server
}
