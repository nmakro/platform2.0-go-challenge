package assets

import "fmt"

//Assets value structs

// AgeGroup represents an age group.
// A valid AgeGroup must have StartAge less than EndYear and StartYear > 18 and EndYear < 100.
type AgeGroup struct {
	StartYear uint32 `json:"start_year"`
	EndYear   uint32 `json:"end_year"`
}

func (a AgeGroup) IsValid() bool {
	return a.StartYear < a.EndYear && a.StartYear > 18 && a.EndYear < 100
}

func NewAgeGroup(from, to uint32) AgeGroup {
	fmt.Println("in new")
	return AgeGroup{
		StartYear: from,
		EndYear:   to}
}

// Gender struct represents Gender.
// Valid values are Male/Female/All.
type Gender struct {
	genderType string
}

func (g Gender) String() string {
	return g.genderType
}

func GenderFromString(g string) Gender {
	switch g {
	case Female.genderType:
		return Female
	case Male.genderType:
		return Male
	case All.genderType:
		return All
	default:
		return Unknown
	}
}

var (
	Male    = Gender{"Male"}
	Female  = Gender{"Female"}
	All     = Gender{"All"}
	Unknown = Gender{""}
)

// DataPoint represents a Chart Data value.
// An instance of a DataPoint has a XAxis e.x year/month etc and a YAxis e.x. HoursSpent or NumOfPurchases.
type DataPoint struct {
	XValue interface{} `json:"xaxis_value"`
	YValue interface{} `json:"yaxis_value"`
	Data   []byte      `json:"data"`
}

// Assets entity structs.

// Audience struct represents an Audience entity.
type Audience struct {
	ID               uint32   `json:"id"`
	SocialMediaHours uint32   `json:"social_media_hours"`
	NumOfPurchases   uint32   `json:"num_of_purchases"`
	AgeGroup         AgeGroup `json:"age_group"`
	Gender           Gender   `json:"gender"`
	BirthCountry     string   `json:"birth_country"`
	Description      string   `json:"description"`
}

// Insight struct represents an Insight entity.
type Insight struct {
	ID          uint32 `json:"id"`
	Topic       string `json:"topic"`
	Text        string `json:"text"`
	Description string `json:"description"`
}

// Chart struct represents a Chart entity.
type Chart struct {
	ID          uint32      `json:"id"`
	Title       string      `json:"title"`
	XAxis       string      `json:"xaxis"`
	YAxis       string      `json:"yaxis"`
	Description string      `json:"description"`
	Data        []DataPoint `json:"data_point"`
}

type Asset struct {
	assetType string
}

var (
	AudienceAsset = Asset{"Audience"}
	InsightAsset  = Asset{"Insight"}
	ChartAsset    = Asset{"Chart"}
	UnknownAsset  = Asset{""}
)

func (a Asset) String() string {
	return a.assetType
}

func FromAssetString(a Asset) Asset {
	switch a.assetType {
	case AudienceAsset.assetType:
		return AudienceAsset
	case InsightAsset.assetType:
		return InsightAsset
	case ChartAsset.assetType:
		return ChartAsset
	default:
		return UnknownAsset
	}
}
