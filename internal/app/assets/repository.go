package assets

import "context"

type AudienceRepository interface {
	Add(ctx context.Context, a Audience) error
	Update(ctx context.Context, AudienceID uint32, desc string) error
	Delete(ctx context.Context, AudienceID uint32) error
	Star(ctx context.Context, AudienceID, UserID uint32) error
	Unstar(ctx context.Context, AudienceID, UserID uint32) error
}

type ChartRepository interface {
	Add(ctx context.Context, c Chart) error
	Update(ctx context.Context, ChartID uint32, desc string) error
	Delete(ctx context.Context, ChartID uint32) error
	Star(ctx context.Context, ChartID, UserID uint32) error
	Unstar(ctx context.Context, ChartID, UserID uint32) error
}

type InsightRepository interface {
	Add(ctx context.Context, i Insight) error
	Update(ctx context.Context, InsightID uint32, desc string) error
	Delete(ctx context.Context, InsightID uint32) error
	Star(ctx context.Context, InsightID, UserID uint32) error
	Unstar(ctx context.Context, InsightID, UserID uint32) error
}
