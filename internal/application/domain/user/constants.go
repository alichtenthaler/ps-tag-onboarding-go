package domain

const (
	ErrorAgeMinimum          = "user does not meet minimum age requirement"
	ErrorEmailFormat         = "user email must be properly formatted"
	ErrorEmailRequired       = "user email required"
	ErrorNameRequired        = "user first/last names required"
	ErrorNameUnique          = "user with the same first and last name already exists"
	ResponseUserNotFound     = "user not found"
	ResponseValidationFailed = "User did not pass validation"
)
