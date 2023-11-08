package service

import (
	"github.com/alichtenthaler/ps-tag-onboarding-go/api/src/application/domain/user"
	"strings"
)

type ValidationError struct {
	Error   string   `json:"error"`
	Details []string `json:"details"`
}

func (s *CreateUserService) validate(user domain.User) []string {
	var errs []string

	if err := s.validateAge(user); err != "" {
		errs = append(errs, err)
	}

	if err := s.validateEmail(user); err != "" {
		errs = append(errs, err)
	}

	if err := s.validateName(user); err != "" {
		errs = append(errs, err)
	}

	return errs
}

func (s *CreateUserService) validateAge(user domain.User) string {
	if user.Age < 18 {
		return domain.ErrorAgeMinimum
	}

	return ""
}

func (s *CreateUserService) validateEmail(user domain.User) string {
	if user.Email == "" {
		return domain.ErrorEmailRequired
	}

	if !strings.Contains(user.Email, "@") {
		return domain.ErrorEmailFormat
	}

	return ""
}

func (s *CreateUserService) validateName(user domain.User) string {
	if user.FirstName == "" || user.LastName == "" {
		return domain.ErrorNameRequired
	}

	return ""
}


