package internal

import "fmt"

type AppError struct {
	msg string
}

// Error function.
func (r *AppError) Error() string {
	return fmt.Sprintf("error: `%s`", r.msg)
}

// NewDBError constructor.
func NewAppError(msg string, err error) *AppError {
	return &AppError{
		msg: msg,
	}
}
