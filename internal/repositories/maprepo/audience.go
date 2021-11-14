package maprepo

import (
	"context"
	"errors"
	"fmt"

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
		return NewDuplicateEntryError()
	}
	return nil
}

func (a *AudienceRepo) Get(ctx context.Context, audienceID uint32) (assets.Audience, error) {
	key := fmt.Sprintf("%d", audienceID)
	v, err := a.audConn.Get(key)
	if err != nil {
		return assets.Audience{}, err
	}
	if aud, ok := v.(assets.Audience); ok {
		return aud, nil
	}
	return assets.Audience{}, fmt.Errorf("error while reading audience from db.")
}

func (a *AudienceRepo) GetMany(ctx context.Context, audienceIDs []uint32) ([]assets.Audience, error) {
	res := make([]assets.Audience, 0, len(audienceIDs))
	for _, id := range audienceIDs {
		a, err := a.Get(ctx, id)
		if err != nil {
			var notFound *ErrEntityNotFound
			if errors.As(err, &notFound) {
				continue
			}
			return []assets.Audience{}, err
		}
		res = append(res, a)
	}
	return res, nil
}

func (a *AudienceRepo) Delete(ctx context.Context, audienceID uint32) error {
	key := fmt.Sprintf("%d", audienceID)
	_, exists := a.audConn.Delete(key)
	if !exists {
		return NewEntityNotFoundError()
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
	notFound := NewEntityNotFoundError()
	if err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("error while reading starred audiences for user with email: %s", userEmail)
	}

	if starred, ok := v.([]uint32); ok {
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
	notFound := NewEntityNotFoundError()

	switch {
	case err != nil && errors.As(err, &notFound):
		return fmt.Errorf("cannot find starred audience assets for user: %s", userEmail)
	case err != nil && !errors.As(err, &notFound): // This Will never evaluate but hypothetically that could be an internal db error.
		return err
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
	return fmt.Errorf("cannot find starred audience asset: %v for user: %s", audienceID, userEmail)
}

func (a *AudienceRepo) GetStarredIDsForUser(ctx context.Context, userEmail string) ([]uint32, error) {
	if userEmail == "" {
		return []uint32{}, fmt.Errorf("user email cannot be empty")
	}

	v, err := a.starConn.Get(userEmail)
	if err != nil {
		return nil, err
	}

	starred, ok := v.([]uint32)
	if !ok {
		return []uint32{}, fmt.Errorf("error while reading starred audiences for user: %s", userEmail)
	}
	return starred, nil
}
