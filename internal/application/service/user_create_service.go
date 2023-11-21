package service

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/port/out"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserService is the service that creates a user in the database
type CreateUserService struct {
	userPort out.SaveUserPort
}

// NewCreateUserService creates a new CreateUserService
func NewCreateUserService(userPort out.SaveUserPort) *CreateUserService {
	return &CreateUserService{
		userPort: userPort,
	}
}

// CreateUser creates a user in the database
func (s *CreateUserService) CreateUser(ctx context.Context, user domain.User) (primitive.ObjectID, *domain.ValidationError, error) {

	errs := user.Validate()
	if s.userPort.ExistsByFirstNameAndLastName(ctx, user.FirstName, user.LastName) {
		errs = append(errs, domain.ErrorNameUnique)
	}

	if len(errs) > 0 {
		return primitive.NilObjectID, &domain.ValidationError{Err: domain.ResponseValidationFailed, Details: errs}, nil
	}

	id, err := s.userPort.SaveUser(ctx, user)
	if err != nil {
		return primitive.NilObjectID, nil, err
	}

	return id, nil, err
}

