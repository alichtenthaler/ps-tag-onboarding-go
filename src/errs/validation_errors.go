package errs

// ValidationError represents a validation error
type ValidationError struct {
	Error   string   `json:"error"`
	Details []string `json:"details"`
}
