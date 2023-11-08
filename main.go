package main

import (
	"context"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/config"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/database"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/rest"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	userService := user.New(dbConnection)

	log.Info().Msg("Starting HTTP server")
	server := startRestServer(configs.Port, userService)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Info().Msg("Shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Msg("Failed to shutdown HTTP server")
	}

	log.Info().Msg("Server was shutdown properly")
}

func startRestServer(port int, userProcessor *user.Service) *rest.Rest {
	restServer := rest.New(
		port,
		userProcessor,
	)

	go restServer.Start()
	return restServer
}
