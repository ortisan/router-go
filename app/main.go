package main

import (
	"context"
	"time"

	"github.com/ortisan/router-go/internal/api"
	"github.com/ortisan/router-go/internal/config"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/ortisan/router-go/internal/loadbalancer"
	"github.com/ortisan/router-go/internal/telemetry"
	"github.com/rs/zerolog"
)

// @title Router API
// @version 2.0
// @description This is an Router APi that balance requests to healthy service endpoints.
// @termsOfService http://swagger.io/terms/

// @contact.name Marcelo
// @contact.url https://github.com/ortisan
// @contact.email tentativafc@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Config telemetry
	tp, err := telemetry.Setup()
	if err != nil {
		panic(errApp.NewGenericError("Error to setup telemetry", err))
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Graceful shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(errApp.NewGenericError("Error shutdown telemetry", err))

		}
	}(ctx)

	// Config load balancer
	if err := loadbalancer.Setup(); err != nil {
		panic(errApp.NewGenericError("Error to setup loadbalancer", err))
	}

	// Config server and routes
	r := api.Setup()

	// Running server
	r.Run(config.ConfigObj.App.ServerAddress) // Listen server
}
