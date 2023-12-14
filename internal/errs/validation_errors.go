package errs

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Err     string   `json:"error"`
	Details []string `json:"details"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%v: %s", e.Err, strings.Join(e.Details, ", "))
}
