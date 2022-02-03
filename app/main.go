package main

import (
	"github.com/ortisan/router-go/config"
	"github.com/ortisan/router-go/integration"
	"github.com/ortisan/router-go/route"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	config, err := config.LoadConfig(".")

	integration.PutValue("services.prefix.app1", "https://jsonplaceholder.typicode.com")
	value, _ := integration.GetValue("services.prefix.app1")
	log.Debug().Msg(value)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to read config file")
	}

	route.ConfigServer(config)
}
