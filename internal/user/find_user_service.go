package user

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindRepository abstracts the database
type FindRepository interface {
	getByID(ctx context.Context, id primitive.ObjectID) (User, error)
}

// FindUserSrv retrieves a user from the database
type FindUserSrv struct {
	repo FindRepository
}

// NewFindUserService creates a new user service
func NewFindUserService(db *mongo.Database) FindUserSrv {
	repo := newRepository(db)
	return FindUserSrv{
		repo: repo,
	}
}

// FindUserById retrieves a user from the database by id
func (s FindUserSrv) FindUserById(ctx context.Context, userId primitive.ObjectID) (*User, error) {

	user, err := s.repo.getByID(ctx, userId)
	if err != nil {
		log.Error().Msgf("error getting user by id in the database: %s", err)
		return nil, err
	}

	if user.ID.IsZero() {
		log.Info().Msgf("no user found with id '%s'", userId)
		return nil, ResponseUserNotFound
	}

	return &user, nil
}
