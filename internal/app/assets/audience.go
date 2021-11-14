package assets

import "context"

func (s AssetService) ValidateAudience(a Audience) bool {
	return a.AgeGroup.IsValid()
}

func (s AssetService) AddAudience(ctx context.Context, a Audience) error {
	if err := ValidateAudience(a); err != nil {
		return err
	}
	return s.AudienceRepo.Add(ctx, a)
}

func (s AssetService) UpdateAudience(ctx context.Context, a Audience) error {
	if err := ValidateAudience(a); err != nil {
		return err
	}
	return s.AudienceRepo.Update(ctx, a)
}

func (s AssetService) DeleteAudience(ctx context.Context, audienceID uint32) error {
	return s.AudienceRepo.Delete(ctx, audienceID)
}

func (s AssetService) StarAudience(ctx context.Context, userEmail string, audienceID uint32) error {
	if _, err := s.userService.GetUser(ctx, userEmail); err != nil {
		return err
	}
	return s.AudienceRepo.Star(ctx, userEmail, audienceID)
}

func (s AssetService) UnstarAudience(ctx context.Context, userEmail string, audienceID uint32) error {
	return s.AudienceRepo.Unstar(ctx, userEmail, audienceID)
}

func (s AssetService) GetAudiencesForUser(ctx context.Context, userEmail string) ([]Audience, error) {
	_, err := s.userService.GetUser(ctx, userEmail)
	if err != nil {
		return []Audience{}, err
	}

	ids, err := s.AudienceRepo.GetStarredIDsForUser(ctx, userEmail)

	if err != nil {
		return []Audience{}, err
	}

	audiences, err := s.AudienceRepo.GetMany(ctx, ids)

	if err != nil {
		return []Audience{}, err
	}

	return audiences, nil
}

type audienceValidator = func(a Audience) error

func ValidateAudienceID(a Audience) error {
	if a.ID == 0 {
		return NewAssetNoIDError()
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
