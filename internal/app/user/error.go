package user

type EmailMissing struct {
}

// Error function.
func (r EmailMissing) Error() string {
	return "user email cannot be empty"
}

func NewEmailMissingError() *EmailMissing {
	return &EmailMissing{}
}
