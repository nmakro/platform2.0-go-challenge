//go:build integration_test

package maprepo_test

import (
	"context"
	"testing"

	"github.com/nmakro/platform2.0-go-challenge/internal/app/user"
	repo "github.com/nmakro/platform2.0-go-challenge/internal/repositories/maprepo"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	userRepo, down := SetUp()

	defer down()

	u := user.AddUserCommand{
		UserName: "test",
		Email:    "test@host.com",
		Password: "@!3aER&4!",
	}

	err := userRepo.Add(context.Background(), u)

	assert.NoError(t, err)
}

func TestSaveUserEmptyEmail(t *testing.T) {
	userRepo, down := SetUp()

	defer down()

	u := user.AddUserCommand{
		UserName: "test",
		Password: "@!3aER&4!",
	}

	err := userRepo.Add(context.Background(), u)

	var expected *user.ErrValidation
	assert.ErrorAs(t, err, &expected)
}

func TestGetUser(t *testing.T) {
	userRepo, down := SetUp()

	defer down()

	u := user.AddUserCommand{
		UserName: "test",
		Email:    "test@host.com",
		Password: "@!3aER&4!",
	}

	err := userRepo.Add(context.Background(), u)
	assert.NoError(t, err)

	usr, err := userRepo.GetByEmail(context.Background(), u.Email)

	usrFromCmd := user.User{
		UserName: u.UserName,
		Email:    u.Email,
	}
	assert.NoError(t, err)
	assert.Equal(t, usrFromCmd, usr)
}

func SetUp() (*repo.UserDBRepo, func()) {
	dbClient := repo.NewClient()
	userRepo := repo.NewUserRepo(dbClient)

	tearDown := func() {
		dbClient.ClearAll()
	}

	return userRepo, tearDown
}
