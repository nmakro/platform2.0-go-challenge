package maprepo

import (
	"github.com/nmakro/platform2.0-go-challenge/internal/app/user"

	cmap "github.com/orcaman/concurrent-map"
)

type UserDBClient struct {
	userDB cmap.ConcurrentMap
}

func NewUserDBClient() *UserDBClient {
	return &UserDBClient{
		userDB: cmap.New(),
	}
}

func (u *UserDBClient) AddUser(user user.User) error {
	u.userDB.Set(user.Email, user)
	return nil
}

func (u *UserDBClient) DeleteUser(user user.User) error {
	u.userDB.Remove(user.Email)
	return nil
}
