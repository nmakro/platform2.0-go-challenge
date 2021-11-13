//go:build integration_test

package maprepo_test

import (
	"context"
	"testing"

	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	repo "github.com/nmakro/platform2.0-go-challenge/internal/app/repositories/maprepo"
	"github.com/stretchr/testify/assert"
)

func TestSaveInsight(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	ins := assets.Insight{
		ID:    123,
		Topic: "sales",
	}

	err := insightRepo.AddInsight(context.Background(), ins)

	assert.NoError(t, err)
}

func TestSaveInsightNoID(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	ins := assets.Insight{
		ID:    0,
		Topic: "sales",
	}

	err := insightRepo.AddInsight(context.Background(), ins)

	expectedErr := assets.NewAssetNoIDError()
	assert.ErrorAs(t, err, &expectedErr)
}

func TestGetInsight(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	ins := assets.Insight{
		ID:    123,
		Topic: "sales",
	}

	err := insightRepo.AddInsight(context.Background(), ins)
	assert.NoError(t, err)

	res, err := insightRepo.Get(context.Background(), ins.ID)
	assert.NoError(t, err)
	assert.Equal(t, ins, res)
}

func TestDeleteInsight(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	ins := assets.Insight{
		ID:    123,
		Topic: "sales",
	}

	err := insightRepo.AddInsight(context.Background(), ins)
	assert.NoError(t, err)

	err = insightRepo.Delete(context.Background(), ins.ID)
	assert.NoError(t, err)

	_, err = insightRepo.Get(context.Background(), ins.ID)
	var notFound *repo.ErrEntityNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestUpdateInsight(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	ins := assets.Insight{
		ID:    123,
		Topic: "sales",
	}
	err := insightRepo.AddInsight(context.Background(), ins)
	assert.NoError(t, err)

	ins.Topic = "new topic"

	err = insightRepo.Update(context.Background(), ins)
	assert.NoError(t, err)

	expected, err := insightRepo.Get(context.Background(), ins.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.Topic, ins.Topic)
}

func TestStarInsightForUser(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	userEmail := "test@host.com"
	insightID := uint32(123)

	err := insightRepo.Star(context.Background(), userEmail, insightID)
	assert.NoError(t, err)

	stared, err := insightRepo.GetStaredInsightsIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	expected := []uint32{123}
	assert.Equal(t, expected, stared)
}

func TestUnStarInsightForUser(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	userEmail := "test@host.com"
	insightID1 := uint32(123)

	err := insightRepo.Star(context.Background(), userEmail, insightID1)
	assert.NoError(t, err)

	insightID2 := uint32(456)
	err = insightRepo.Star(context.Background(), userEmail, insightID2)
	assert.NoError(t, err)

	err = insightRepo.Unstar(context.Background(), userEmail, insightID1)
	assert.NoError(t, err)

	stared, err := insightRepo.GetStaredInsightsIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	expected := []uint32{insightID2}
	assert.Equal(t, expected, stared)
}

func TestGetStaredInsightIDsForUser(t *testing.T) {
	insightRepo, down := SetUpInsight()

	defer down()

	userEmail := "test@host.com"
	insightID1 := uint32(1)
	insightID2 := uint32(2)
	insightID3 := uint32(3)

	err := insightRepo.Star(context.Background(), userEmail, insightID1)
	assert.NoError(t, err)

	err = insightRepo.Star(context.Background(), userEmail, insightID2)
	assert.NoError(t, err)

	err = insightRepo.Star(context.Background(), userEmail, insightID3)
	assert.NoError(t, err)

	res, err := insightRepo.GetStaredInsightsIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	assert.Equal(t, []uint32{insightID1, insightID2, insightID3}, res)
}

func SetUpInsight() (*repo.InsightRepo, func()) {
	insightClient := repo.NewClient()
	starClient := repo.NewClient()
	audienceRepo := repo.NewAInsightRepo(insightClient, starClient)

	tearDown := func() {
		insightClient.ClearAll()
		starClient.ClearAll()
	}

	return audienceRepo, tearDown
}
