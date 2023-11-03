package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/adapter/in/web"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockUserValidator struct {
	domain.ValidatorI
	ValidateFunc func(ctx context.Context, user domain.User) []string
}

func (v MockUserValidator) validate(ctx context.Context, user domain.User) []string {
	if v.ValidateFunc != nil {
		return v.ValidateFunc(ctx, user)
	}

	return v.ValidatorI.validate(ctx, user)
}

type MockUserRepository struct {
	UserRepositoryI
	CreateFunc  func(ctx context.Context, user domain.User) (primitive.ObjectID, error)
	GetByIDFunc func(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	ExistsByFirstNameAndLastNameFunc func(ctx context.Context, firstName, lastName string) bool
}

func (r MockUserRepository) create(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, user)
	}

	return r.UserRepositoryI.create(ctx, user)
}

func (r MockUserRepository) getByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	if r.GetByIDFunc != nil {
		return r.GetByIDFunc(ctx, id)
	}

	return r.UserRepositoryI.getByID(ctx, id)
}

func (r MockUserRepository) existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {
	if r.ExistsByFirstNameAndLastNameFunc != nil {
		return r.ExistsByFirstNameAndLastNameFunc(ctx, firstName, lastName)
	}

	return r.UserRepositoryI.existsByFirstNameAndLastName(ctx, firstName, lastName)
}

func TestCreateUserHandlerOK(t *testing.T) {
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":"j@j.com"}`
	var user domain.User
	err := json.Unmarshal([]byte(payload), &user)
	if err != nil {
		t.Fatal(err)
	}

	userService := &Service{
		repo: MockUserRepository{
			CreateFunc: func(ctx context.Context, user user.User) (primitive.ObjectID, error) {
				return primitive.NewObjectID(), nil
			},
			ExistsByFirstNameAndLastNameFunc: func(ctx context.Context, firstName, lastName string) bool {
				return false
			},
		},
		validator: MockUserValidator{
			ValidateFunc: func(ctx context.Context, user user.User) []string {
				return []string{}
			},
		},
	}

	req, err := http.NewRequest("POST", "/save", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	userService.CreateUser(res, req)

	var responseUser user.User
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
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":"j@j.com"}`
	var user domain.User
	err := json.Unmarshal([]byte(payload), &user)
	if err != nil {
		t.Fatal(err)
	}

	userService := &Service{
		repo: MockUserRepository{
			CreateFunc: func(ctx context.Context, user user.User) (primitive.ObjectID, error) {
				return primitive.NewObjectID(), nil
			},
		},
		validator: MockUserValidator{
			ValidateFunc: func(ctx context.Context, user user.User) []string {
				return []string{user.ErrorNameUnique, user.ErrorEmailRequired}
			},
		},
	}

	req, err := http.NewRequest("POST", "/save", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	userService.CreateUser(res, req)

	var responseError web.ValidationError
	err = json.NewDecoder(res.Body).Decode(&responseError)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, user.ErrorNameUnique, responseError.Details[0])
	assert.Equal(t, user.ErrorEmailRequired, responseError.Details[1])
	assert.Equal(t, user.ResponseValidationFailed, responseError.Error)
}

func TestUserValidate(t *testing.T) {

	mockRepoNameConflictFalse := MockUserRepository{
		ExistsByFirstNameAndLastNameFunc: func(ctx context.Context, firstName, lastName string) bool {
			return false
		},
	}

	mockRepoNameConflictTrue := MockUserRepository{
		ExistsByFirstNameAndLastNameFunc: func(ctx context.Context, firstName, lastName string) bool {
			return true
		},
	}

	testCases := []struct {
		name            string
		user            domain.User
		validationError []string
		validatorRepo   MockUserRepository
	}{
		{
			name: "Missing user first name",
			user: domain.User{
				LastName: "ann",
				Email:    "a@a.com",
				Age:      22,
			},
			validationError: []string{domain.ErrorNameRequired},
		},
		{
			name: "Missing user last name",
			user: domain.User{
				FirstName: "ann",
				Email:     "a@a.com",
				Age:       22,
			},
			validationError: []string{domain.ErrorNameRequired},
		},
		{
			name: "User with the same first and last name already exists",
			user: domain.User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "a@a.com",
				Age:       22,
			},
			validationError: []string{domain.ErrorNameUnique},
			validatorRepo:   mockRepoNameConflictTrue,
		},
		{
			name: "Missing user email",
			user: domain.User{
				FirstName: "a",
				LastName:  "ann",
				Age:       22,
			},
			validationError: []string{domain.ErrorEmailRequired},
			validatorRepo:   mockRepoNameConflictFalse,
		},
		{
			name: "User email not in a proper format",
			user: domain.User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "aa.com",
				Age:       18,
			},
			validationError: []string{domain.ErrorEmailFormat},
			validatorRepo:   mockRepoNameConflictFalse,
		},
		{
			name: "Minimum age required",
			user: domain.User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "a@a.com",
				Age:       17,
			},
			validationError: []string{domain.ErrorAgeMinimum},
			validatorRepo:   mockRepoNameConflictFalse,
		},
		{
			name: "User fails validation on multiple fields",
			user: domain.User{
				LastName: "ann",
				Email:    "aa.com",
				Age:      17,
			},
			validationError: []string{domain.ErrorAgeMinimum, domain.ErrorEmailFormat, domain.ErrorNameRequired},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			validator := domain.Validator{
				repo: tc.validatorRepo,
			}

			errs := validator.validate(context.Background(), tc.user)
			assert.Equal(t, tc.validationError, errs)
		})
	}
}

func TestFindUserByIDHandlerOK(t *testing.T) {
	userID := primitive.NewObjectID()
	userToBeReturned := domain.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Johnson",
		Age:       30,
		Email:     "j@j.com",
	}

	userService := &Service{
		repo: MockUserRepository{
			GetByIDFunc: func(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
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

	var respUser domain.User
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

	var respError web.GenericError
	err = json.NewDecoder(res.Body).Decode(&respError)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, domain.ResponseUserNotFound, respError.Error)
}
