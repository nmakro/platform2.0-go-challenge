package assets

import "github.com/nmakro/platform2.0-go-challenge/internal/app/user"

type AssetService struct {
	userService  *user.UserService
	audienceRepo AudienceRepository
	chartRepo    ChartRepository
	insightRepo  InsightRepository
}

func NewAssetService(userService *user.UserService, audience AudienceRepository, chart ChartRepository, insight InsightRepository) *AssetService {
	return &AssetService{
		userService:  userService,
		audienceRepo: audience,
		chartRepo:    chart,
		insightRepo:  insight,
	}
}
