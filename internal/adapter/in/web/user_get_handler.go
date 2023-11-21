package web

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/port/in"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

// GetUserHandler is an HTTP handler for getting a user.
type GetUserHandler struct {
	uc in.GetUserUseCase
}

// NewGetUserHandler creates a new GetUserHandler.
func NewGetUserHandler(uc in.GetUserUseCase) *GetUserHandler {
	return &GetUserHandler{
		uc: uc,
	}
}

// HandleGetUser handles the HTTP request for getting a user.
func (h *GetUserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)

	userID := params["userId"]

	user, err := h.uc.GetUser(r.Context(), userID)
	if err != nil {
		log.Error().Msgf("error getting user by id in the database: %s", err.Error())
		SendGenericError(w, http.StatusInternalServerError, errs.GenericError{Err: err.Error()})
		return
	}

	if user.ID.IsZero() {
		log.Info().Msgf("no user found with id '%s'", userID)
		SendGenericError(w, http.StatusNotFound, errs.GenericError{Err: errs.ResponseUserNotFound.Error()})
		return
	}

	SendResponse(w, http.StatusOK, user)
}
