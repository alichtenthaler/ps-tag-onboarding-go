package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/user"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockFindUserService struct {
	FindUserService
	FindUserByIdFunc func(ctx context.Context, userId primitive.ObjectID) (*user.User, error)
}

func (r MockFindUserService) FindUserById(ctx context.Context, id primitive.ObjectID) (*user.User, error) {
	if r.FindUserByIdFunc != nil {
		return r.FindUserByIdFunc(ctx, id)
	}

	return r.FindUserService.FindUserById(ctx, id)
}

func TestFindUserByIDHandlerOK(t *testing.T) {
	userID := primitive.NewObjectID()
	userToBeReturned := user.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Johnson",
		Age:       30,
		Email:     "j@j.com",
	}

	findUserHandler := &FindUserHandler{
		service: MockFindUserService{
			FindUserByIdFunc: func(ctx context.Context, id primitive.ObjectID) (*user.User, error) {
				return &userToBeReturned, nil
			},
		},
	}

	userIdURLParam := userID.Hex()
	req, err := http.NewRequest("GET", fmt.Sprintf("/find/%s", userIdURLParam), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"userId": userIdURLParam})

	res := httptest.NewRecorder()
	findUserHandler.FindUser(res, req)

	var respUser user.User
	err = json.NewDecoder(res.Body).Decode(&respUser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, userToBeReturned, respUser)
}

func TestFindUserByIDHandlerNotFound(t *testing.T) {
	userID := primitive.NewObjectID()

	findUserHandler := &FindUserHandler{
		service: MockFindUserService{
			FindUserByIdFunc: func(ctx context.Context, id primitive.ObjectID) (*user.User, error) {
				return nil, user.ResponseUserNotFound
			},
		},
	}

	userIdURLParam := userID.Hex()
	req, err := http.NewRequest("GET", fmt.Sprintf("/find/%s", userIdURLParam), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"userId": userIdURLParam})

	res := httptest.NewRecorder()
	findUserHandler.FindUser(res, req)

	var respError GenericError
	err = json.NewDecoder(res.Body).Decode(&respError)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, user.ResponseUserNotFound.Error(), respError.Error)
}
