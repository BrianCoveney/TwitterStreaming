package main

import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"github.com/ChimeraCoder/anaconda"
	"github.com/nats-io/nats"
	"net/url"
	"io/ioutil"
	"strings"
	"os"
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


func readKeys() [] string  {
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

	curTweet := tr.Tweet{}
	curTweet.Text = tweet

	data, err := proto.Marshal(&curTweet)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)
}


func getStream() string {
	api := auth()

	urlValues := url.Values{}
	urlValues.Set("track", "Brexit")
	twitterStream := api.PublicStreamFilter(urlValues)

	for t := range twitterStream.C {
		switch v := t.(type) {
		case anaconda.Tweet:
			tweetText := v.Text
			return tweetText
		}
	}
	return tweetText
}



