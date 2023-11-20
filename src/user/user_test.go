package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/errs"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/response"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockUserRepository struct {
	Repository
	CreateFunc                       func(ctx context.Context, user User) (primitive.ObjectID, error)
	GetByIDFunc                      func(ctx context.Context, id primitive.ObjectID) (User, error)
	ExistsByFirstNameAndLastNameFunc func(ctx context.Context, firstName, lastName string) bool
}

func (r MockUserRepository) create(ctx context.Context, user User) (primitive.ObjectID, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, user)
	}

	return r.Repository.create(ctx, user)
}

func (r MockUserRepository) getByID(ctx context.Context, id primitive.ObjectID) (User, error) {
	if r.GetByIDFunc != nil {
		return r.GetByIDFunc(ctx, id)
	}

	return r.Repository.getByID(ctx, id)
}

func (r MockUserRepository) existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {
	if r.ExistsByFirstNameAndLastNameFunc != nil {
		return r.ExistsByFirstNameAndLastNameFunc(ctx, firstName, lastName)
	}

	return r.Repository.existsByFirstNameAndLastName(ctx, firstName, lastName)
}

func TestCreateUserHandlerOK(t *testing.T) {
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":"j@j.com"}`
	var user User
	err := json.Unmarshal([]byte(payload), &user)
	if err != nil {
		t.Fatal(err)
	}

	userService := &Service{
		repo: MockUserRepository{
			CreateFunc: func(ctx context.Context, user User) (primitive.ObjectID, error) {
				return primitive.NewObjectID(), nil
			},
			ExistsByFirstNameAndLastNameFunc: func(ctx context.Context, firstName, lastName string) bool {
				return false
			},
		},
	}

	req, err := http.NewRequest("POST", "/save", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	userService.CreateUser(res, req)

	var responseUser User
	err = json.NewDecoder(res.Body).Decode(&responseUser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Equal(t, user.FirstName, responseUser.FirstName)
	assert.Equal(t, user.LastName, responseUser.LastName)
	assert.Equal(t, user.Age, responseUser.Age)
	assert.Equal(t, user.Email, responseUser.Email)
}

func TestCreateUserHandlerFailValidation(t *testing.T) {
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":""}`
	var user User
	err := json.Unmarshal([]byte(payload), &user)
	if err != nil {
		t.Fatal(err)
	}

	userService := &Service{
		repo: MockUserRepository{
			CreateFunc: func(ctx context.Context, user User) (primitive.ObjectID, error) {
				return primitive.NewObjectID(), nil
			},
			ExistsByFirstNameAndLastNameFunc: func(ctx context.Context, firstName, lastName string) bool {
				return true
			},
		},
	}

	req, err := http.NewRequest("POST", "/save", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	userService.CreateUser(res, req)

	var responseError errs.ValidationError
	err = json.NewDecoder(res.Body).Decode(&responseError)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, ErrorEmailRequired.Error(), responseError.Details[0])
	assert.Equal(t, ErrorNameUnique.Error(), responseError.Details[1])
	assert.Equal(t, ResponseValidationFailed.Error(), responseError.Err)
}

func TestUserValidate(t *testing.T) {

	testCases := []struct {
		name            string
		user            User
		validationError []string
		validatorRepo   MockUserRepository
	}{
		{
			name: "Missing user first name",
			user: User{
				LastName: "ann",
				Email:    "a@a.com",
				Age:      22,
			},
			validationError: []string{ErrorNameRequired.Error()},
		},
		{
			name: "Missing user last name",
			user: User{
				FirstName: "ann",
				Email:     "a@a.com",
				Age:       22,
			},
			validationError: []string{ErrorNameRequired.Error()},
		},
		{
			name: "Missing user email",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Age:       22,
			},
			validationError: []string{ErrorEmailRequired.Error()},
		},
		{
			name: "User email not in a proper format",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "aa.com",
				Age:       18,
			},
			validationError: []string{ErrorEmailFormat.Error()},
		},
		{
			name: "Minimum age required",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "a@a.com",
				Age:       17,
			},
			validationError: []string{ErrorAgeMinimum.Error()},
		},
		{
			name: "User fails validation on multiple fields",
			user: User{
				LastName: "ann",
				Email:    "aa.com",
				Age:      17,
			},
			validationError: []string{ErrorAgeMinimum.Error(), ErrorEmailFormat.Error(), ErrorNameRequired.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			validationErr := tc.user.validate()
			assert.Equal(tt, errs.ValidationError{Err: ResponseValidationFailed.Error(), Details: tc.validationError}, validationErr)
		})
	}
}

func TestFindUserByIDHandlerOK(t *testing.T) {
	userID := primitive.NewObjectID()
	userToBeReturned := User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Johnson",
		Age:       30,
		Email:     "j@j.com",
	}

	userService := &Service{
		repo: MockUserRepository{
			GetByIDFunc: func(ctx context.Context, id primitive.ObjectID) (User, error) {
				return userToBeReturned, nil
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
	userService.FindUserById(res, req)

	var respUser User
	err = json.NewDecoder(res.Body).Decode(&respUser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, userToBeReturned, respUser)
}

func TestFindUserByIDHandlerNotFound(t *testing.T) {
	userID := "a"

	userService := &Service{}

	req, err := http.NewRequest("GET", fmt.Sprintf("/find/%s", userID), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"userId": userID})

	res := httptest.NewRecorder()
	userService.FindUserById(res, req)

	var respError response.GenericError
	err = json.NewDecoder(res.Body).Decode(&respError)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, ResponseUserNotFound.Error(), respError.Error)
}
