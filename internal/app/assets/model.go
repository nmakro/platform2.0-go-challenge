package assets

//Assets value structs

// AgeGroup represents an age group.
// A valid AgeGroup must have StartAge less than EndYear and StartYear > 18 and EndYear < 100.
type AgeGroup struct {
	StartYear uint32
	EndYear   uint32
}

func (a AgeGroup) IsValid() bool {
	return a.StartYear < a.EndYear && a.StartYear > 18 && a.EndYear < 100
}

func NewAgeGroup(from, to uint32) AgeGroup {
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
	XValue interface{}
	YValue interface{}
	Data   []byte
}

// Assets entity structs.

// Audience struct represents an Audience entity.
type Audience struct {
	ID               uint32
	SocialMediaHours uint32
	NumOfPurchases   uint32
	AgeGroup         AgeGroup
	Gender           Gender
	BirthCountry     string
	Description      string
}

// Insight struct represents an Insight entity.
type Insight struct {
	ID          uint32
	Topic       string
	Text        string
	Description string
}

// Chart struct represents a Chart entity.
type Chart struct {
	ID          uint32
	Title       string
	XAxis       string
	YAxis       string
	Description string
	Data        []DataPoint
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
