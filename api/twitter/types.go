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
	apikey := a.APIKey
	tweet_count := a.TweetCount
	user_handle := a.UserHandle

	// Twitter client
	config := oauth1.NewConfig("consumerKey", "consumerSecret")
	token := oauth1.NewToken("accessToken", "accessSecret")
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// make request
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: user_handle,
		Count:      tweet_count})
	if err != nil {
		return err
	}

	o.Tweets = tweets

	return nil
}

func (o *TwitterResponse) Transform() interface{} {
	return nil
}
