package web

import (
	"encoding/json"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/domain/user"
	"net/http"

	"github.com/rs/zerolog/log"
)

type GenericError struct {
	Error string `json:"error"`
}

// SendResponse sends an HTTP response to the client. The data is encoded as JSON.
// The provided status code must be a valid HTTP 1xx-5xx status code.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func SendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
}

// SendError sends an HTTP error response to the client. The error is a GenericError.
// The provided status code must be a valid HTTP 1xx-5xx status code.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func SendError(w http.ResponseWriter, statusCode int, err error) {
	SendResponse(w, statusCode, GenericError{
		Error: err.Error(),
	})
}

// SendValidationError sends an HTTP error response to the client. The error is a ValidationError.
// The provided status code must be a valid HTTP 1xx-5xx status code.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func SendValidationError(w http.ResponseWriter, statusCode int, err domain.ValidationError) {
	SendResponse(w, statusCode, err)
}

