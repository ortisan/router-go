package main

import (
	"github.com/ortisan/router-go/internal/integration"
	"github.com/ortisan/router-go/internal/loadbalancer"
	"github.com/ortisan/router-go/internal/route"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	integration.PutValue("/services/prefix/app1", "https://jsonplaceholder.typicode.com,https://abcxpto.com") // TODO Just for testing

	// values, _ := integration.GetValues("/services/prefix/app1")                                                                    // TODO Just for testing
	// log.Debug().Strs("values", values).Msg("Values from etcd loaded.")                                                             // TODO Just for testing

	// mapValues, _ := integration.GetValuesPrefixed("/services/prefix/") // TODO Just for testing
	// for key, value := range mapValues {
	// 	log.Debug().Str(key, value).Msg("Values from etcd loaded.") // TODO Just for testing
	// }

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

	loadbalancer.ConfigLoadBalancer()

	route.ConfigServer()
}
