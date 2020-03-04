package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pisign/pisign-backend/types"
)

type TwitterResponse struct {
	Tweets []twitter.Tweet
}

func (o *TwitterResponse) Update(arguments interface{}) error {
	a := (arguments).(types.TwitterConfig)
	consumer_key := a.ConsumerKey
	consumer_secret := a.ConsumerSecret
	access_token := a.AccessToken
	access_secret := a.AccessSecret
	tweet_count := a.TweetCount
	user_handle := a.UserHandle

	// Twitter client
	config := oauth1.NewConfig(consumer_key, consumer_secret)
	token := oauth1.NewToken(access_token, access_secret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// make request
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: user_handle,
		Count:      tweet_count})
	if err != nil {
		return err
	}

	// store tweets in the twitter response
	o.Tweets = tweets

	return nil
}

func (o *TwitterResponse) Transform() interface{} {
	twitter_response := types.TwitterResponse{
		Tweets: o.Tweets,
	}
	return &twitter_response
}
