package assets

import "context"

func (s AssetService) AddInsight(ctx context.Context, i Insight) error {
	if err := ValidateInsight(i); err != nil {
		return err
	}
	return s.insightRepo.Add(ctx, i)
}

func (s AssetService) GetInsight(ctx context.Context, insightID uint32) (Insight, error) {
	return s.insightRepo.Get(ctx, insightID)
}

func (s AssetService) DeleteInsight(ctx context.Context, insightID uint32) error {
	return s.insightRepo.Delete(ctx, insightID)
}

func (s AssetService) UpdateInsight(ctx context.Context, insight Insight) error {
	if err := ValidateInsight(insight); err != nil {
		return err
	}
	return s.insightRepo.Update(ctx, insight)
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
