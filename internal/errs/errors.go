package errs

import "fmt"

// Error represents an application error
type Error struct {
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

// GenericError represents a simple error to be returned in the API
type GenericError struct {
	Err string `json:"error"`
}

func (e GenericError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}
