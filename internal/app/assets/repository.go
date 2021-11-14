package assets

import "context"

type AudienceRepository interface {
	Add(ctx context.Context, a Audience) error
	Update(ctx context.Context, a Audience) error
	Get(ctx context.Context, AudienceID uint32) (Audience, error)
	GetMany(ctx context.Context, AudienceIDs []uint32) ([]Audience, error)
	Delete(ctx context.Context, AudienceID uint32) error
	Star(ctx context.Context, userEmail string, audienceID uint32) error
	Unstar(ctx context.Context, userEmail string, audienceID uint32) error
	GetStaredIDsForUser(ctx context.Context, userEmail string) ([]uint32, error)
}

type ChartRepository interface {
	Add(ctx context.Context, c Chart) error
	Update(ctx context.Context, c Chart) error
	Get(ctx context.Context, ChartID uint32) (Chart, error)
	GetMany(ctx context.Context, ChartIDs []uint32) ([]Chart, error)
	Delete(ctx context.Context, ChartID uint32) error
	Star(ctx context.Context, UserEmail string, ChartID uint32) error
	Unstar(ctx context.Context, UserEmail string, ChartID uint32) error
	GetStaredIDsForUser(ctx context.Context, userEmail string) ([]uint32, error)
}

type InsightRepository interface {
	Add(ctx context.Context, i Insight) error
	Update(ctx context.Context, i Insight) error
	Get(ctx context.Context, insightID uint32) (Insight, error)
	GetMany(ctx context.Context, insightIDs []uint32) ([]Insight, error)
	Delete(ctx context.Context, insightID uint32) error
	Star(ctx context.Context, userEmail string, insightID uint32) error
	Unstar(ctx context.Context, userEmail string, insightID uint32) error
	GetStaredIDsForUser(ctx context.Context, userEmail string) ([]uint32, error)
}
