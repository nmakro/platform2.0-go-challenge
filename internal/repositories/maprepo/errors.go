package maprepo

type ErrNotFound struct {
}

func NewNotFoundError() *ErrNotFound {
	return &ErrNotFound{}
}

func (e ErrNotFound) Error() string {
	return ""
}

type ErrInternalRepository struct {
	inner error
	msg   string
}

func NewInternalRepositoryError(msg string, inner error) *ErrInternalRepository {
	return &ErrInternalRepository{msg: msg, inner: inner}
}

func (e ErrInternalRepository) Error() string {
	return e.msg
}

func (r *ErrInternalRepository) Unwrap() error {
	return r.inner
}

const UnknownError = "unknown internal repository error"
