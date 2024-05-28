package social

import (
	"fmt"
	"strings"

	"github.com/dghubble/go-twitter/twitter"

	worldHealthOrg "main/worldHealthOrg"
)

func getVideoURL(t twitter.Tweet) string {
	if !hasMedia(t) {
		return ""
	}

	if !hasVideoMedia(t) {
		return ""
	}

	if len(t.ExtendedEntities.Media[0].VideoInfo.Variants) == 0 {
		return ""
	}

	return t.ExtendedEntities.Media[0].VideoInfo.Variants[0].URL
}

func getImageURL(t twitter.Tweet) string {
	if hasMedia(t) {
		// todo: convert this to an array
		return t.Entities.Media[0].MediaURL
	}

	return ""
}

func getHashtags(t twitter.Tweet) Hashtags {
	var hashtags Hashtags

	if len(t.Entities.Hashtags) > 0 {
		for _, ht := range t.Entities.Hashtags {
			hashtagURL := getHashtagURL(ht.Text)

			hastagMap := &Hashtag{
				Text:       ht.Text,
				StartIndex: ht.Indices.Start(),
				EndIndex:   ht.Indices.End(),

				DisplayURL: hashtagURL.DisplayURL,
				URL:        hashtagURL.URL,
			}

			hashtags = append(hashtags, *hastagMap)
		}
	}

	return hashtags
}

func getHashtagURL(key string) HashtagURL {
	hashtagURL := HashtagURL{
		DisplayURL: "#" + key,
		URL:        "https://twitter.com/hashtag/" + key + "?src=hash",
	}

	return hashtagURL
}

func getURLs(t twitter.Tweet) TweetURLs {
	var urls TweetURLs

	if t.Entities != nil && len(t.Entities.Urls) > 0 {
		for _, tu := range t.Entities.Urls {
			urlMap := &TweetURL{
				DisplayURL:  tu.DisplayURL,
				URL:         tu.URL,
				ExpandedURL: tu.ExpandedURL,

				StartIndex: tu.Indices.Start(),
				EndIndex:   tu.Indices.End(),
			}

			urls = append(urls, *urlMap)
		}
	}

	if t.ExtendedEntities != nil && len(t.ExtendedEntities.Media) > 0 {
		for _, xeu := range t.ExtendedEntities.Media {
			xURLMap := &TweetURL{
				DisplayURL:  xeu.DisplayURL,
				URL:         xeu.URL,
				ExpandedURL: xeu.ExpandedURL,

				StartIndex: xeu.Indices.Start(),
				EndIndex:   xeu.Indices.End(),
			}

			urls = append(urls, *xURLMap)
		}
	}

	if t.Text != "" {
		words := strings.Fields(t.Text)

		for _, word := range words {
			if strings.Contains(word, "@") {
				screenName := ""

				woAtTag := strings.Replace(word, "@", "", -1)

				if strings.Contains(woAtTag, ":") {
					screenName = strings.Replace(woAtTag, ":", "", -1)

				} else {
					screenName = woAtTag
				}

				tagURL := &TweetURL{
					DisplayURL: "@" + screenName,
					URL:        "https://twitter.com/" + screenName,
				}

				urls = append(urls, *tagURL)
			}
		}
	}

	return urls
}

func hasMedia(t twitter.Tweet) bool {
	if t.ExtendedEntities == nil || t.ExtendedEntities.Media == nil {
		return false
	}

	length := len(t.ExtendedEntities.Media)

	if length > 0 {
		return true
	}

	return false
}

func hasVideoMedia(tweet twitter.Tweet) bool {
	if tweet.ExtendedEntities == nil {
		return false
	}

	if len(tweet.ExtendedEntities.Media) == 0 {
		return false
	}

	if strings.ToUpper(tweet.ExtendedEntities.Media[0].Type) == "VIDEO" {
		return true
	}

	return false
}

func handleTimelineResponse(ts []twitter.Tweet) Tweets {
	var tweets Tweets

	for _, tweet := range ts {
		//map just the data we want.  keep it simple

		userMap := &User{
			ScreenName:                tweet.User.ScreenName,
			ProfileBackgroundImageURL: tweet.User.ProfileBackgroundImageURL,
			Verified:                  tweet.User.Verified,
			ProfileBackgroundColor:    tweet.User.ProfileBackgroundColor,
			Location:                  tweet.User.Location,
			Description:               tweet.User.Description,
			FollowersCount:            tweet.User.FollowersCount,
			ID:                        tweet.User.ID,
			Name:                      tweet.User.Name,
			ProfileImageURL:           tweet.User.ProfileImageURL,
		}

		tweetMap := &Tweet{
			ID:            tweet.ID,
			Text:          tweet.Text,
			FavoriteCount: tweet.FavoriteCount,
			RetweetCount:  tweet.RetweetCount,
			DateCreated:   tweet.CreatedAt,
			Favorited:     tweet.Favorited,
			Retweeted:     tweet.Retweeted,

			VideoURL: getVideoURL(tweet),
			ImageURL: getImageURL(tweet),

			User: userMap,

			Hashtags: getHashtags(tweet),
			URLS:     getURLs(tweet),

			HasVideo: hasVideoMedia(tweet),
		}

		// add to the list
		tweets = append(tweets, *tweetMap)
	}

	return tweets
}

// FetchHomeTimeline returns the home timeline tweets
func FetchHomeTimeline() (Tweets, error) {
	if twitterClient, e := TwitterInstance(); e == nil {
		// TODO: replace hardcoded Count value
		if ts, _, err := twitterClient.Timelines.HomeTimeline(&twitter.HomeTimelineParams{Count: 50}); err == nil {
			return handleTimelineResponse(ts), nil

		} else {
			return nil, err
		}

	} else {
		return nil, e
	}
}

// FetchUserTimeline returns the user timeline tweets
func FetchUserTimeline() (Tweets, error) {
	if twitterClient, e := TwitterInstance(); e == nil {
		// TODO: replace hardcoded count
		if ts, _, err := twitterClient.Timelines.UserTimeline(&twitter.UserTimelineParams{Count: 50}); err == nil {
			return handleTimelineResponse(ts), nil

		} else {
			return nil, err
		}

	} else {
		return nil, e
	}
}

// RunTasks will start all twitter tasks
func RunTasks() {
	fmt.Println("Task is being performed")

	if infantNutrition, e := worldHealthOrg.FetchInfantNutrition("USA"); e == nil {
		fmt.Println(infantNutrition)
	}
}

// GetFollowers returns a list of Followers
// func GetFollowers() (Followers, error) {
// 	if twitterClient, e := twitterRepository.TwitterInstance(); e == nil {

// 		if fl, _, err := twitterClient.Followers.List(&twitter.FollowerListParams{}); err == nil {
// 			var followers Followers

// 			for _, fi := range fl {
// 				followers = append(followers, Follower{})
// 			}
// 		}
// 	}
// }
