package rest

import (
	"context"
	"encoding/json"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/user"
	"github.com/rs/zerolog/log"
	"net/http"
)

// CreateUserService provides application operations
type CreateUserService interface {
	CreateUser(ctx context.Context, user *user.User) error
}

// CreateUserHandler defines routes and HTTP handlers for user operations
type CreateUserHandler struct {
	service CreateUserService
}

// NewCreateUserHandler creates a new CreateUserHandler
func NewCreateUserHandler(service CreateUserService) CreateUserHandler {
	return CreateUserHandler{service: service}
}

// CreateUser handles the user creation (/user/save) endpoint
func (h CreateUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error().Msg(err.Error())
		SendError(w, errs.ValidationError{Err: user.ResponseValidationFailed.Message, Details: []string{err.Error()}})
		return
	}

	err = h.service.CreateUser(r.Context(), &u)
	if err != nil {
		SendError(w, err)
		return
	}

	SendResponse(w, http.StatusCreated, u)
}
