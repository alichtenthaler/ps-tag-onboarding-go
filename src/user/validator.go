package user

import (
	"context"
	"strings"
)

type ValidatorI interface {
	validate(ctx context.Context, user User) []string
}

type Validator struct {
	repo UserRepositoryI
}

func newValidator(repo *UserRepository) *Validator {
	return &Validator{
		repo: repo,
	}
}

func (v *Validator) validate(ctx context.Context, user User) []string {
	var errs []string

	if err := v.validateAge(user); err != "" {
		errs = append(errs, err)
	}

	if err := v.validateEmail(user); err != "" {
		errs = append(errs, err)
	}

	if err := v.validateName(ctx, user); err != "" {
		errs = append(errs, err)
	}

	return errs
}

func (v *Validator) validateAge(user User) string {
	if user.Age < 18 {
		return ErrorAgeMinimum
	}

	return ""
}

func (v *Validator) validateEmail(user User) string {
	if user.Email == "" {
		return ErrorEmailRequired
	}

	if !strings.Contains(user.Email, "@") {
		return ErrorEmailFormat
	}

	return ""
}

func (v *Validator) validateName(ctx context.Context, user User) string {
	if user.FirstName == "" || user.LastName == "" {
		return ErrorNameRequired
	}

	if v.repo.existsByFirstNameAndLastName(ctx, user.FirstName, user.LastName) {
		return ErrorNameUnique
	}

	return ""
}


