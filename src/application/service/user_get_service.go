package service

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/port/out"
)

// GetUserService is the service that retrieves a user from the database
type GetUserService struct {
	userPort out.GetUserPort
}

// NewGetUserService creates a new instance of the service
func NewGetUserService(userPort out.GetUserPort) *GetUserService {
	return &GetUserService{
		userPort: userPort,
	}
}

// GetUser retrieves a user from the database if it exists or returns an error if it does not
func (s *GetUserService) GetUser(ctx context.Context, id string) (domain.User, error) {
	return s.userPort.GetUserById(ctx, id)
}