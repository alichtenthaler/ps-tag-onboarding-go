package domain

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/internal/errs"
	"regexp"
)

func (u *User) Validate() errs.ValidationError {
	var errDetails []string

	if err := u.validateAge(); err != "" {
		errDetails = append(errDetails, err)
	}

	if err := u.validateEmail(); err != "" {
		errDetails = append(errDetails, err)
	}

	if err := u.validateName(); err != "" {
		errDetails = append(errDetails, err)
	}

	return errs.ValidationError{Err: errs.ResponseValidationFailed.Message, Details: errDetails}
}

func (u *User) validateAge() string {
	if u.Age < 18 {
		return errs.ErrorAgeMinimum.Error()
	}

	return ""
}

func (u *User) validateEmail() string {
	if u.Email == "" {
		return errs.ErrorEmailRequired.Error()
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		return errs.ErrorEmailFormat.Error()
	}

	return ""
}

func (u *User) validateName() string {
	if u.FirstName == "" || u.LastName == "" {
		return errs.ErrorNameRequired.Error()
	}

	return ""
}


