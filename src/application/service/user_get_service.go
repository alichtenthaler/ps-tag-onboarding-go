package service

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/port/out"
)

type GetUserService struct {
	userPort out.GetUserPort
}

func NewGetUserService(userPort out.GetUserPort) *GetUserService {
	return &GetUserService{
		userPort: userPort,
	}
}

func (s *GetUserService) GetUser(ctx context.Context, id string) (domain.User, error) {
	return s.userPort.GetUserById(ctx, id)
}