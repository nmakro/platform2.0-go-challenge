package main

import (
	"github.com/nmakro/platform2.0-go-challenge/internal/app/assets"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/session"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/user"
	"github.com/nmakro/platform2.0-go-challenge/internal/repositories/maprepo"
	repo "github.com/nmakro/platform2.0-go-challenge/internal/repositories/maprepo"
)

type AppImplemention struct {
	assetService   *assets.AssetService
	userService    *user.UserService
	sessionService *session.SessioService
}

func NewApp() *AppImplemention {
	audienceRepoClient := repo.NewClient()
	audienceStarredClient := repo.NewClient()
	audienceRepo := maprepo.NewAudienceRepo(audienceRepoClient, audienceStarredClient)

	chartRepoClient := repo.NewClient()
	chartStarredClient := repo.NewClient()
	chartRepo := maprepo.NewChartRepo(chartRepoClient, chartStarredClient)

	insightRepoClient := repo.NewClient()
	insightStarredClient := repo.NewClient()
	insightsRepo := maprepo.NewAInsightRepo(insightRepoClient, insightStarredClient)

	userRepoClient := repo.NewClient()
	userRepo := maprepo.NewUserRepo(userRepoClient)
	userService := user.NewService(userRepo)

	assetService := assets.NewAssetService(userService, audienceRepo, chartRepo, insightsRepo)

	sessionRepo := maprepo.NewSessionRepo()
	sessionSerive := session.NewSessionService(sessionRepo)

	return &AppImplemention{
		assetService:   assetService,
		userService:    userService,
		sessionService: sessionSerive,
	}
}
