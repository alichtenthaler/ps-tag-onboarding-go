package user

import (
	"strings"
)

func (up *Processor) validate(user User) []string {
	var errs []string

	if err := up.validateAge(user); err != "" {
		errs = append(errs, err)
	}

	if err := up.validateEmail(user); err != "" {
		errs = append(errs, err)
	}

	if err := up.validateName(user); err != "" {
		errs = append(errs, err)
	}

	return errs
}

func (up *Processor) validateAge(user User) string {
	if user.Age < 18 {
		return ErrorAgeMinimum
	}

	return ""
}

func (up *Processor) validateEmail(user User) string {
	if user.Email == "" {
		return ErrorEmailRequired
	}

	if !strings.Contains(user.Email, "@") {
		return ErrorEmailFormat
	}

	return ""
}

func (up *Processor) validateName(user User) string {
	if user.FirstName == "" || user.LastName == "" {
		return ErrorNameRequired
	}

	if up.existsByFirstNameAndLastName(user.FirstName, user.LastName) {
		return ErrorNameUnique
	}

	return ""
}


