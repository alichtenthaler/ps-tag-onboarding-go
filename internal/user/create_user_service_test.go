package user

import (
	"context"
	"encoding/json"

	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"

	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockCreateUserRepository struct {
	CreateRepository
	CreateFunc                       func(ctx context.Context, user User) (primitive.ObjectID, error)
	ExistsByFirstNameAndLastNameFunc func(ctx context.Context, firstName, lastName string) bool
}

func (r MockCreateUserRepository) create(ctx context.Context, user User) (primitive.ObjectID, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, user)
	}

	return r.CreateRepository.create(ctx, user)
}

func (r MockCreateUserRepository) existsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {
	if r.ExistsByFirstNameAndLastNameFunc != nil {
		return r.ExistsByFirstNameAndLastNameFunc(ctx, firstName, lastName)
	}

	return r.CreateRepository.existsByFirstNameAndLastName(ctx, firstName, lastName)
}

func TestCreateUserServiceOK(t *testing.T) {
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":"j@j.com"}`
	var user User
	err := json.Unmarshal([]byte(payload), &user)
	if err != nil {
		t.Fatal(err)
	}

	userService := &CreateUserSrv{
		repo: MockCreateUserRepository{
			CreateFunc: func(ctx context.Context, user User) (primitive.ObjectID, error) {
				return primitive.NewObjectID(), nil
			},
			ExistsByFirstNameAndLastNameFunc: func(ctx context.Context, firstName, lastName string) bool {
				return false
			},
		},
	}

	err = userService.CreateUser(context.TODO(), &user)
	assert.NoError(t, err)
	assert.NotNil(t, user.ID)
}

func TestCreateUserServiceFailValidation(t *testing.T) {
	payload := `{"firstName":"John","lastName":"Johnson","age":30,"email":""}`
	var user User
	err := json.Unmarshal([]byte(payload), &user)
	if err != nil {
		t.Fatal(err)
	}

	userService := &CreateUserSrv{
		repo: MockCreateUserRepository{
			CreateFunc: func(ctx context.Context, user User) (primitive.ObjectID, error) {
				return primitive.NewObjectID(), nil
			},
			ExistsByFirstNameAndLastNameFunc: func(ctx context.Context, firstName, lastName string) bool {
				return true
			},
		},
	}

	err = userService.CreateUser(context.TODO(), &user)
	assert.NotNil(t, err)

	assert.Equal(t, ErrorEmailRequired.Error(), err.(errs.ValidationError).Details[0])
	assert.Equal(t, ErrorNameUnique.Error(), err.(errs.ValidationError).Details[1])
	assert.Equal(t, ResponseValidationFailed.Error(), err.(errs.ValidationError).Err)
}

func TestUserValidate(t *testing.T) {

	testCases := []struct {
		name            string
		user            User
		validationError []string
		validatorRepo MockCreateUserRepository
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
