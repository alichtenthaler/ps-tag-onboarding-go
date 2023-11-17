package user

import "github.com/alichtenthaler/ps-tag-onboarding-go/api/src/errs"

var (
	ErrorAgeMinimum          = errs.Error{Message: "user does not meet minimum age requirement"}
	ErrorEmailFormat         = errs.Error{Message: "user email must be properly formatted"}
	ErrorEmailRequired       = errs.Error{Message: "user email required"}
	ErrorNameRequired        = errs.Error{Message: "user first/last names required"}
	ErrorNameUnique          = errs.Error{Message: "user with the same first and last name already exists"}
	ResponseUserNotFound     = errs.Error{Message: "user not found"}
	ResponseValidationFailed = errs.Error{Message: "User did not pass validation"}
)
