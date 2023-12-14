package user

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type MockFindUserRepository struct {
	FindRepository
	GetByIDFunc                      func(ctx context.Context, id primitive.ObjectID) (User, error)
}

func (r MockFindUserRepository) getByID(ctx context.Context, id primitive.ObjectID) (User, error) {
	if r.GetByIDFunc != nil {
		return r.GetByIDFunc(ctx, id)
	}

	return r.FindRepository.getByID(ctx, id)
}

func TestFindUserByIDServiceOK(t *testing.T) {
	userID := primitive.NewObjectID()
	user := User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Johnson",
		Age:       30,
		Email:     "j@j.com",
	}

	userService := &FindUserSrv{
		repo: MockFindUserRepository{
			GetByIDFunc: func(ctx context.Context, id primitive.ObjectID) (User, error) {
				return user, nil
			},
		},
	}

	userReturned, err := userService.FindUserById(context.TODO(), userID)

	assert.NoError(t, err)
	assert.Equal(t, &user, userReturned)
}

func TestFindUserByIDServiceNotFound(t *testing.T) {
	userID := primitive.NewObjectID()

	userService := &FindUserSrv{
		repo: MockFindUserRepository{
			GetByIDFunc: func(ctx context.Context, id primitive.ObjectID) (User, error) {
				return User{}, nil
			},
		},
	}

	user, err := userService.FindUserById(context.TODO(), userID)
	assert.Nil(t, user)
	assert.Equal(t, ResponseUserNotFound, err)
}

