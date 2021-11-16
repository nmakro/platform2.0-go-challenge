package maprepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
)

type AudienceRepo struct {
	audConn  *Client
	starConn *Client
}

func NewAudienceRepo(audConn, starConn *Client) *AudienceRepo {
	return &AudienceRepo{
		audConn:  audConn,
		starConn: starConn,
	}
}

func (a *AudienceRepo) Add(ctx context.Context, audience assets.Audience) error {
	err := assets.ValidateAudience(audience)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", audience.ID)
	ok := a.audConn.Insert(key, audience)
	if !ok {
		errMsg := fmt.Sprintf("asset audience with id: %v already exists", audience.ID)
		return app.NewDuplicateEntryError(errMsg)
	}
	return nil
}

func (a *AudienceRepo) Get(ctx context.Context, audienceID uint32) (assets.Audience, error) {
	key := fmt.Sprintf("%d", audienceID)
	v, err := a.audConn.Get(key)
	var notFound *ErrNotFound

	switch {
	case err != nil && errors.As(err, &notFound):
		errMsg := fmt.Sprintf("asset audience with id: %v was not found", audienceID)
		return assets.Audience{}, app.NewEntityNotFoundError(errMsg)
	case err != nil && !errors.As(err, &notFound):
		errMsg := "unknown internal error"
		return assets.Audience{}, NewInternalRepositoryError(errMsg, nil)
	}

	if aud, ok := v.(assets.Audience); ok {
		return aud, nil
	}

	errMsg := fmt.Sprintf("error while reading audience with id: %d from db", audienceID)
	return assets.Audience{}, NewInternalRepositoryError(errMsg, nil)
}

func (a *AudienceRepo) GetMany(ctx context.Context, audienceIDs []uint32) ([]assets.Audience, error) {
	res := make([]assets.Audience, 0, len(audienceIDs))
	for _, id := range audienceIDs {
		a, err := a.Get(ctx, id)
		if err != nil {
			var notFound *app.ErrEntityNotFound
			if errors.As(err, &notFound) {
				continue
			}
			return []assets.Audience{}, fmt.Errorf("error while reading audiences: %w", err)
		}
		res = append(res, a)
	}
	return res, nil
}

func (a *AudienceRepo) List(ctx context.Context) ([]assets.Audience, error) {
	keys := a.audConn.Keys()
	res := make([]assets.Audience, 0, len(keys))
	var notFound *ErrNotFound
	for i := 0; i < len(keys); i++ {
		v, err := a.audConn.Get(keys[i])
		if err != nil {
			if errors.As(err, &notFound) {
				continue
			}
			return []assets.Audience{}, fmt.Errorf("error while reading audiences: %w", err)
		}

		aud, ok := v.(assets.Audience)
		if !ok {
			errMsg := fmt.Sprintf("error while reading audience with id: %s from db", keys[i])
			return []assets.Audience{}, NewInternalRepositoryError(errMsg, nil)

		}
		res = append(res, aud)
	}
	return res, nil
}

func (a *AudienceRepo) Delete(ctx context.Context, audienceID uint32) error {
	key := fmt.Sprintf("%d", audienceID)
	_, exists := a.audConn.Delete(key)
	if !exists {
		errMsg := fmt.Sprintf("asset audience with id: %v does not exist", audienceID)
		return app.NewEntityNotFoundError(errMsg)
	}
	return nil
}

func (a *AudienceRepo) Update(ctx context.Context, audience assets.Audience) error {
	aud, err := a.Get(ctx, audience.ID)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", aud.ID)
	a.audConn.Upsert(key, audience)
	return nil
}

func (a *AudienceRepo) Star(ctx context.Context, userEmail string, audienceID uint32) error {
	if audienceID == 0 || userEmail == "" {
		return fmt.Errorf("audience id and user email cannot be empty")
	}

	v, err := a.starConn.Get(userEmail)
	var notFound *ErrNotFound

	if err != nil && !errors.As(err, &notFound) {
		errMsg := fmt.Sprintf("error while reading starred audiences for user with email: %s:", userEmail)
		return NewInternalRepositoryError(errMsg, nil)
	}

	if starred, ok := v.([]uint32); ok {
		for i := range starred {
			if audienceID == starred[i] {
				return nil
			}
		}
		starred = append(starred, audienceID)
		a.starConn.Upsert(userEmail, starred)
	} else {
		res := make([]uint32, 0, 20)
		res = append(res, audienceID)
		a.starConn.Upsert(userEmail, res)
	}

	return nil
}

func (a *AudienceRepo) Unstar(ctx context.Context, userEmail string, audienceID uint32) error {
	if audienceID == 0 || userEmail == "" {
		return fmt.Errorf("audience id and user email cannot be empty")
	}

	v, err := a.starConn.Get(userEmail)

	var notFound *ErrNotFound
	switch {
	case err != nil && errors.As(err, &notFound):
		errMsg := fmt.Sprintf("cannot find starred audience assets for user: %s", userEmail)
		return app.NewEntityNotFoundError(errMsg)
	case err != nil && !errors.As(err, &notFound): // This Will never evaluate but hypothetically that could be an internal db error.
		return NewInternalRepositoryError(UnknownError, err)
	default:
		if starred, ok := v.([]uint32); ok {
			found := false
			for i := range starred {
				if starred[i] == audienceID {
					starred[i] = starred[len(starred)-1]
					found = true
					break
				}
			}
			if found {
				a.starConn.Upsert(userEmail, starred[:len(starred)-1])
				return nil
			}
		}
	}
	errMsg := fmt.Sprintf("cannot find starred audience assets for user: %s", userEmail)
	return app.NewEntityNotFoundError(errMsg)
}

func (a *AudienceRepo) GetStarredIDsForUser(ctx context.Context, userEmail string) ([]uint32, error) {
	if userEmail == "" {
		return []uint32{}, fmt.Errorf("user email cannot be empty")
	}

	v, err := a.starConn.Get(userEmail)
	if err != nil {
		errMsg := fmt.Sprintf("no starred audien found for user: %s", userEmail)
		return nil, app.NewEntityNotFoundError(errMsg)
	}

	starred, ok := v.([]uint32)
	if !ok {
		return []uint32{}, fmt.Errorf("error while reading starred audiences for user: %s", userEmail)
	}
	return starred, nil
}
