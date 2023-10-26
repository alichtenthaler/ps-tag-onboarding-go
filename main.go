package main

import (
	"context"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/config"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/database"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/rest"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/user"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dbConnection *mongo.Database

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Timestamp().Logger()

	config.Load()

	log.Info().Msg("Waiting for connection to database...")
	dbConnection = createDBConnection()
	log.Info().Msg("Connected to the database")

	userService := user.New(dbConnection)
	server := startRestServer(userService)

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

func createDBConnection() *mongo.Database {
	db, err := database.Connect(config.StringDBConnection)
	if err != nil {
		log.Fatal().Msgf("Connection to the database failed: %s", err)
	}

	return db
}

func startRestServer(userProcessor *user.Service) *rest.Rest {
	restServer := rest.New(
		userProcessor,
	)

	go restServer.Start()
	return restServer
}
