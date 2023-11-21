package errs

var (
	ErrorAgeMinimum          = Error{Message: "user does not meet minimum age requirement"}
	ErrorEmailFormat         = Error{Message: "user email must be properly formatted"}
	ErrorEmailRequired       = Error{Message: "user email required"}
	ErrorNameRequired        = Error{Message: "user first/last names required"}
	ErrorNameUnique          = Error{Message: "user with the same first and last name already exists"}
	ResponseUserNotFound     = Error{Message: "user not found"}
	ResponseValidationFailed = Error{Message: "User did not pass validation"}
)
