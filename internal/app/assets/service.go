package assets

import (
	"context"
)

// type AudienceService interface {
// 	Add(ctx context.Context, a Audience) error
// 	Update(ctx context.Context, audienceID uint32) error
// 	Delete(ctx context.Context, audienceID uint32) error
// 	Star(ctx context.Context, audienceID, userID uint32) error
// 	Unstar(ctx context.Context, audienceID, userID uint32) error
// 	ListForUser(ctx context.Context, userID uint32, pg pagination.PageInfo) ([]Audience, error)
// 	Validate(a Audience) bool
// }

// type ChartService interface {
// 	Add(ctx context.Context, c Chart) error
// 	Update(ctx context.Context, chartID uint32) error
// 	Delete(ctx context.Context, chartID uint32) error
// 	Star(ctx context.Context, chartID, userID uint32) error
// 	Unstar(ctx context.Context, chartID, userID uint32) error
// 	ListForUser(ctx context.Context, userID uint32, pg pagination.PageInfo) ([]Chart, error)
// }

// type InsightService interface {
// 	Add(ctx context.Context, i InsightRepository) error
// 	Update(ctx context.Context, insightID uint32) error
// 	Delete(ctx context.Context, insightID uint32) error
// 	Star(ctx context.Context, insightID, userID uint32) error
// 	Unstar(ctx context.Context, insightID, userID uint32) error
// 	ListForUser(ctx context.Context, userID uint32, pg pagination.PageInfo) ([]Insight, error)
// }

type AssetService struct {
	AudienceRepo AudienceRepository
	ChartRepo    ChartRepository
	InsightRepo  InsightRepository
}

func NewAssetService(audience AudienceRepository, chart ChartRepository, insight InsightRepository) *AssetService {
	return &AssetService{
		AudienceRepo: audience,
		ChartRepo:    chart,
		InsightRepo:  insight,
	}
}

func (s AssetService) ValidateAudience(a Audience) bool {
	return a.AgeGroup.IsValid()
}

func (s AssetService) AddAudience(ctx context.Context, a Audience) error {
	return s.AudienceRepo.AddAudience(ctx, a)
}

func (s AssetService) UpdateAudience(ctx context.Context, audienceID uint32, desc string) error {
	return s.AudienceRepo.Update(ctx, audienceID, desc)
}

func (s AssetService) DeleteAudience(ctx context.Context, audienceID uint32) error {
	return s.AudienceRepo.Delete(ctx, audienceID)
}

func (s AssetService) StarAudience(ctx context.Context, userEmail string, audienceID uint32) error {
	return s.AudienceRepo.Star(ctx, userEmail, audienceID)
}

func (s AssetService) UnstarAudience(ctx context.Context, userEmail string, audienceID uint32) error {
	return s.AudienceRepo.Unstar(ctx, userEmail, audienceID)
}

func (s AssetService) AddChart(ctx context.Context, c Chart) error {
	return s.ChartRepo.Add(ctx, c)
}

func (s AssetService) UpdateChart(ctx context.Context, chartID uint32, desc string) error {
	return s.ChartRepo.Update(ctx, chartID, desc)
}

func (s AssetService) DeleteChart(ctx context.Context, chartID uint32) error {
	return s.ChartRepo.Delete(ctx, chartID)
}

func (s AssetService) StarChart(ctx context.Context, chartID, userID uint32) error {
	return s.ChartRepo.Star(ctx, chartID, userID)
}

func (s AssetService) UnstarChart(ctx context.Context, chartID, userID uint32) error {
	return s.ChartRepo.Unstar(ctx, chartID, userID)
}

func (s AssetService) AddInsight(ctx context.Context, i Insight) error {
	return s.InsightRepo.Add(ctx, i)
}

func (s AssetService) DeleteInsight(ctx context.Context, insightID uint32) error {
	return s.InsightRepo.Delete(ctx, insightID)
}

func (s AssetService) UpdateInsight(ctx context.Context, insightID uint32, desc string) error {
	return s.InsightRepo.Update(ctx, insightID, desc)
}

func (s AssetService) StartInsight(ctx context.Context, insightID uint32, userID uint32) error {
	return s.InsightRepo.Star(ctx, insightID, userID)
}

func (s AssetService) UnstarInsight(ctx context.Context, insightID uint32, userID uint32) error {
	return s.InsightRepo.Unstar(ctx, insightID, userID)
}
