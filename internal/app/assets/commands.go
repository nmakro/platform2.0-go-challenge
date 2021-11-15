package assets

import (
	"sync/atomic"
)

type AddAudienceCommand struct {
	Description      string `json:"description"`
	BirthCountry     string `json:"birth_country"`
	Gender           string `json:"gender" validate:"oneof=male female"`
	AgeGroupFrom     uint32 `json:"age_group_from"`
	AgeGroupTo       uint32 `json:"age_group_to"`
	NumOfPurchases   uint32 `json:"purchases"`
	SocialMediaHours uint32 `json:"social_hours"`
}

func (a *AddAudienceCommand) BuildFromCmd() Audience {
	ageGroup := NewAgeGroup(a.AgeGroupFrom, a.AgeGroupTo)

	audience := Audience{
		ID:               atomic.AddUint32(&AudienceIndex, 1),
		SocialMediaHours: a.SocialMediaHours,
		NumOfPurchases:   a.NumOfPurchases,
		AgeGroup:         ageGroup,
		Gender:           GenderFromString(a.Gender),
		BirthCountry:     a.BirthCountry,
		Description:      a.Description,
	}

	return audience
}

type UpdateAudienceCommand struct {
	Description      *string `json:"description"`
	BirthCountry     *string `json:"birth_country"`
	Gender           *string `json:"gender" validate:"oneof=male female"`
	AgeGroupFrom     *uint32 `json:"age_group_from"`
	AgeGroupTo       *uint32 `json:"age_group_to"`
	NumOfPurchases   *uint32 `json:"num_of_purchases"`
	SocialMediaHours *uint32 `json:"social_media_hours"`
}

var (
	AudienceIndex uint32
	InsightIndex  uint32
	ChartIndex    uint32
)
