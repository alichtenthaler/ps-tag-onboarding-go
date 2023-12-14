package errs

import "fmt"

// Error represents an application error
type Error struct {
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
