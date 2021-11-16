package maprepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/nmakro/platform2.0-go-challenge/internal/app"
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
		errMsg := fmt.Sprintf("asset insight with id: %v already exists", insight.ID)
		return app.NewDuplicateEntryError(errMsg)
	}
	return nil
}

func (a *InsightRepo) Get(ctx context.Context, insightID uint32) (assets.Insight, error) {
	key := fmt.Sprintf("%d", insightID)
	v, err := a.insightConn.Get(key)
	var notFound *ErrNotFound

	switch {
	case err != nil && errors.As(err, &notFound):
		errMsg := fmt.Sprintf("asset insight with id: %v was not found", insightID)
		return assets.Insight{}, app.NewEntityNotFoundError(errMsg)
	case err != nil && !errors.As(err, &notFound):
		errMsg := "unknown internal error"
		return assets.Insight{}, NewInternalRepositoryError(errMsg, nil)
	}

	if ins, ok := v.(assets.Insight); ok {
		return ins, nil
	}

	errMsg := fmt.Sprintf("error while reading insight with id: %d from db", insightID)
	return assets.Insight{}, NewInternalRepositoryError(errMsg, nil)
}

func (a *InsightRepo) GetMany(ctx context.Context, insightIDs []uint32) ([]assets.Insight, error) {
	res := make([]assets.Insight, 0, len(insightIDs))
	for _, id := range insightIDs {
		a, err := a.Get(ctx, id)
		if err != nil {
			var notFound *app.ErrEntityNotFound
			if errors.As(err, &notFound) {
				continue
			}
			return []assets.Insight{}, fmt.Errorf("error while reading insights: %w", err)
		}
		res = append(res, a)
	}
	return res, nil
}

func (a *InsightRepo) List(ctx context.Context) ([]assets.Insight, error) {
	keys := a.insightConn.Keys()
	res := make([]assets.Insight, 0, len(keys))
	var notFound *ErrNotFound
	for i := 0; i < len(keys); i++ {
		v, err := a.insightConn.Get(keys[i])
		if err != nil {
			if errors.As(err, &notFound) {
				continue
			}
			return []assets.Insight{}, fmt.Errorf("error while reading insights: %w", err)
		}

		aud, ok := v.(assets.Insight)
		if !ok {
			errMsg := fmt.Sprintf("error while reading insight with id: %s from db", keys[i])
			return []assets.Insight{}, NewInternalRepositoryError(errMsg, nil)

		}
		res = append(res, aud)
	}
	return res, nil
}

func (i *InsightRepo) Delete(ctx context.Context, insightID uint32) error {
	key := fmt.Sprintf("%d", insightID)
	_, exists := i.insightConn.Delete(key)
	if !exists {
		return NewNotFoundError()
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
	var notFound *ErrNotFound

	if err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("error while reading starred audiences for user with email: %s", userEmail)
	}

	if starred, ok := v.([]uint32); ok {
		starred = append(starred, insightID)
		a.starConn.Upsert(userEmail, starred)
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
	var notFound *app.ErrEntityNotFound

	switch {
	case err != nil && errors.As(err, &notFound):
		return fmt.Errorf("cannot find starred insights assets for user: %s", userEmail)
	case err != nil && !errors.As(err, &notFound):
		return err
	default:
		if starred, ok := v.([]uint32); ok {
			found := false
			for i := range starred {
				if starred[i] == insightID {
					starred[i] = starred[len(starred)-1]
					found = true
					break
				}
			}
			if found {
				i.starConn.Upsert(userEmail, starred[:len(starred)-1])
				return nil
			}
		}
	}
	return fmt.Errorf("cannot find starred insight asset: %v for user: %s", insightID, userEmail)
}

func (i *InsightRepo) GetStarredIDsForUser(ctx context.Context, userEmail string) ([]uint32, error) {
	if userEmail == "" {
		return []uint32{}, fmt.Errorf("user email cannot be empty")
	}

	v, err := i.starConn.Get(userEmail)
	if err != nil {
		return nil, err
	}

	starred, ok := v.([]uint32)
	if !ok {
		return []uint32{}, fmt.Errorf("error while reading starred insights for user: %s", userEmail)
	}
	return starred, nil
}
