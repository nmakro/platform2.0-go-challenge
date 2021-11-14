package maprepo

type ErrDuplicateEntry struct {
}

func NewDuplicateEntryError() *ErrDuplicateEntry {
	return &ErrDuplicateEntry{}
}

func (e ErrDuplicateEntry) Error() string {
	return "duplicate entry error"
}

type ErrEntityNotFound struct {
}

func NewEntityNotFoundError() *ErrEntityNotFound {
	return &ErrEntityNotFound{}
}

func (e ErrEntityNotFound) Error() string {
	return "entity not found error"
}
