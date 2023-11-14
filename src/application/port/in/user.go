package in

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserUseCase interface {
	CreateUser(ctx context.Context, user domain.User) (primitive.ObjectID, *domain.ValidationError, error)
}

type GetUserUseCase interface {
	GetUser(ctx context.Context, id string) (domain.User, error)
}
