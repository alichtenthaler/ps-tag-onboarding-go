package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Processor struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Processor {
	return &Processor{db}
}

// CreateUser creates a user and saves it in the database
func (up *Processor) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Msg(err.Error())
		response.SendValidationError(w, http.StatusBadRequest, response.ValidationError{Error: ResponseValidationFailed, Details: []string{err.Error()}})
		return
	}

	if errs := up.validate(user); len(errs) > 0 {
		log.Error().Msgf("error validating user: %s", strings.Join(errs, ", "))
		response.SendValidationError(w, http.StatusBadRequest, response.ValidationError{Error: ResponseValidationFailed, Details: errs})
		return
	}

	user.ID, err = up.create(user)
	if err != nil {
		log.Error().Msgf("error saving user in the database: %s", err.Error())
		response.SendError(w, http.StatusInternalServerError, err)
		return
	}

	response.SendResponse(w, http.StatusCreated, user)
}

// FindUserById returns a user by id
func (up *Processor) FindUserById(w http.ResponseWriter, r *http.Request) {

	var err error
	params := mux.Vars(r)

	userID := params["userId"]
	if userID == "" {
		err = errors.New("no user id provided")
		log.Error().Msg(err.Error())
		response.SendError(w, http.StatusBadRequest, err)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Error().Msgf("error converting user id to object id: %s", err.Error())
		response.SendError(w, http.StatusNotFound, errors.New(ResponseUserNotFound))
		return
	}

	user, err := up.getByID(context.Background(), objectID)
	if err != nil {
		log.Error().Msgf("error getting user by id in the database: %s", err.Error())
		response.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID.IsZero() {
		log.Warn().Msgf("No user found with id '%s'", userID)
		response.SendError(w, http.StatusNotFound, errors.New(ResponseUserNotFound))
		return
	}

	response.SendResponse(w, http.StatusOK, user)
}
