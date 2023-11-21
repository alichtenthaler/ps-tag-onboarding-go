package web

import (
	"encoding/json"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"net/http"

	"github.com/rs/zerolog/log"
)


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

// SendGenericError sends an HTTP error response to the client. The error is a GenericError.
// The provided status code must be a valid HTTP 1xx-5xx status code.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func SendGenericError(w http.ResponseWriter, statusCode int, err errs.GenericError) {
	SendResponse(w, statusCode, err)
}

// SendValidationError sends an HTTP error response to the client. The error is a ValidationError.
// The provided status code must be a valid HTTP 1xx-5xx status code.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func SendValidationError(w http.ResponseWriter, statusCode int, err errs.ValidationError) {
	SendResponse(w, statusCode, err)
}

