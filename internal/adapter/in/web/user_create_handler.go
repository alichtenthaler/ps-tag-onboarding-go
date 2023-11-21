package web

import (
	"encoding/json"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/port/in"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/rs/zerolog/log"
	"net/http"
)

// CreateUserHandler is an HTTP handler for creating a user.
type CreateUserHandler struct {
	uc in.CreateUserUseCase
}

// NewCreateUserHandler creates a new CreateUserHandler.
func NewCreateUserHandler(uc in.CreateUserUseCase) *CreateUserHandler {
	return &CreateUserHandler{
		uc: uc,
	}
}

// HandleCreteUser handles the HTTP request for creating a user.
func (h *CreateUserHandler) HandleCreteUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Msg(err.Error())
		SendValidationError(w, http.StatusBadRequest, errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: []string{err.Error()}})
		return
	}

	var validationErr errs.ValidationError
	user.ID, validationErr, err = h.uc.CreateUser(r.Context(), user)
	if len(validationErr.Details) > 0 {
		log.Error().Msgf("error validating user: %s", validationErr.Error())
		SendValidationError(w, http.StatusBadRequest, validationErr)
		return
	}
	if err != nil {
		log.Error().Msgf("error saving user in the database: %s", err.Error())
		SendGenericError(w, http.StatusInternalServerError, errs.GenericError{Err: err.Error()})
		return
	}

	SendResponse(w, http.StatusCreated, user)
}
