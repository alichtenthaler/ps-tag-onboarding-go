package out

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SaveUserPort interface {
	SaveUser(ctx context.Context, user domain.User) (primitive.ObjectID, error)
	ExistsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool
}

type GetUserPort interface {
	GetUserById(ctx context.Context, id string) (domain.User, error)
}