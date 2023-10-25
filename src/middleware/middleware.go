package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Logger logs every request received by the API
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msgf("%s %s %s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}