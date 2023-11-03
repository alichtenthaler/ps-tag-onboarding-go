package user

import (
	"encoding/json"
	"errors"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/adapter/in/web"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Service struct {
	repo      UserRepositoryI
	validator domain.ValidatorI
}

func New(db *mongo.Database) *Service {
	repo := newRepository(db)
	return &Service{
		repo:      repo,
		validator: domain.newValidator(repo),
	}
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Msg(err.Error())
		web.SendValidationError(w, http.StatusBadRequest, web.ValidationError{Error: user.ResponseValidationFailed, Details: []string{err.Error()}})
		return
	}

	if errs := s.validator.validate(r.Context(), user); len(errs) > 0 {
		log.Error().Msgf("error validating user: %s", strings.Join(errs, ", "))
		web.SendValidationError(w, http.StatusBadRequest, web.ValidationError{Error: user.ResponseValidationFailed, Details: errs})
		return
	}

	user.ID, err = s.repo.create(r.Context(), user)
	if err != nil {
		log.Error().Msgf("error saving user in the database: %s", err.Error())
		web.SendError(w, http.StatusInternalServerError, err)
		return
	}

	web.SendResponse(w, http.StatusCreated, user)
}

func (s *Service) FindUserById(w http.ResponseWriter, r *http.Request) {

	var err error
	params := mux.Vars(r)

	userID := params["userId"]
	if userID == "" {
		err = errors.New("no user id provided")
		log.Error().Msg(err.Error())
		web.SendError(w, http.StatusBadRequest, err)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Error().Msgf("error converting user id to object id: %s", err.Error())
		web.SendError(w, http.StatusNotFound, errors.New(domain.ResponseUserNotFound))
		return
	}

	user, err := s.repo.getByID(r.Context(), objectID)
	if err != nil {
		log.Error().Msgf("error getting user by id in the database: %s", err.Error())
		web.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID.IsZero() {
		log.Warn().Msgf("No user found with id '%s'", userID)
		web.SendError(w, http.StatusNotFound, errors.New(user.ResponseUserNotFound))
		return
	}

	web.SendResponse(w, http.StatusOK, user)
}
