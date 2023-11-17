package web

import (
	"encoding/json"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/port/in"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
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
		SendValidationError(w, http.StatusBadRequest, domain.ValidationError{Err: domain.ResponseValidationFailed, Details: []string{err.Error()}})
		return
	}

	var validationErr *domain.ValidationError
	user.ID, validationErr, err = h.uc.CreateUser(r.Context(), user)
	if validationErr != nil {
		log.Error().Msgf("error validating user: %s", strings.Join(validationErr.Details, ", "))
		SendValidationError(w, http.StatusBadRequest, domain.ValidationError{Err: validationErr.Err, Details: validationErr.Details})
		return
	}
	if err != nil {
		log.Error().Msgf("error saving user in the database: %s", err.Error())
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	SendResponse(w, http.StatusCreated, user)
}
