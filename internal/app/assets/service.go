package assets

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
