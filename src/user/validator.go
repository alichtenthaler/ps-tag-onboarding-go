package user

import (
	"strings"
)

func (u *User) validate() []string {
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

	if !strings.Contains(u.Email, "@") {
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


