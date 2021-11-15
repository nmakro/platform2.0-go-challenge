package assets

type ErrValidation struct {
	msg string
}

// Error function.
func (e *ErrValidation) Error() string {
	return e.msg
}

func NewErrValidation(msg string) *ErrValidation {
	return &ErrValidation{msg: msg}
}
