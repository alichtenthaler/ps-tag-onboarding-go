package domain

import (
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
			validationError: []string{ErrorNameRequired},
		},
		{
			name: "Missing user last name",
			user: User{
				FirstName: "ann",
				Email:     "a@a.com",
				Age:       22,
			},
			validationError: []string{ErrorNameRequired},
		},
		{
			name: "Missing user email",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Age:       22,
			},
			validationError: []string{ErrorEmailRequired},
		},
		{
			name: "User email not in a proper format",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "aa.com",
				Age:       18,
			},
			validationError: []string{ErrorEmailFormat},
		},
		{
			name: "Minimum age required",
			user: User{
				FirstName: "a",
				LastName:  "ann",
				Email:     "a@a.com",
				Age:       17,
			},
			validationError: []string{ErrorAgeMinimum},
		},
		{
			name: "User fails validation on multiple fields",
			user: User{
				LastName: "ann",
				Email:    "aa.com",
				Age:      17,
			},
			validationError: []string{ErrorAgeMinimum, ErrorEmailFormat, ErrorNameRequired},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			validationErrs := tc.user.Validate()
			assert.Equal(t, tc.validationError, validationErrs)
		})
	}
}
