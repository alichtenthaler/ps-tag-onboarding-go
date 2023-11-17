package in

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserUseCase abstracts the service that creates a user in the database
type CreateUserUseCase interface {
	CreateUser(ctx context.Context, user domain.User) (primitive.ObjectID, *domain.ValidationError, error)
}

// GetUserUseCase abstracts the service that retrieves a user from the database
type GetUserUseCase interface {
	GetUser(ctx context.Context, id string) (domain.User, error)
}
