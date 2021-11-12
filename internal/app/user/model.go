package user

import (
	"fmt"
	"net/mail"
	"unicode"
)

type User struct {
	UserName  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"-"`
}

type userValidator = func(u User) error

func validateEmail(u User) error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return fmt.Errorf("user email %s is not valid", u.Email)
	}
	return nil
}

func validateName(u User) error {
	msg := "only letters are allowed in first and last name"
	for _, v := range u.FirstName {
		if !unicode.IsLetter(v) {
			return fmt.Errorf(msg)
		}
	}

	for _, v := range u.LastName {
		if !unicode.IsLetter(v) {
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func validatePassword(u User) error {
	if len(u.Password) < 8 {
		return fmt.Errorf("password must be at least eight characters long")
	}
	hasUpper := false
	hasDigit := false
	for _, v := range u.Password {
		if unicode.IsDigit(v) {
			hasDigit = true
		}
		if unicode.IsUpper(v) {
			hasUpper = true
		}
	}
	if !hasDigit {
		return fmt.Errorf("password must have at least one digit")
	}
	if !hasUpper {
		return fmt.Errorf("password must have at least one upper letter")
	}
	return nil
}

var userValidations = []userValidator{
	validateEmail,
	validateName,
	validatePassword,
}

func ValidateUser(u User) error {
	for _, f := range userValidations {
		if err := f(u); err != nil {
			return err
		}
	}
	return nil
}
