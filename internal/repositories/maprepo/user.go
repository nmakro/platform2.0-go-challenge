package maprepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/user"
	"github.com/nmakro/platform2.0-go-challenge/pkg/security"
)

type UserDBRepo struct {
	conn *Client
}

func NewUserRepo(client *Client) *UserDBRepo {
	return &UserDBRepo{
		conn: client,
	}
}

func (u *UserDBRepo) Add(ctx context.Context, usr user.AddUserCommand) error {
	if usr.Email == "" {
		err := user.NewErrValidation("user email cannot be empty")
		return err
	}

	hashed, err := security.HashPassword(usr.Password)
	if err != nil {
		errMsg := "error hashing user password"
		return NewInternalRepositoryError(errMsg, err)
	}

	usr.Password = hashed
	if ok := u.conn.Insert(usr.Email, usr); !ok {
		errMsg := fmt.Sprintf("user with email: %s already exists", usr.Email)
		return app.NewDuplicateEntryError(errMsg)
	}
	return nil
}

func (u *UserDBRepo) GetByEmail(ctx context.Context, email string) (user.User, error) {
	res, err := u.conn.Get(email)

	if err != nil {
		var notFound *ErrNotFound
		if errors.As(err, &notFound) {
			errMsg := fmt.Sprintf("user with email: %s not found", email)
			return user.User{}, app.NewEntityNotFoundError(errMsg)
		}
		errMsg := fmt.Sprintf("unknown internal error while getting user with email: %s", email)
		return user.User{}, NewInternalRepositoryError(errMsg, nil)
	}

	v, ok := res.(user.AddUserCommand)
	if !ok {
		errMsg := fmt.Sprintf("unknown internal error while reading user with email: %s", email)
		return user.User{}, NewInternalRepositoryError(errMsg, nil)
	}
	user := user.User{
		UserName:  v.UserName,
		FirstName: v.FirstName,
		LastName:  v.LastName,
		Email:     v.Email,
	}
	return user, nil
}

func (u *UserDBRepo) GetWithPassword(ctx context.Context, email string) (user.AddUserCommand, error) {
	res, err := u.conn.Get(email)

	if err != nil {
		var notFound *ErrNotFound
		if errors.As(err, &notFound) {
			errMsg := fmt.Sprintf("user with email: %s not found", email)
			return user.AddUserCommand{}, app.NewEntityNotFoundError(errMsg)
		}
		errMsg := fmt.Sprintf("unknown internal error while getting user with email: %s", email)
		return user.AddUserCommand{}, NewInternalRepositoryError(errMsg, nil)
	}

	v, ok := res.(user.AddUserCommand)
	if !ok {
		errMsg := fmt.Sprintf("unknown internal error while reading user with email: %s", email)
		return user.AddUserCommand{}, NewInternalRepositoryError(errMsg, nil)
	}

	return v, nil
}

func (u *UserDBRepo) Delete(ctx context.Context, userEmail string) error {
	_, exists := u.conn.Delete(userEmail)

	if !exists {
		errMsg := fmt.Sprintf("user with email: %s not found", userEmail)
		return app.NewEntityNotFoundError(errMsg)
	}
	return nil
}

func (u *UserDBRepo) List(ctx context.Context) ([]user.User, error) {
	keys := u.conn.Keys()
	res := make([]user.User, 0, len(keys))
	var notFound *ErrNotFound
	for i := 0; i < len(keys); i++ {
		v, err := u.conn.Get(keys[i])
		if err != nil {
			if errors.As(err, &notFound) {
				continue
			}
			return []user.User{}, fmt.Errorf("error while reading users: %w", err)
		}

		usr, ok := v.(user.AddUserCommand)
		if !ok {
			errMsg := fmt.Sprintf("error while reading user with email: %s from db", keys[i])
			return []user.User{}, NewInternalRepositoryError(errMsg, nil)

		}
		res = append(res, user.User{UserName: usr.UserName, Email: usr.Email, FirstName: usr.FirstName, LastName: usr.FirstName})
		fmt.Println(res)
	}
	return res, nil
}
