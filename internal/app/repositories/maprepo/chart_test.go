//go:build integration_test

package maprepo_test

import (
	"context"
	"testing"

	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	repo "github.com/nmakro/platform2.0-go-challenge/internal/app/repositories/maprepo"
	"github.com/stretchr/testify/assert"
)

func TestSaveChart(t *testing.T) {
	chartRepo, down := SetUpChart()

	defer down()

	chart := assets.Chart{
		ID:    123,
		Title: "test chart",
	}

	err := chartRepo.AddChart(context.Background(), chart)

	assert.NoError(t, err)
}

func TestSaveChartNoID(t *testing.T) {
	chartRepo, down := SetUpChart()

	defer down()

	chart := assets.Chart{
		ID:    0,
		Title: "test chart",
	}

	err := chartRepo.AddChart(context.Background(), chart)

	expectedErr := assets.NewAssetNoIDError()
	assert.ErrorAs(t, err, &expectedErr)
}

func TestGetChart(t *testing.T) {
	chartRepo, down := SetUpChart()

	defer down()

	chart := assets.Chart{
		ID:    123,
		Title: "test chart",
	}

	err := chartRepo.AddChart(context.Background(), chart)
	assert.NoError(t, err)

	res, err := chartRepo.Get(context.Background(), chart.ID)
	assert.NoError(t, err)
	assert.Equal(t, chart, res)
}

func TestDeleteChart(t *testing.T) {
	chartRepo, down := SetUpChart()

	defer down()

	chart := assets.Chart{
		ID:    123,
		Title: "test chart",
	}

	err := chartRepo.AddChart(context.Background(), chart)
	assert.NoError(t, err)

	err = chartRepo.Delete(context.Background(), chart.ID)
	assert.NoError(t, err)

	_, err = chartRepo.Get(context.Background(), chart.ID)
	var notFound *repo.ErrEntityNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestStarChartForUser(t *testing.T) {
	chartRepo, down := SetUpChart()

	defer down()

	userEmail := "test@host.com"
	chartID := uint32(123)

	err := chartRepo.Star(context.Background(), userEmail, chartID)
	assert.NoError(t, err)

	stared, err := chartRepo.GetStaredChartIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	expected := []uint32{123}
	assert.Equal(t, expected, stared)
}

func TestUnStarChartForUser(t *testing.T) {
	chartRepo, down := SetUpChart()

	defer down()

	userEmail := "test@host.com"
	chartID1 := uint32(123)

	err := chartRepo.Star(context.Background(), userEmail, chartID1)
	assert.NoError(t, err)

	chartID2 := uint32(456)
	err = chartRepo.Star(context.Background(), userEmail, chartID2)
	assert.NoError(t, err)

	err = chartRepo.Unstar(context.Background(), userEmail, chartID1)
	assert.NoError(t, err)

	stared, err := chartRepo.GetStaredChartIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	expected := []uint32{chartID2}
	assert.Equal(t, expected, stared)
}

func TestGetStaredChartsIDsForUser(t *testing.T) {
	chartRepo, down := SetUpChart()

	defer down()

	userEmail := "test@host.com"
	chartID1 := uint32(1)
	chartID2 := uint32(2)
	chartID3 := uint32(3)

	err := chartRepo.Star(context.Background(), userEmail, chartID1)
	assert.NoError(t, err)

	err = chartRepo.Star(context.Background(), userEmail, chartID2)
	assert.NoError(t, err)

	err = chartRepo.Star(context.Background(), userEmail, chartID3)
	assert.NoError(t, err)

	res, err := chartRepo.GetStaredChartIDsForUser(context.Background(), userEmail)
	assert.NoError(t, err)
	assert.Equal(t, []uint32{chartID1, chartID2, chartID3}, res)
}

func SetUpChart() (*repo.ChartRepo, func()) {
	chartClient := repo.NewClient()
	starClient := repo.NewClient()
	chartRepo := repo.NewChartRepo(chartClient, starClient)

	tearDown := func() {
		chartClient.ClearAll()
		starClient.ClearAll()
	}

	return chartRepo, tearDown
}
