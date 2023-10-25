package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

var (
	StringDBConnection = ""
	DBName             = ""
	DBUser             = ""
	DBPass             = ""
	DBHost             = ""
	DBPort             = ""
	Port               = 0
)

// Load initializes environment variables
func Load() {
	log.Info().Msg("Loading configs...")
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}

	DBName = os.Getenv("DB_NAME")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")

	StringDBConnection = fmt.Sprintf("mongodb://%s:%s@%s:%s",
		DBUser,
		DBPass,
		DBHost,
		DBPort,
	)

	log.Info().Msg("Configs loaded.")
}
