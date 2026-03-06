package command

type ValidationError struct {
	Code    string
	Message string
}

func NewValidationError(code, message string) *ValidationError {
	return &ValidationError{
		Code:    code,
		Message: message,
	}
}

func (e *ValidationError) Error() string {
	return e.Message
}
