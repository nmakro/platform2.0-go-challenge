package assets

import "context"

func (s AssetService) ValidateAudience(a Audience) bool {
	return a.AgeGroup.IsValid()
}

func (s AssetService) AddAudience(ctx context.Context, a Audience) error {
	if err := ValidateAudience(a); err != nil {
		return err
	}
	return s.audienceRepo.Add(ctx, a)
}

func (s AssetService) GetAudience(ctx context.Context, audienceID uint32) (Audience, error) {
	return s.audienceRepo.Get(ctx, audienceID)
}

func (s AssetService) UpdateAudience(ctx context.Context, a Audience) error {
	old, err := s.GetAudience(ctx, a.ID)
	if err != nil {
		return err
	}

	if a.Description != "" {
		old.Description = a.Description
	}

	if a.BirthCountry != "" {
		old.BirthCountry = a.BirthCountry
	}

	if a.Gender != Unknown {
		old.Gender = a.Gender
	}

	if a.NumOfPurchases != 0 {
		old.NumOfPurchases = a.NumOfPurchases
	}

	if a.SocialMediaHours != 0 {
		old.SocialMediaHours = a.SocialMediaHours
	}

	if a.AgeGroup.EndYear != 0 && a.AgeGroup.StartYear != 0 {
		old.AgeGroup = a.AgeGroup
	}

	if err := ValidateAudience(old); err != nil {
		return err
	}
	return s.audienceRepo.Update(ctx, old)
}

func (s AssetService) DeleteAudience(ctx context.Context, audienceID uint32) error {
	return s.audienceRepo.Delete(ctx, audienceID)
}

func (s AssetService) StarAudience(ctx context.Context, userEmail string, audienceID uint32) error {
	if _, err := s.userService.GetUser(ctx, userEmail); err != nil {
		return err
	}
	return s.audienceRepo.Star(ctx, userEmail, audienceID)
}

func (s AssetService) UnstarAudience(ctx context.Context, userEmail string, audienceID uint32) error {
	return s.audienceRepo.Unstar(ctx, userEmail, audienceID)
}

func (s AssetService) GetAudiencesForUser(ctx context.Context, userEmail string) ([]Audience, error) {
	_, err := s.userService.GetUser(ctx, userEmail)
	if err != nil {
		return []Audience{}, err
	}

	ids, err := s.audienceRepo.GetStarredIDsForUser(ctx, userEmail)

	if err != nil {
		return []Audience{}, err
	}

	audiences, err := s.audienceRepo.GetMany(ctx, ids)

	if err != nil {
		return []Audience{}, err
	}

	return audiences, nil
}

type audienceValidator = func(a Audience) error

func ValidateAudienceID(a Audience) error {
	if a.ID == 0 {
		msg := "audience id cannot be empty"
		return NewErrValidation(msg)
	}
	return nil
}

var audienceValidators = []audienceValidator{
	ValidateAudienceID,
}

func ValidateAudience(a Audience) error {
	for _, f := range audienceValidators {
		if err := f(a); err != nil {
			return err
		}
	}
	return nil
}
