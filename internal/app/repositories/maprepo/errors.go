package maprepo

import "fmt"

type DBError struct {
	inner error
	msg   string
}

// Error function.
func (r *DBError) Error() string {
	if r.inner != nil {
		return fmt.Sprintf("Repository error: `%s` with inner error `%s`", r.msg, r.inner.Error())
	}
	return r.msg
}

// NewDBError constructor.
func NewRepositoryError(msg string, err error) *DBError {
	return &DBError{
		msg:   msg,
		inner: err,
	}
}
