package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

// Configuration represents the application configuration
type Configuration struct {
	Port        int
	Environment string
	DBConfig    DBConfig
	LogFormat   string
	LogLevel    zerolog.Level
}

// DBConfig represents the database configuration
type DBConfig struct {
	Name string
	User string
	Pass string
	Host string
}

// Load loads the application configuration from the environment variables and returns it
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
			Name: os.Getenv("DB_NAME"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Host: os.Getenv("DB_HOST"),
		},
		LogFormat: os.Getenv("LOG_FORMAT"),
		LogLevel:  logLevel,
	}
}
