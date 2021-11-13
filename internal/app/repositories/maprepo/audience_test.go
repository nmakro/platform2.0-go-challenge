//go:build integration_test

package maprepo_test

import (
	"context"
	"testing"

	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	repo "github.com/nmakro/platform2.0-go-challenge/internal/app/repositories/maprepo"
	"github.com/stretchr/testify/assert"
)

func TestSaveAudience(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	aud := assets.Audience{
		ID:               153,
		SocialMediaHours: 168,
		Gender:           assets.Male,
	}

	err := audienceRepo.AddAudience(context.Background(), aud)

	assert.NoError(t, err)
}

func TestSaveAudienceNoID(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	aud := assets.Audience{
		ID:               0,
		SocialMediaHours: 168,
		Gender:           assets.Male,
	}

	err := audienceRepo.AddAudience(context.Background(), aud)

	expectedErr := assets.NewAssetNoIDError()
	assert.ErrorAs(t, err, &expectedErr)
}

func TestGetAudience(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	aud := assets.Audience{
		ID:               1,
		SocialMediaHours: 168,
		Gender:           assets.Male,
	}

	err := audienceRepo.AddAudience(context.Background(), aud)
	assert.NoError(t, err)

	res, err := audienceRepo.Get(context.Background(), aud.ID)
	assert.NoError(t, err)
	assert.Equal(t, aud, res)
}

func TestDeleteAudience(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	aud := assets.Audience{
		ID:               142,
		SocialMediaHours: 168,
		Gender:           assets.Male,
	}

	err := audienceRepo.AddAudience(context.Background(), aud)
	assert.NoError(t, err)

	err = audienceRepo.Delete(context.Background(), aud.ID)
	assert.NoError(t, err)

	_, err = audienceRepo.Get(context.Background(), aud.ID)
	var notFound *repo.ErrEntityNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestUpdateAudience(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	aud := assets.Audience{
		ID:               142,
		SocialMediaHours: 168,
		Gender:           assets.Male,
	}

	err := audienceRepo.AddAudience(context.Background(), aud)
	assert.NoError(t, err)

	aud.Gender = assets.Female

	err = audienceRepo.Update(context.Background(), aud)
	assert.NoError(t, err)

	expected, err := audienceRepo.Get(context.Background(), aud.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.Gender, assets.Female)
}

func TestStarAudienceForUser(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	userEmail := "test@host.com"
	audienceID := uint32(123)

	err := audienceRepo.Star(context.Background(), userEmail, audienceID)
	assert.NoError(t, err)

	stared, err := audienceRepo.GetStaredAudienceIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	expected := []uint32{123}
	assert.Equal(t, expected, stared)
}

func TestUnStarAudienceForUser(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	userEmail := "test@host.com"
	audienceID1 := uint32(123)

	err := audienceRepo.Star(context.Background(), userEmail, audienceID1)
	assert.NoError(t, err)

	audienceID2 := uint32(456)
	err = audienceRepo.Star(context.Background(), userEmail, audienceID2)
	assert.NoError(t, err)

	err = audienceRepo.Unstar(context.Background(), userEmail, audienceID1)
	assert.NoError(t, err)

	stared, err := audienceRepo.GetStaredAudienceIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	expected := []uint32{audienceID2}
	assert.Equal(t, expected, stared)
}

func TestGetStaredAudienceIDsForUser(t *testing.T) {
	audienceRepo, down := SetUpAudience()

	defer down()

	userEmail := "test@host.com"
	audienceID1 := uint32(1)
	audienceID2 := uint32(2)
	audienceID3 := uint32(3)

	err := audienceRepo.Star(context.Background(), userEmail, audienceID1)
	assert.NoError(t, err)

	err = audienceRepo.Star(context.Background(), userEmail, audienceID2)
	assert.NoError(t, err)

	err = audienceRepo.Star(context.Background(), userEmail, audienceID3)
	assert.NoError(t, err)

	res, err := audienceRepo.GetStaredAudienceIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	assert.Equal(t, []uint32{audienceID1, audienceID2, audienceID3}, res)
}

func SetUpAudience() (*repo.AudienceRepo, func()) {
	audClient := repo.NewClient()
	starClient := repo.NewClient()
	audienceRepo := repo.NewAudienceRepo(audClient, starClient)

	tearDown := func() {
		audClient.ClearAll()
		starClient.ClearAll()
	}

	return audienceRepo, tearDown
}
