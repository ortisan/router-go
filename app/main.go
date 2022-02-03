package main

import (
	"github.com/ortisan/router-go/internal/integration"
	"github.com/ortisan/router-go/internal/loadbalancer"
	"github.com/ortisan/router-go/internal/route"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	integration.PutValue("/services/prefix/app1", "https://jsonplaceholder.typicode.com")                                          // TODO Just for testing
	integration.PutValue("/services/prefix/app2", "https://jsonplaceholder.typicode.com,https://jsonplaceholderxpto.typicode.com") // TODO Just for testing

	// values, _ := integration.GetValues("/services/prefix/app1")                                                                    // TODO Just for testing
	// log.Debug().Strs("values", values).Msg("Values from etcd loaded.")                                                             // TODO Just for testing

	// mapValues, _ := integration.GetValuesPrefixed("/services/prefix/") // TODO Just for testing
	// for key, value := range mapValues {
	// 	log.Debug().Str(key, value).Msg("Values from etcd loaded.") // TODO Just for testing
	// }

	loadbalancer.ConfigLoadBalancer()

	route.ConfigServer()
}
