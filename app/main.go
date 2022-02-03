package main

import (
	"github.com/ortisan/router-go/internal/integration"
	"github.com/ortisan/router-go/internal/route"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	integration.PutValue("services.prefix.app1", "https://jsonplaceholder.typicode.com") // TODO Just for testing
	value, _ := integration.GetValue("services.prefix.app1")                             // TODO Just for testing
	log.Debug().Msg(value)                                                               // TODO Just for testing

	route.ConfigServer()
}
