package assets

import "context"

func (s AssetService) AddChart(ctx context.Context, c Chart) error {
	if err := ValidateChart(c); err != nil {
		return err
	}
	return s.chartRepo.Add(ctx, c)
}

func (s AssetService) UpdateChart(ctx context.Context, chartID uint32, c UpdateChartCommand) error {
	old, err := s.GetChart(ctx, chartID)
	if err != nil {
		return err
	}

	if c.Description != nil {
		old.Description = *c.Description
	}

	if c.Title != nil {
		old.Title = *c.Title
	}

	if c.XAxis != nil {
		old.XAxis = *c.XAxis
	}

	if c.YAxis != nil {
		old.XAxis = *c.YAxis
	}

	if c.Data != nil {
		old.Data = *c.Data
	}

	if err := ValidateChart(old); err != nil {
		return err
	}
	return s.chartRepo.Update(ctx, old)
}

func (s AssetService) GetChart(ctx context.Context, chartID uint32) (Chart, error) {
	return s.chartRepo.Get(ctx, chartID)
}

func (s AssetService) GetAllCharts(ctx context.Context) ([]Chart, error) {
	return s.chartRepo.List(ctx)
}

func (s AssetService) DeleteChart(ctx context.Context, chartID uint32) error {
	return s.chartRepo.Delete(ctx, chartID)
}

func (s AssetService) StarChart(ctx context.Context, userEmail string, chartID uint32) error {
	if _, err := s.userService.GetUser(ctx, userEmail); err != nil {
		return err
	}
	return s.chartRepo.Star(ctx, userEmail, chartID)
}

func (s AssetService) GetChartsForUser(ctx context.Context, userEmail string) ([]Chart, error) {
	_, err := s.userService.GetUser(ctx, userEmail)
	if err != nil {
		return []Chart{}, err
	}

	ids, err := s.chartRepo.GetStarredIDsForUser(ctx, userEmail)

	if err != nil {
		return []Chart{}, err
	}

	charts, err := s.chartRepo.GetMany(ctx, ids)

	if err != nil {
		return []Chart{}, err
	}

	return charts, nil
}

func (s AssetService) UnstarChart(ctx context.Context, userEmail string, chartID uint32) error {
	return s.chartRepo.Unstar(ctx, userEmail, chartID)
}

type chartValidator = func(c Chart) error

func ValidateChartID(c Chart) error {
	if c.ID == 0 {
		return NewErrValidation("chart id cannot be empty")
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
