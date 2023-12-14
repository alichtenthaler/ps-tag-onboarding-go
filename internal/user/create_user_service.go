package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/rs/zerolog/log"
)

// CreateUserSrv validates and saves the user in the database
type CreateUserSrv struct {
	repo CreateRepository
}

// CreateRepository abstracts the database
type CreateRepository interface {
	create(ctx context.Context, user User) (primitive.ObjectID, error)
	existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool
}

// NewCreateUserService creates a new user service
func NewCreateUserService(db *mongo.Database) CreateUserSrv {
	repo := newRepository(db)
	return CreateUserSrv{
		repo: repo,
	}
}

// CreateUser validate and saves the user in the database
func (s CreateUserSrv) CreateUser(ctx context.Context, user *User) error {

	validationError := user.validate()
	if s.repo.existsByFirstNameAndLastName(ctx, user.FirstName, user.LastName) {
		validationError.Details = append(validationError.Details, ErrorNameUnique.Error())
	}

	if len(validationError.Details) > 0 {
		return validationError
	}

	var err error
	user.ID, err = s.repo.create(ctx, *user)
	if err != nil {
		log.Error().Msg("error saving user in the database")
		return err
	}

	return nil
}

