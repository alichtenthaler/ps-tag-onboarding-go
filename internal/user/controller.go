package user

import (
	"context"
	"encoding/json"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Service handles the user endpoints
type Service struct {
	repo Repository
}

// Repository abstracts the database
type Repository interface {
	create(ctx context.Context, user User) (primitive.ObjectID, error)
	getByID(ctx context.Context, id primitive.ObjectID) (User, error)
	existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool
}

// New creates a new user service
func New(db *mongo.Database) *Service {
	repo := newRepository(db)
	return &Service{
		repo: repo,
	}
}

// CreateUser handles the user creation (/user/save) endpoint
func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Msg(err.Error())
		response.SendValidationError(w, http.StatusBadRequest, errs.ValidationError{Err: ResponseValidationFailed.Message, Details: []string{err.Error()}})
		return
	}

	validationError := user.validate()
	if s.repo.existsByFirstNameAndLastName(r.Context(), user.FirstName, user.LastName) {
		validationError.Details = append(validationError.Details, ErrorNameUnique.Error())
	}

	if len(validationError.Details) > 0 {
		log.Error().Msg(validationError.Error())
		response.SendValidationError(w, http.StatusBadRequest, validationError)
		return
	}

	user.ID, err = s.repo.create(r.Context(), user)
	if err != nil {
		log.Error().Msgf("error saving user in the database: %s", err.Error())
		response.SendError(w, http.StatusInternalServerError, errs.Error{Message: err.Error()})
		return
	}

	response.SendResponse(w, http.StatusCreated, user)
}

// FindUserById handles the user find by id (/user/find/{userId}) endpoint
func (s *Service) FindUserById(w http.ResponseWriter, r *http.Request) {

	var err error
	params := mux.Vars(r)

	userID := params["userId"]

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Error().Msgf("error converting user id to object id: %s", err.Error())
		response.SendError(w, http.StatusNotFound, ResponseUserNotFound)
		return
	}

	user, err := s.repo.getByID(r.Context(), objectID)
	if err != nil {
		log.Error().Msgf("error getting user by id in the database: %s", err.Error())
		response.SendError(w, http.StatusInternalServerError, errs.Error{Message: err.Error()})
		return
	}

	if user.ID.IsZero() {
		log.Info().Msgf("no user found with id '%s'", userID)
		response.SendError(w, http.StatusNotFound, ResponseUserNotFound)
		return
	}

	response.SendResponse(w, http.StatusOK, user)
}
