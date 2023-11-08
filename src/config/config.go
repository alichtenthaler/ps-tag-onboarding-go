package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type Configuration struct {
	Port        int
	Environment string
	DBConfig    DBConfig
	LogFormat   string
	LogLevel    zerolog.Level
}

type DBConfig struct {
	DBName          string
	DBConnectionURI string
}

func Load() *Configuration {
	log.Info().Msg("Loading configs...")
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal().Msgf("Error loading configs from .env file: %s", err.Error())
	}

	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		port = 9000
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	if dbName == "" || dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" {
		log.Fatal().Msg("Error loading configs from .env file: one or more DB configs are missing")
	}

	stringDBConnection := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
	)

	logLevel, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Error().Msgf("Error parsing log level: %s. Setting it to InfoLevel", err.Error())
		logLevel = zerolog.InfoLevel
	}

	log.Info().Msg("Configs loaded.")

	return &Configuration{
		Port:        port,
		Environment: os.Getenv("ENVIRONMENT"),
		DBConfig: DBConfig{
			DBName:          dbName,
			DBConnectionURI: stringDBConnection,
		},
		LogFormat: os.Getenv("LOG_FORMAT"),
		LogLevel:  logLevel,
	}
}
