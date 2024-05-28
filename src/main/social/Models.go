package social

// MediaEntity contract
type MediaEntity struct {
	ID             int64
	SourceStatusID int64

	IDStr             string
	MediaURL          string
	MediaURLHttps     string
	SourceStatusIDStr string
	Type              string

	Info MediaInfo
}

// MediaInfo contract.
type MediaInfo struct {
	AspectRatio    [2]int
	DurationMillis int

	MediaVariant []MediaVariant
}

// MediaVariant contract
type MediaVariant struct {
	ContentType string
	URL         string

	Bitrate int
}

// ExtendedEntity contract
type ExtendedEntity struct {
	Media []MediaEntity
}

// Tweet holds tweet data
type Tweet struct {
	ID                int64
	InReplyToStatusID int64
	InReplyToUserID   int64
	QuotedStatusID    int64

	FavoriteCount int
	RetweetCount  int

	Favorited         bool
	PossiblySensitive bool
	Retweeted         bool
	Truncated         bool
	WithheldCopyright bool
	HasVideo          bool

	FilterLevel          string
	IDStr                string
	InReplyToScreenName  string
	InReplyToStatusIDStr string
	InReplyToUserIDStr   string
	Lang                 string
	Source               string
	Text                 string
	WithheldScope        string
	QuotedStatusIDStr    string
	ImageURL             string
	VideoURL             string
	DateCreated          string

	WithheldInCountries []string

	User *User

	Hashtags Hashtags
	URLS     TweetURLs
}

// Tweets holds a list of tweets
type Tweets []Tweet

// Follower contract
type Follower struct {
	CreatedAt                      string
	Description                    string
	Email                          string
	IDStr                          string
	Lang                           string
	Location                       string
	Name                           string
	ProfileBackgroundColor         string
	ProfileBackgroundImageURL      string
	ProfileBackgroundImageURLHttps string
	ProfileBannerURL               string
	ProfileImageURL                string
	ProfileImageURLHttps           string
	ProfileLinkColor               string
	ProfileSidebarBorderColor      string
	ProfileSidebarFillColor        string
	ProfileTextColor               string
	ScreenName                     string
	Timezone                       string
	URL                            string
	WithholdScope                  string

	WithheldInCountries []string

	ContributorsEnabled       bool
	DefaultProfile            bool
	DefaultProfileImage       bool
	FollowRequestSent         bool
	Following                 bool
	GeoEnabled                bool
	IsTranslator              bool
	Notifications             bool
	ProfileBackgroundTile     bool
	ProfileUseBackgroundImage bool
	Protected                 bool
	ShowAllInlineMedia        bool
	Verified                  bool

	FavouritesCount int
	FollowersCount  int
	FriendsCount    int
	ListedCount     int
	StatusesCount   int
	UtcOffset       int

	ID int64

	Status *Tweet
}

// Followers contract
type Followers []Follower

// User contract
type User struct {
	// TODO: group them
	ContributorsEnabled            bool
	CreatedAt                      string
	DefaultProfile                 bool
	DefaultProfileImage            bool
	Description                    string
	Email                          string
	FavouritesCount                int
	FollowRequestSent              bool
	Following                      bool
	FollowersCount                 int
	FriendsCount                   int
	GeoEnabled                     bool
	ID                             int64
	IDStr                          string
	IsTranslator                   bool
	Lang                           string
	ListedCount                    int
	Location                       string
	Name                           string
	Notifications                  bool
	ProfileBackgroundColor         string
	ProfileBackgroundImageURL      string
	ProfileBackgroundImageURLHttps string
	ProfileBackgroundTile          bool
	ProfileBannerURL               string
	ProfileImageURL                string
	ProfileImageURLHttps           string
	ProfileLinkColor               string
	ProfileSidebarBorderColor      string
	ProfileSidebarFillColor        string
	ProfileTextColor               string
	ProfileUseBackgroundImage      bool
	Protected                      bool
	ScreenName                     string
	ShowAllInlineMedia             bool
	StatusesCount                  int
	Timezone                       string
	URL                            string
	UtcOffset                      int
	Verified                       bool
	WithheldInCountries            []string
	WithholdScope                  string
}

// Users contract
type Users []User

// Hashtag contract
type Hashtag struct {
	Text       string
	URL        string
	DisplayURL string

	StartIndex int
	EndIndex   int
}

// Hashtags array
type Hashtags []Hashtag

// HashtagURL contract
type HashtagURL struct {
	DisplayURL string
	URL        string
}

// TweetURL contract
type TweetURL struct {
	DisplayURL  string
	ExpandedURL string
	URL         string

	StartIndex int
	EndIndex   int
}

// TweetURLs array
type TweetURLs []TweetURL

type TwitterConfiguration struct {
	Url               string
	Version           string
	TokenUrl          string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	GrantType         string
	Handle            string

	HandleId int
}

type TwitterConfig struct {
	Twitter TwitterConfiguration
}
