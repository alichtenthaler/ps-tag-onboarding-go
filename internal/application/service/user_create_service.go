package service

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/port/out"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
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
func (s *CreateUserService) CreateUser(ctx context.Context, user domain.User) (primitive.ObjectID, errs.ValidationError, error) {

	var validationErrs errs.ValidationError

	validationErrs = user.Validate()
	if s.userPort.ExistsByFirstNameAndLastName(ctx, user.FirstName, user.LastName) {
		validationErrs.Details = append(validationErrs.Details, errs.ErrorNameUnique.Message)
	}

	if len(validationErrs.Details) > 0 {
		return primitive.NilObjectID, errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: validationErrs.Details}, nil
	}

	id, err := s.userPort.SaveUser(ctx, user)
	if err != nil {
		return primitive.NilObjectID, validationErrs, err
	}

	return id, validationErrs, err
}

