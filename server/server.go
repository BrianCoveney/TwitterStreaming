package main

import (
	"flag"
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"
	"log"
	"os"
	"strings"
	"io/ioutil"
	"net/http"
)

var nc *nats.Conn

func main() {

	uri := os.Getenv("NATS_URI")

	var err error

	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to NATS server " + uri)

	nc.QueueSubscribe("TwitterByText", "TwitterTeller", GetTweetStream)
	select {} // Block forever
}


func readKeys() []string {
	myKeysFile, err := ioutil.ReadFile("my-keys")
	if err != nil {
		fmt.Println("There was a problem with the twitter api keys")
	}
	return strings.Split(string(myKeysFile), "\n")
}

func getHttpClient() *http.Client {
	myKeys := readKeys()
	cKey, cSecret := myKeys[0], myKeys[1]
	aToken, aSecret := myKeys[2], myKeys[3]

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", cKey, "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", cSecret, "Twitter Consumer Secret")
	accessToken := flags.String("access-token", aToken, "Twitter Access Token")
	accessSecret := flags.String("access-secret", aSecret, "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	return httpClient
}


func GetTweetStream(m *nats.Msg) {
	httpClient := getHttpClient()

	client := twitter.NewClient(httpClient)

	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"brexit"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	/*** Type shift tweets to string ***/
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {

		curTweet := tr.Tweet{}
		curTweet.Text = tweet.Text

		data, err := proto.Marshal(&curTweet)
		if err != nil {
			fmt.Println(err)
			return
		}
		nc.Publish(m.Reply, data)
	}

	// Pass the Demux each message or give it the entire Stream.Message
	demux.HandleChan(stream.Messages)

}
