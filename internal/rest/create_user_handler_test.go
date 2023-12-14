package rest

import (
	"context"
	"encoding/json"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCreateUserService struct {
	CreateUserService
	CreateUserFunc func(ctx context.Context, u *user.User) error
}

func (r MockCreateUserService) CreateUser(ctx context.Context, u *user.User) error {
	if r.CreateUserFunc != nil {
		return r.CreateUserFunc(ctx, u)
	}

	return r.CreateUserService.CreateUser(ctx, u)
}

func TestCreateUserHandlerOK(t *testing.T) {
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":"j@j.com"}`
	var u user.User
	err := json.Unmarshal([]byte(payload), &u)
	if err != nil {
		t.Fatal(err)
	}

	createUserHandler := &CreateUserHandler{
		service: MockCreateUserService{
			CreateUserFunc: func(ctx context.Context, u *user.User) error {
				return nil
			},
		},
	}

	req, err := http.NewRequest("POST", "/save", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	createUserHandler.CreateUser(res, req)

	var responseUser user.User
	err = json.NewDecoder(res.Body).Decode(&responseUser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Equal(t, u.FirstName, responseUser.FirstName)
	assert.Equal(t, u.LastName, responseUser.LastName)
	assert.Equal(t, u.Age, responseUser.Age)
	assert.Equal(t, u.Email, responseUser.Email)
}

func TestCreateUserHandlerFailValidation(t *testing.T) {
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":""}`
	var u user.User
	err := json.Unmarshal([]byte(payload), &u)
	if err != nil {
		t.Fatal(err)
	}

	createUserHandler := &CreateUserHandler{
		service: MockCreateUserService{
			CreateUserFunc: func(ctx context.Context, u *user.User) error {
				return errs.ValidationError{Err: user.ResponseValidationFailed.Error(), Details: []string{user.ErrorEmailRequired.Error()}}
			},
		},
	}

	req, err := http.NewRequest("POST", "/save", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	createUserHandler.CreateUser(res, req)

	var responseError errs.ValidationError
	err = json.NewDecoder(res.Body).Decode(&responseError)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, user.ErrorEmailRequired.Error(), responseError.Details[0])
	assert.Equal(t, user.ResponseValidationFailed.Error(), responseError.Err)
}
