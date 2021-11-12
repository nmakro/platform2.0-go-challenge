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

func (u *UserDBRepo) AddUser(ctx context.Context, usr user.User) error {
	if usr.Email == "" {
		err := user.NewEmailMissingError()
		return err
	}

	if u.conn.Insert(usr.Email, usr) {
	} else {
		return fmt.Errorf("user email already exists")
	}
	return nil
}

func (u *UserDBRepo) GetByEmail(ctx context.Context, email string) (user.User, error) {
	res, exists := u.conn.Get(email)

	if exists {
		if v, ok := res.(user.User); ok {
			return v, nil
		}
		return user.User{}, fmt.Errorf("error in reading user")
	}
	return user.User{}, fmt.Errorf("user not found")
}

func (u *UserDBRepo) Delete(ctx context.Context, userEmail string) (user.User, error) {
	v, exists := u.conn.Delete(userEmail)

	if exists {
		if usr, ok := v.(user.User); ok {
			return usr, nil
		} else {
			return user.User{}, fmt.Errorf("error while reading user from db in user delete")
		}
	}
	return user.User{}, nil
}
