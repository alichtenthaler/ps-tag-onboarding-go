package main

import (
	"io"
	"os"

	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/config"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/database"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/rest"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/user"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	configs := config.Load()

	var logOutput io.Writer
	logOutput = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	}
	if configs.LogFormat == "json" {
		logOutput = os.Stderr
	}
	log.Logger = zerolog.New(logOutput).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(configs.LogLevel)

	log.Info().Msg("Waiting for connection to database....")

	dbConnection, err := database.Connect(configs.DBConfig)
	if err != nil {
		log.Fatal().Msgf("Connection to the database failed: %s", err)
	}

	log.Info().Msg("Connected to the database")

	createUserService := user.NewCreateUserService(dbConnection)
	createUserHandler := rest.NewCreateUserHandler(createUserService)

	findUserService := user.NewFindUserService(dbConnection)
	findUserHandler := rest.NewFindUserHandler(findUserService)

	log.Info().Msg("Starting HTTP server")

	rest.New(rest.NewRouter(createUserHandler, findUserHandler), rest.WithPort(configs.Port)).Start()
}
