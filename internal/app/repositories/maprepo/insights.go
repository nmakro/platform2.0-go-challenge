package maprepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
)

type InsightRepo struct {
	insightConn *Client
	starConn    *Client
}

func NewAInsightRepo(insightConn, starConn *Client) *InsightRepo {
	return &InsightRepo{
		insightConn: insightConn,
		starConn:    starConn,
	}
}

func (i *InsightRepo) Add(ctx context.Context, insight assets.Insight) error {
	err := assets.ValidateInsightID(insight)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", insight.ID)
	ok := i.insightConn.Insert(key, insight)
	if !ok {
		return NewDuplicateEntryError()
	}
	return nil
}

func (a *InsightRepo) Get(ctx context.Context, insightID uint32) (assets.Insight, error) {
	key := fmt.Sprintf("%d", insightID)
	v, err := a.insightConn.Get(key)
	if err != nil {
		return assets.Insight{}, err
	}
	if ins, ok := v.(assets.Insight); ok {
		return ins, nil
	}
	return assets.Insight{}, fmt.Errorf("error while reading insight from db.")
}

func (i *InsightRepo) Delete(ctx context.Context, insightID uint32) error {
	key := fmt.Sprintf("%d", insightID)
	_, exists := i.insightConn.Delete(key)
	if !exists {
		return NewEntityNotFoundError()
	}
	return nil
}

func (i *InsightRepo) Update(ctx context.Context, insight assets.Insight) error {
	ins, err := i.Get(ctx, insight.ID)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", ins.ID)
	i.insightConn.Upsert(key, insight)
	return nil
}

func (a *InsightRepo) Star(ctx context.Context, userEmail string, insightID uint32) error {
	if insightID == 0 || userEmail == "" {
		return fmt.Errorf("insight id and user email cannot be empty")
	}

	v, err := a.starConn.Get(userEmail)
	notFound := NewEntityNotFoundError()
	if err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("error while reading stared audiences for user with email: %s", userEmail)
	}

	if stared, ok := v.([]uint32); ok {
		stared = append(stared, insightID)
		a.starConn.Upsert(userEmail, stared)
	} else {
		res := make([]uint32, 0, 20)
		res = append(res, insightID)
		a.starConn.Upsert(userEmail, res)
	}

	return nil
}

func (i *InsightRepo) Unstar(ctx context.Context, userEmail string, insightID uint32) error {
	if insightID == 0 || userEmail == "" {
		return fmt.Errorf("insight id and user email cannot be empty")
	}

	v, err := i.starConn.Get(userEmail)
	notFound := NewEntityNotFoundError()

	switch {
	case err != nil && errors.As(err, &notFound):
		return fmt.Errorf("cannot find stared insights assets for user: %s", userEmail)
	case err != nil && errors.As(err, &notFound):
		return err
	default:
		if stared, ok := v.([]uint32); ok {
			found := false
			for i := range stared {
				if stared[i] == insightID {
					stared[i] = stared[len(stared)-1]
					found = true
					break
				}
			}
			if found {
				i.starConn.Upsert(userEmail, stared[:len(stared)-1])
				return nil
			}
		}
	}
	return fmt.Errorf("cannot find stared insight asset: %v for user: %s", insightID, userEmail)
}

func (i *InsightRepo) GetStaredInsightsIDsForUser(ctx context.Context, userEmail string) ([]uint32, error) {
	if userEmail == "" {
		return []uint32{}, fmt.Errorf("user email cannot be empty")
	}

	v, err := i.starConn.Get(userEmail)
	if err != nil {
		return nil, err
	}

	stared, ok := v.([]uint32)
	if !ok {
		return []uint32{}, fmt.Errorf("error while reading stared insights for user: %s", userEmail)
	}
	return stared, nil
}
