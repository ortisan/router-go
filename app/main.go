package main

import (
	"github.com/ortisan/router-go/route"
	"github.com/rs/zerolog"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	route.CreateRoutes()
}
