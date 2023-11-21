package main

import (
	"context"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/adapter/in/web"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/adapter/out/mongo"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/service"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/config"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/database"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/rest"
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

	log.Info().Msg("Waiting for connection to database...")

	dbConnection, err := database.Connect(configs.DBConfig)
	if err != nil {
		log.Fatal().Msgf("Connection to the database failed: %s", err)
	}

	log.Info().Msg("Connected to the database")

	userPersistenceAdapter := mongo.NewUserPersistenceAdapter(dbConnection)
	createUserService := service.NewCreateUserService(userPersistenceAdapter)
	getUserService := service.NewGetUserService(userPersistenceAdapter)

	createUserHandler := web.NewCreateUserHandler(createUserService)
	getUserHandler := web.NewGetUserHandler(getUserService)

	log.Info().Msg("Starting HTTP server")
	server := startRestServer(configs.Port, getUserHandler, createUserHandler)

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

func startRestServer(port int, getUserHandler *web.GetUserHandler, createUserHandler *web.CreateUserHandler) *rest.Rest {
	restServer := rest.New(
		port,
		getUserHandler,
		createUserHandler,
	)

	go restServer.Start()
	return restServer
}
