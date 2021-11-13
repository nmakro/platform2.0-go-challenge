package assets

import "context"

func (s AssetService) AddChart(ctx context.Context, c Chart) error {
	return s.ChartRepo.Add(ctx, c)
}

func (s AssetService) UpdateChart(ctx context.Context, chart Chart) error {
	return s.ChartRepo.Update(ctx, chart)
}

func (s AssetService) GetChart(ctx context.Context, chartID uint32) (Chart, error) {
	return s.ChartRepo.Get(ctx, chartID)
}

func (s AssetService) DeleteChart(ctx context.Context, chartID uint32) error {
	return s.ChartRepo.Delete(ctx, chartID)
}

func (s AssetService) StarChart(ctx context.Context, userEmail string, chartID uint32) error {
	return s.ChartRepo.Star(ctx, userEmail, chartID)
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
