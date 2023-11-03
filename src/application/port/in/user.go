package in

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserUseCase interface {
	CreateUser(ctx context.Context, user domain.User) (primitive.ObjectID, service.ValidationError, error)
}

type GetUserUseCase interface {
	GetUser(ctx context.Context, id string) (domain.User, error)
}
