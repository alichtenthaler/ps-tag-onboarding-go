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

	errs := s.validate(user)
	if s.userPort.ExistsByFirstNameAndLastName(ctx, user.FirstName, user.LastName) {
		errs = append(errs, domain.ErrorNameUnique)
	}

	if len(errs) > 0 {
		return primitive.NilObjectID, ValidationError{Error: domain.ResponseValidationFailed, Details: errs}, nil
	}

	id, err := s.userPort.SaveUser(ctx, user)
	if err != nil {
		return primitive.NilObjectID, ValidationError{}, err
	}

	return id, ValidationError{}, err
}

