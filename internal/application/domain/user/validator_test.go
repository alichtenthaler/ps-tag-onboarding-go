package domain

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserValidate(t *testing.T) {

	testCases := []struct {
		name            string
		user            User
		validationError []string

	}{
		{
			name: "Missing user first name",
			user: User{
				LastName: "ann",
				Email:    "a@a.com",
				Age:      22,
			},
			validationError: []string{errs.ErrorNameRequired.Error()},
		},
		{
			name: "Missing user last name",
			user: User{
				FirstName: "ann",
				Email:     "a@a.com",
				Age:       22,
			},
			validationError: []string{errs.ErrorNameRequired.Error()},
		},
		{
			name: "Missing user email",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Age:       22,
			},
			validationError: []string{errs.ErrorEmailRequired.Error()},
		},
		{
			name: "User email not in a proper format",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "aa.com",
				Age:       18,
			},
			validationError: []string{errs.ErrorEmailFormat.Error()},
		},
		{
			name: "Minimum age required",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "a@a.com",
				Age:       17,
			},
			validationError: []string{errs.ErrorAgeMinimum.Error()},
		},
		{
			name: "User fails validation on multiple fields",
			user: User{
				LastName: "ann",
				Email:    "aa.com",
				Age:      17,
			},
			validationError: []string{errs.ErrorAgeMinimum.Error(), errs.ErrorEmailFormat.Error(), errs.ErrorNameRequired.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			validationErr := tc.user.Validate()
			assert.Equal(t, errs.ValidationError{Err: errs.ResponseValidationFailed.Error(), Details: tc.validationError}, validationErr)
		})
	}
}
