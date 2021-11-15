package maprepo

import (
	"context"
	"fmt"

	"github.com/nmakro/platform2.0-go-challenge/internal/app/user"
)

type UserDBRepo struct {
	conn *Client
}

func NewUserRepo(client *Client) *UserDBRepo {
	return &UserDBRepo{
		conn: client,
	}
}

func (u *UserDBRepo) Add(ctx context.Context, usr user.User) error {
	if usr.Email == "" {
		err := user.NewEmailMissingError()
		return err
	}

	if ok := u.conn.Insert(usr.Email, usr); !ok {
		return fmt.Errorf("user email already exists")
	}
	return nil
}

func (u *UserDBRepo) GetByEmail(ctx context.Context, email string) (user.User, error) {
	res, err := u.conn.Get(email)

	if err != nil {
		return user.User{}, err
	}

	v, ok := res.(user.User)
	if !ok {
		return user.User{}, fmt.Errorf("error in get user, read unexpected value %v", res)
	}
	return v, nil
}

func (u *UserDBRepo) Delete(ctx context.Context, userEmail string) error {
	v, exists := u.conn.Delete(userEmail)

	if exists {
		if _, ok := v.(user.User); ok {
			return nil
		} else {
			return fmt.Errorf("error while reading user from db in user delete")
		}
	}
	return nil
}