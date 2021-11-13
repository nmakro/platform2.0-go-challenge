package maprepo

import (
	"context"
	"errors"
	"fmt"

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

func (c *ChartRepo) AddChart(ctx context.Context, chart assets.Chart) error {
	err := assets.ValidateChart(chart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", chart.ID)
	ok := c.chartConn.Insert(key, chart)
	if !ok {
		return NewDuplicateEntryError()
	}
	return nil
}

func (c *ChartRepo) Get(ctx context.Context, chartID uint32) (assets.Chart, error) {
	key := fmt.Sprintf("%d", chartID)
	v, err := c.chartConn.Get(key)
	if err != nil {
		return assets.Chart{}, err
	}
	if chart, ok := v.(assets.Chart); ok {
		return chart, nil
	}
	return assets.Chart{}, fmt.Errorf("error while reading chart from db.")
}

func (a *ChartRepo) Delete(ctx context.Context, chartID uint32) error {
	key := fmt.Sprintf("%d", chartID)
	_, exists := a.chartConn.Delete(key)
	if !exists {
		return NewEntityNotFoundError()
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
	notFound := NewEntityNotFoundError()
	if err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("error while reading stared charts for user with email: %s", userEmail)
	}

	if stared, ok := v.([]uint32); ok {
		stared = append(stared, chartID)
		c.starConn.Upsert(userEmail, stared)
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
	notFound := NewEntityNotFoundError()

	switch {
	case err != nil && errors.As(err, &notFound):
		return fmt.Errorf("cannot find stared audience assets for user: %s", userEmail)
	case err != nil && errors.As(err, &notFound):
		return err
	default:
		if stared, ok := v.([]uint32); ok {
			found := false
			for i := range stared {
				if stared[i] == chartID {
					stared[i] = stared[len(stared)-1]
					found = true
					break
				}
			}
			if found {
				c.starConn.Upsert(userEmail, stared[:len(stared)-1])
				return nil
			}
		}
	}
	return fmt.Errorf("cannot find stared chart asset: %v for user: %s", chartID, userEmail)
}

func (c *ChartRepo) GetStaredChartIDsForUser(ctx context.Context, userEmail string) ([]uint32, error) {
	if userEmail == "" {
		return []uint32{}, fmt.Errorf("user email cannot be empty")
	}

	v, err := c.starConn.Get(userEmail)
	if err != nil {
		return nil, err
	}

	stared, ok := v.([]uint32)
	if !ok {
		return []uint32{}, fmt.Errorf("error while reading stared charts for user: %s", userEmail)
	}
	return stared, nil
}
