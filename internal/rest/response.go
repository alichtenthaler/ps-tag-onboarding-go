package rest

import (
	"encoding/json"
	"errors"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/user"
	"net/http"

	"github.com/rs/zerolog/log"
)

// GenericError represents an error to be return to the client
type GenericError struct {
	Error string `json:"error"`
}

// SendResponse sends an HTTP response to the client. The data is encoded as JSON.
// The provided status code must be a valid HTTP 1xx-5xx status code.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func SendResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Error().Msgf("error encoding response data to the client: %s", err.Error())
			http.Error(w, "Internal Server Err", http.StatusInternalServerError)
		}
	}
}

// SendError sends an HTTP error response to the client.
// The provided status code must be a valid HTTP 1xx-5xx status code.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func SendError(w http.ResponseWriter, err error) {
	handleError(w, err)
}

func handleError(w http.ResponseWriter, err error) {

	log.Error().Msgf("error processing the request: %s", err)

	var response interface{}
	var httpCode int

	var validationError errs.ValidationError

	switch {
	case errors.As(err, &validationError):
		response = errs.ValidationError{
			Err:     err.(errs.ValidationError).Err,
			Details: err.(errs.ValidationError).Details,
		}
		httpCode = http.StatusBadRequest
	case errors.Is(err, user.ResponseUserNotFound):
		response = GenericError{
			Error: err.Error(),
		}
		httpCode = http.StatusNotFound
	default:
		response = GenericError{
			Error: err.Error(),
		}
		httpCode = http.StatusInternalServerError
	}

	SendResponse(w, httpCode, response)
}
