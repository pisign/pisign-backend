package types

import (
	"github.com/dghubble/go-twitter/twitter"
)

// TwitterResponse main format for data coming out of twitter api
type TwitterResponse struct {
	BaseMessage
	Tweets []twitter.Tweet
}

// TwitterConfig configuration arguments for twitter api
type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	UserHandle     string
	TweetCount     int
}
