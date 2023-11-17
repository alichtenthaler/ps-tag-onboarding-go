package domain

import (
	"regexp"
)

type ValidationError struct {
	Err string `json:"error"`
	Details []string `json:"details"`
}

func (ve ValidationError) Error() string {
	return ve.Err
}

func (u *User) Validate() []string {
	var errs []string

	if err := u.validateAge(); err != "" {
		errs = append(errs, err)
	}

	if err := u.validateEmail(); err != "" {
		errs = append(errs, err)
	}

	if err := u.validateName(); err != "" {
		errs = append(errs, err)
	}

	return errs
}

func (u *User) validateAge() string {
	if u.Age < 18 {
		return ErrorAgeMinimum
	}

	return ""
}

func (u *User) validateEmail() string {
	if u.Email == "" {
		return ErrorEmailRequired
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		return ErrorEmailFormat
	}

	return ""
}

func (u *User) validateName() string {
	if u.FirstName == "" || u.LastName == "" {
		return ErrorNameRequired
	}

	return ""
}


