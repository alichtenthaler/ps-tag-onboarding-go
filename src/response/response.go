package response

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ValidationError struct {
	Error   string   `json:"error"`
	Details []string `json:"details"`
}

type GenericError struct {
	Error string `json:"error"`
}

func SendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
}

func SendError(w http.ResponseWriter, statusCode int, err error) {
	SendResponse(w, statusCode, GenericError{
		Error: err.Error(),
	})
}

func SendValidationError(w http.ResponseWriter, statusCode int, err ValidationError) {
	SendResponse(w, statusCode, ValidationError{
		Error:   err.Error,
		Details: err.Details,
	})
}
