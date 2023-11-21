package service

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/application/domain/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type mockGetUserByIdFunc func(ctx context.Context, id string) (domain.User, error)

type mockGetUserPort struct {
	GetUserByIdFunc mockGetUserByIdFunc
}

func newMockGetUserPort(getByIdFunc mockGetUserByIdFunc) *mockGetUserPort {
	return &mockGetUserPort{
		GetUserByIdFunc: getByIdFunc,
	}
}

func (mp mockGetUserPort) GetUserById(ctx context.Context, id string) (domain.User, error) {
	return mp.GetUserByIdFunc(ctx, id)
}

func TestGetUserServiceOK(t *testing.T) {
	userID := primitive.NewObjectID()
	userToBeReturned := domain.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Johnson",
		Age:       30,
		Email:     "j@j.com",
	}

	mockGetUserPort := newMockGetUserPort(
		func(ctx context.Context, id string) (domain.User, error) {
			return userToBeReturned, nil
		},
	)

	userService := NewGetUserService(mockGetUserPort)
	user, err := userService.GetUser(context.Background(), userID.Hex())

	assert.Equal(t, userToBeReturned, user)
	assert.Nil(t, err)
}

func TestGetUserServiceNotFound(t *testing.T) {
	mockGetUserPort := newMockGetUserPort(
		func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{}, nil
		},
	)

	userService := NewGetUserService(mockGetUserPort)
	user, err := userService.GetUser(context.Background(), "a")

	assert.Equal(t, domain.User{}, user)
	assert.True(t, user.ID.IsZero())
	assert.Nil(t, err)
}
