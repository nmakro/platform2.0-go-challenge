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

	u := user.User{
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

	u := user.User{
		UserName: "test",
		Password: "@!3aER&4!",
	}

	err := userRepo.Add(context.Background(), u)

	expected := user.NewEmailMissingError()
	assert.ErrorAs(t, err, &expected)
}

func TestGetUser(t *testing.T) {
	userRepo, down := SetUp()

	defer down()

	u := user.User{
		UserName: "test",
		Email:    "test@host.com",
		Password: "@!3aER&4!",
	}

	err := userRepo.Add(context.Background(), u)
	assert.NoError(t, err)

	usr, err := userRepo.GetByEmail(context.Background(), u.Email)

	assert.NoError(t, err)
	assert.Equal(t, u, usr)
}

func SetUp() (*repo.UserDBRepo, func()) {
	dbClient := repo.NewClient()
	userRepo := repo.NewUserRepo(dbClient)

	tearDown := func() {
		dbClient.ClearAll()
	}

	return userRepo, tearDown
}
