package main


import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"
	"os"
	"flag"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/oauth1"
	"log"
	"github.com/dghubble/go-twitter/twitter"
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


func GetTweetStream(m *nats.Msg)  {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "GVfEgw6AQc9T7kJtgYXTGruA3", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "njhflLVqEmpt54NYFkaDL7vfBMaYbUQJ7mst3UyE36LlURsP6T", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "885340272-oidxTehfUKu1KgucgVuvuQGwnffLYjG6Os1QYj0M", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "YUnTVJeYJfcQKw08VrPnlVGKDKCvhpQRz101sxEX9Xy8Z", "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

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


