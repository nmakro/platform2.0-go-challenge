package maprepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
)

type ChartRepo struct {
	chartConn *Client
	starConn  *Client
}

func NewChartRepo(chartConn, starConn *Client) *ChartRepo {
	return &ChartRepo{
		chartConn: chartConn,
		starConn:  starConn,
	}
}

func (c *ChartRepo) Add(ctx context.Context, chart assets.Chart) error {
	err := assets.ValidateChart(chart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", chart.ID)
	ok := c.chartConn.Insert(key, chart)
	if !ok {
		errMsg := fmt.Sprintf("asset chart with id: %v already exists", chart.ID)
		return app.NewDuplicateEntryError(errMsg)
	}
	return nil
}

func (c *ChartRepo) Get(ctx context.Context, chartID uint32) (assets.Chart, error) {
	key := fmt.Sprintf("%d", chartID)
	v, err := c.chartConn.Get(key)
	var notFound *ErrNotFound

	switch {
	case err != nil && errors.As(err, &notFound):
		errMsg := fmt.Sprintf("asset chart with id: %v was not found", chartID)
		return assets.Chart{}, app.NewEntityNotFoundError(errMsg)
	case err != nil && !errors.As(err, &notFound):
		errMsg := "unknown internal error"
		return assets.Chart{}, NewInternalRepositoryError(errMsg)
	}

	if err != nil {
		return assets.Chart{}, err
	}
	if chart, ok := v.(assets.Chart); ok {
		return chart, nil
	}
	errMsg := fmt.Sprintf("error while reading chart with id: %d from db", chartID)
	return assets.Chart{}, NewInternalRepositoryError(errMsg)
}

func (c *ChartRepo) GetMany(ctx context.Context, chartIDs []uint32) ([]assets.Chart, error) {
	res := make([]assets.Chart, 0, len(chartIDs))
	for _, id := range chartIDs {
		chart, err := c.Get(ctx, id)
		if err != nil {
			var notFound *app.ErrEntityNotFound
			if errors.As(err, &notFound) {
				continue
			}
			return []assets.Chart{}, fmt.Errorf("error while reading charts: %w", err)
		}
		res = append(res, chart)
	}
	return res, nil
}

func (a *ChartRepo) Delete(ctx context.Context, chartID uint32) error {
	key := fmt.Sprintf("%d", chartID)
	_, exists := a.chartConn.Delete(key)
	if !exists {
		errMsg := fmt.Sprintf("asset chart with id: %v was not found", chartID)
		return app.NewEntityNotFoundError(errMsg)
	}
	return nil
}

func (a *ChartRepo) Update(ctx context.Context, chart assets.Chart) error {
	c, err := a.Get(ctx, chart.ID)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", c.ID)
	a.chartConn.Upsert(key, chart)
	return nil
}

func (c *ChartRepo) Star(ctx context.Context, userEmail string, chartID uint32) error {
	if chartID == 0 || userEmail == "" {
		return fmt.Errorf("chart id and user email cannot be empty")
	}

	v, err := c.starConn.Get(userEmail)
	var notFound *ErrNotFound

	if err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("error while reading starred charts for user with email: %s", userEmail)
	}

	if starred, ok := v.([]uint32); ok {
		starred = append(starred, chartID)
		c.starConn.Upsert(userEmail, starred)
	} else {
		res := make([]uint32, 0, 20)
		res = append(res, chartID)
		c.starConn.Upsert(userEmail, res)
	}

	return nil
}

func (c *ChartRepo) Unstar(ctx context.Context, userEmail string, chartID uint32) error {
	if chartID == 0 || userEmail == "" {
		return fmt.Errorf("chart id and user email cannot be empty")
	}

	v, err := c.starConn.Get(userEmail)
	var notFound *app.ErrEntityNotFound

	switch {
	case err != nil && errors.As(err, &notFound):
		return fmt.Errorf("cannot find starred audience assets for user: %s", userEmail)
	case err != nil && !errors.As(err, &notFound):
		return err
	default:
		if starred, ok := v.([]uint32); ok {
			found := false
			for i := range starred {
				if starred[i] == chartID {
					starred[i] = starred[len(starred)-1]
					found = true
					break
				}
			}
			if found {
				c.starConn.Upsert(userEmail, starred[:len(starred)-1])
				return nil
			}
		}
	}
	return fmt.Errorf("cannot find starred chart asset: %v for user: %s", chartID, userEmail)
}

func (c *ChartRepo) GetStarredIDsForUser(ctx context.Context, userEmail string) ([]uint32, error) {
	if userEmail == "" {
		return []uint32{}, fmt.Errorf("user email cannot be empty")
	}

	v, err := c.starConn.Get(userEmail)
	if err != nil {
		return nil, err
	}

	starred, ok := v.([]uint32)
	if !ok {
		return []uint32{}, fmt.Errorf("error while reading starred charts for user: %s", userEmail)
	}
	return starred, nil
}
