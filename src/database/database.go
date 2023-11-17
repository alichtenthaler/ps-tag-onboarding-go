package database

import (
	"context"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// Connect opens a connection with mongo and returns it
func Connect(config config.DBConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := createClientOptions(config)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(config.Name), nil
}

func createClientOptions(config config.DBConfig) *options.ClientOptions {
	clientOptions := options.Client()

	if config.Host == "" {
		log.Fatal().Msg("No host provided for database connection")
	}

	clientOptions.SetHosts([]string{config.Host})

	if config.User != "" && config.Pass != "" {
		clientOptions.SetAuth(options.Credential{
			Username: config.User,
			Password: config.Pass,
		})
	}

	return clientOptions
}
