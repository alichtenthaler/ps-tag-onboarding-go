package rest

import (
	"context"
	"net/http"

	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindUserService provides application operations to retrieve a user from the database
type FindUserService interface {
	FindUserById(ctx context.Context, userId primitive.ObjectID) (*user.User, error)
}

// FindUserHandler defines routes and HTTP handlers for user operations
type FindUserHandler struct {
	service FindUserService
}

// NewFindUserHandler creates a new FindUserHandler
func NewFindUserHandler(service FindUserService) FindUserHandler {
	return FindUserHandler{service: service}
}

// FindUser handles the user find (/user/find/{userId}) endpoint
func (h FindUserHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)

	userID := params["userId"]

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Error().Msgf("error converting user id to object id")
		SendError(w, user.ResponseUserNotFound)
		return
	}

	u, err := h.service.FindUserById(r.Context(), objectID)
	if err != nil {
		SendError(w, err)
		return
	}

	SendResponse(w, http.StatusOK, u)
}
