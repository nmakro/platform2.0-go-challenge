package assets

import (
	"context"
)

func (s AssetService) AddInsight(ctx context.Context, i Insight) error {
	if err := ValidateInsight(i); err != nil {
		return err
	}
	return s.insightRepo.Add(ctx, i)
}

func (s AssetService) GetInsight(ctx context.Context, insightID uint32) (Insight, error) {
	return s.insightRepo.Get(ctx, insightID)
}

func (s AssetService) ListInsightAssets(ctx context.Context) ([]Insight, error) {
	return s.insightRepo.List(ctx)
}

func (s AssetService) DeleteInsight(ctx context.Context, insightID uint32) error {
	return s.insightRepo.Delete(ctx, insightID)
}

func (s AssetService) UpdateInsight(ctx context.Context, insightID uint32, insight UpdateInsightCommand) error {
	old, err := s.GetInsight(ctx, insightID)
	if err != nil {
		return err
	}

	if insight.Description != nil {
		old.Description = *insight.Description
	}

	if insight.Text != nil {
		old.Text = *insight.Text
	}

	if insight.Topic != nil {
		old.Topic = *insight.Topic
	}

	if err := ValidateInsight(old); err != nil {
		return err
	}
	return s.insightRepo.Update(ctx, old)
}

func (s AssetService) StartInsight(ctx context.Context, userEmail string, insightID uint32) error {
	_, err := s.userService.GetUser(ctx, userEmail)
	if err != nil {
		return err
	}
	return s.insightRepo.Star(ctx, userEmail, insightID)
}

func (s AssetService) UnstarInsight(ctx context.Context, userEmail string, insightID uint32) error {
	return s.insightRepo.Unstar(ctx, userEmail, insightID)
}

func (s AssetService) GetInsightsForUser(ctx context.Context, userEmail string) ([]Insight, error) {
	_, err := s.userService.GetUser(ctx, userEmail)
	if err != nil {
		return []Insight{}, err
	}

	ids, err := s.insightRepo.GetStarredIDsForUser(ctx, userEmail)

	if err != nil {
		return []Insight{}, err
	}

	insights, err := s.insightRepo.GetMany(ctx, ids)

	if err != nil {
		return []Insight{}, err
	}

	return insights, nil
}

type insightValidator = func(i Insight) error

func ValidateInsightID(i Insight) error {
	if i.ID == 0 {
		return NewErrValidation("insight id cannot be empty")
	}
	return nil
}

var insightValidators = []insightValidator{
	ValidateInsightID,
}

func ValidateInsight(i Insight) error {
	for _, f := range insightValidators {
		if err := f(i); err != nil {
			return err
		}
	}
	return nil
}
