package data

type TweetTracing struct {
	PostedBy         string
	Id               int64
	Text             string
	Retweet          bool
	OriginalPostedBy string
}

var TweetsTracing = map[int64]TweetTracing{
	1: {
		PostedBy:         "grafbaza5",
		Id:               1,
		Text:             "Neki tekst",
		Retweet:          false,
		OriginalPostedBy: "grafbaza5",
	},
	2: {
		PostedBy:         "grafbaza4",
		Id:               2,
		Text:             "Novi tweet",
		Retweet:          false,
		OriginalPostedBy: "grafbaza4",
	},
	3: {
		PostedBy:         "grafbaza5",
		Id:               3,
		Text:             "Moj 2. tweet",
		Retweet:          false,
		OriginalPostedBy: "grafbaza5",
	},
}
