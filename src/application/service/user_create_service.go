package service

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/port/out"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserService struct {
	userPort out.SaveUserPort
}

func NewCreateUserService(userPort out.SaveUserPort) *CreateUserService {
	return &CreateUserService{
		userPort: userPort,
	}
}

func (s *CreateUserService) CreateUser(ctx context.Context, user domain.User) (primitive.ObjectID, ValidationError, error) {

	validationErr := s.validate(ctx, user)
	if len(validationErr.Details) > 0 {
		return primitive.NilObjectID, validationErr, nil
	}

	id, err := s.userPort.SaveUser(ctx, user)

	return id, ValidationError{}, err
}

