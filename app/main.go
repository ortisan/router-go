package main

import (
	"context"
	"time"

	"github.com/ortisan/router-go/internal/api"
	"github.com/ortisan/router-go/internal/integration"
	"github.com/ortisan/router-go/internal/loadbalancer"
	"github.com/ortisan/router-go/internal/telemetry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	// Mock values
	integration.PutValue("/services/prefix/app1", "https://jsonplaceholder.typicode.com,https://abcxpto.com") // TODO Just for testing

	values, _ := integration.GetValues("/services/prefix/app1")        // TODO Just for testing
	log.Debug().Strs("values", values).Msg("Values from etcd loaded.") // TODO Just for testing

	mapValues, _ := integration.GetValuesPrefixed("/services/prefix/") // TODO Just for testing
	for key, value := range mapValues {                                // TODO Just for testing
		log.Debug().Str(key, value).Msg("Values from etcd loaded.") // TODO Just for testing
	}

	var res, err = integration.PutCacheValue("teste", "1")
	if err != nil {
		panic(err)
	}
	log.Debug().Str("result", res).Msg("Put value on cache.") // TODO Just for testing

	res, err = integration.GetCacheValue("teste")
	if err != nil {
		panic(err)
	}
	log.Debug().Str("result", res).Msg("Get value on cache.") // TODO Just for testing

	// Config telemetry
	tp, err := telemetry.ConfigTracerProvider()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Graceful shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	// Config load balancer
	loadbalancer.ConfigLoadBalancer()

	// Config server and routes
	api.ConfigServer()
}
