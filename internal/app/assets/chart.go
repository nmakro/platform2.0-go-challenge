package assets

import "context"

func (s AssetService) AddChart(ctx context.Context, c Chart) error {
	if err := ValidateChart(c); err != nil {
		return err
	}
	return s.ChartRepo.Add(ctx, c)
}

func (s AssetService) UpdateChart(ctx context.Context, c Chart) error {
	if err := ValidateChart(c); err != nil {
		return err
	}
	return s.ChartRepo.Update(ctx, c)
}

func (s AssetService) GetChart(ctx context.Context, chartID uint32) (Chart, error) {
	return s.ChartRepo.Get(ctx, chartID)
}

func (s AssetService) DeleteChart(ctx context.Context, chartID uint32) error {
	return s.ChartRepo.Delete(ctx, chartID)
}

func (s AssetService) StarChart(ctx context.Context, userEmail string, chartID uint32) error {
	if _, err := s.userService.GetUser(ctx, userEmail); err != nil {
		return err
	}
	return s.ChartRepo.Star(ctx, userEmail, chartID)
}

func (s AssetService) GetChartsForUser(ctx context.Context, userEmail string) ([]Chart, error) {
	_, err := s.userService.GetUser(ctx, userEmail)
	if err != nil {
		return []Chart{}, err
	}

	ids, err := s.ChartRepo.GetStarredIDsForUser(ctx, userEmail)

	if err != nil {
		return []Chart{}, err
	}

	charts, err := s.ChartRepo.GetMany(ctx, ids)

	if err != nil {
		return []Chart{}, err
	}

	return charts, nil
}

func (s AssetService) UnstarChart(ctx context.Context, userEmail string, chartID uint32) error {
	return s.ChartRepo.Unstar(ctx, userEmail, chartID)
}

type chartValidator = func(c Chart) error

func ValidateChartID(c Chart) error {
	if c.ID == 0 {
		return NewAssetNoIDError()
	}
	return nil
}

var chartValidators = []chartValidator{
	ValidateChartID,
}

func ValidateChart(c Chart) error {
	for _, f := range chartValidators {
		if err := f(c); err != nil {
			return err
		}
	}
	return nil
}
