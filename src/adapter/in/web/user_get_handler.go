package web

import (
	"errors"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/port/in"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GetUserHandler struct {
	uc in.GetUserUseCase
}

func NewGetUserHandler(uc in.GetUserUseCase) *GetUserHandler {
	return &GetUserHandler{
		uc: uc,
	}
}

func (h *GetUserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)

	userID := params["userId"]
	if userID == "" {
		err = errors.New("no user id provided")
		log.Error().Msg(err.Error())
		SendError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.uc.GetUser(r.Context(), userID)
	if err != nil {
		log.Error().Msgf("error getting user by id in the database: %s", err.Error())
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID.IsZero() {
		log.Warn().Msgf("No user found with id '%s'", userID)
		SendError(w, http.StatusNotFound, errors.New(domain.ResponseUserNotFound))
		return
	}

	SendResponse(w, http.StatusOK, user)
}
