package main

import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

)

var nc *nats.Conn
var tweetText = ""

func main() {
	uri := os.Getenv("NATS_URI")

	var err error

	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to NATS server " + uri)

	nc.QueueSubscribe("TwitterByText", "TwitterTeller", publishTweetFromStream)
	select {} // Block forever

}

func readKeys() []string {
	myKeysFile, err := ioutil.ReadFile("my-keys")
	if err != nil {
		fmt.Println("There was a problem with the twitter api keys")
	}
	return strings.Split(string(myKeysFile), "\n")
}

func auth() *anaconda.TwitterApi {
	myKeys := readKeys()
	consumerKey, consumerSecret := myKeys[0], myKeys[1]
	accessToken, accessSecret := myKeys[2], myKeys[3]

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)

	return api
}

func publishTweetFromStream(m *nats.Msg) {

	tweet := getStream()

	//curTweet := &tr.Tweet{}
	//curTweet.Text = tweet
	//fmt.Println("twitter-server ", curTweet.Text)


	//tweets := &tr.Twitter{TwitterText:tweet}
	tweets := &tr.Twitter{}
	tweets.TwitterText = tweet

	//for _, t := range tweet {
	//	x = append(tweets.TwitterText, t)
	//
	//}

	data, err := proto.Marshal(tweets)
	if err != nil {
		fmt.Println("Fails here ", err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)


}

func getStream() []string {
	api := auth()

	urlValues := url.Values{}
	urlValues.Set("track", "brexit")
	twitterStream := api.PublicStreamFilter(urlValues)

	var tweets []string

	for t := range twitterStream.C {
		switch v := t.(type) {
		case anaconda.Tweet:
			tweetText := v.Text
			fmt.Println("twitter-server-tweet", tweetText)
			//tweets = append(tweets, tweetText)

			return tweets
		}
	}
	return tweets
}
