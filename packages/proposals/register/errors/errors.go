package errors

import "fmt"

type MissingKeyError struct {
	Context string
	Field   string
}

func (e *MissingKeyError) Error() string {
	return fmt.Sprintf("[%s] Missing field: %s", e.Context, e.Field)
}

type InvalidFieldError struct {
	Field string
}

func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("Invalid field: %s", e.Field)
}

type EmptyFieldError struct {
	Field string
}

func (e *EmptyFieldError) Error() string {
	return fmt.Sprintf("Empty field: %s", e.Field)
}

func NewErrorResponse(message string) map[string]any {
	return map[string]any{
		"body": map[string]string{
			"error": message,
		},
	}
}
