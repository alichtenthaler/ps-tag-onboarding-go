package service

import (
	"context"
	domain "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/port/out"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type mockSaveUserFunc func(ctx context.Context, user domain.User) (primitive.ObjectID, error)
type mockExistsByFirstNameAndLastNameFunc func(ctx context.Context, firstName, lastName string) bool

type mockSaveUserPort struct {
	SaveUserFunc                     mockSaveUserFunc
	ExistsByFirstNameAndLastNameFunc mockExistsByFirstNameAndLastNameFunc
}

func newMockSaveUserPort(saveFunc mockSaveUserFunc, userExists mockExistsByFirstNameAndLastNameFunc) *mockSaveUserPort {
	return &mockSaveUserPort{
		SaveUserFunc:                     saveFunc,
		ExistsByFirstNameAndLastNameFunc: userExists,
	}
}

func (mp mockSaveUserPort) SaveUser(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
	return mp.SaveUserFunc(ctx, user)
}

func (mp mockSaveUserPort) ExistsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) bool {
	return mp.ExistsByFirstNameAndLastNameFunc(ctx, firstName, lastName)
}

func TestCreateUserServiceOK(t *testing.T) {
	saveUserPort := newMockSaveUserPort(
		func(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
			return primitive.ObjectID{}, nil
		},
		func(ctx context.Context, firstName, lastName string) bool {
			return false
		},
	)

	user := domain.User{
		FirstName: "John",
		LastName:  "Johnson",
		Email:     "j@j.com",
		Age:       30,
	}

	userService := NewCreateUserService(saveUserPort)
	_, validationErr, err := userService.CreateUser(context.Background(), user)

	assert.Nil(t, validationErr)
	assert.Nil(t, err)
}

func TestCreateUserServiceFailValidation(t *testing.T) {
	mockSaveUserFunc := func(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
		return primitive.ObjectID{}, nil
	}

	mockExistsNameFuncFalse := func(ctx context.Context, firstName, lastName string) bool {
		return false
	}

	mockExistsNameFuncTrue := func(ctx context.Context, firstName, lastName string) bool {
		return true
	}

	mockSaveUserPortWithoutNameConflict := newMockSaveUserPort(
		mockSaveUserFunc,
		mockExistsNameFuncFalse,
	)

	mockSaveUserPortWithNameConflict := newMockSaveUserPort(
		mockSaveUserFunc,
		mockExistsNameFuncTrue,
	)

	testCases := []struct {
		name            string
		user            domain.User
		validationError []string
		saveUserPort    out.SaveUserPort
	}{
		{
			name: "Missing user first name",
			user: domain.User{
				LastName: "ann",
				Email:    "a@a.com",
				Age:      22,
			},
			validationError: []string{domain.ErrorNameRequired},
			saveUserPort:    mockSaveUserPortWithoutNameConflict,
		},
		{
			name: "Missing user last name",
			user: domain.User{
				FirstName: "ann",
				Email:     "a@a.com",
				Age:       22,
			},
			validationError: []string{domain.ErrorNameRequired},
			saveUserPort:    mockSaveUserPortWithoutNameConflict,
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
			saveUserPort:    mockSaveUserPortWithNameConflict,
		},
		{
			name: "Missing user email",
			user: domain.User{
				FirstName: "a",
				LastName:  "ann",
				Age:       22,
			},
			validationError: []string{domain.ErrorEmailRequired},
			saveUserPort:    mockSaveUserPortWithoutNameConflict,
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
			saveUserPort:    mockSaveUserPortWithoutNameConflict,
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
			saveUserPort:    mockSaveUserPortWithoutNameConflict,
		},
		{
			name: "User fails validation on multiple fields",
			user: domain.User{
				LastName: "ann",
				Email:    "aa.com",
				Age:      17,
			},
			validationError: []string{domain.ErrorAgeMinimum, domain.ErrorEmailFormat, domain.ErrorNameRequired},
			saveUserPort:    mockSaveUserPortWithoutNameConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {

			userService := NewCreateUserService(tc.saveUserPort)
			id, validationErr, err := userService.CreateUser(context.Background(), tc.user)

			assert.Equal(tt, primitive.NilObjectID, id)
			assert.Equal(t,  domain.ValidationError{Err:domain.ResponseValidationFailed, Details:tc.validationError}, *validationErr)
			assert.Nil(tt, err)
		})
	}
}
