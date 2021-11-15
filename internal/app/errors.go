package app

type ErrDuplicateEntry struct {
	msg string
}

func NewDuplicateEntryError(msg string) *ErrDuplicateEntry {
	return &ErrDuplicateEntry{msg: msg}
}

func (e ErrDuplicateEntry) Error() string {
	return e.msg
}

type ErrEntityNotFound struct {
	inner error
	msg   string
}

func NewEntityNotFoundError(msg string) *ErrEntityNotFound {
	return &ErrEntityNotFound{msg: msg}
}

func (e ErrEntityNotFound) Error() string {
	return e.msg
}

func (e ErrEntityNotFound) Unwrap() error {
	return e.inner
}
