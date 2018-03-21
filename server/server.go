package main

import (
	"fmt"
	"github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"
	"os"
	"flag"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/oauth1"
	"log"
	"os/signal"
	"syscall"
	"github.com/dghubble/go-twitter/twitter"
)


var nc *nats.Conn
//var client *twitter.Client


func main() {

	uri := os.Getenv("NATS_URI")

	var err error

	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to NATS server " + uri)

	nc.QueueSubscribe("TwitterByText", "TwitterTeller", ReplyWithTweet)
	select {} // Block forever
}


func ReplyWithTweet(m *nats.Msg) {
	curTweet := Transport.Tweet{Text: "Some text"}

	data, err := proto.Marshal(&curTweet)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to anything", m.Reply)
	nc.Publish(m.Reply, data)
}


//func replyWithUserId(m *nats.Msg) {
//
//	log.Print("Log message")
//
//	myTweet := tr.Tweet{}
//	err := proto.Unmarshal(m.Data, &myTweet)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	myTweet.Text = twitterTest[myTweet.Text]
//	data, err := proto.Marshal(&myTweet)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println("Replying to ", m.Reply)
//	nc.Publish(m.Reply, data)
//}
//
//
//
//
func GetTweetStream(m *nats.Msg)  {

	log.Print("Server message")

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
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"brexit"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Type shift tweets to string
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		ts := tweet.Text
		log.Print("Log server message stream", ts)


		myTweet := Transport.Tweet{}
		err := proto.Unmarshal(m.Data, &myTweet)
		if err != nil {
			fmt.Println(err)
			return
		}

		myTweet.Text = ts
		data, err := proto.Marshal(&myTweet)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Replying to ", m.Reply)
		nc.Publish(m.Reply, data)


	}


	// Pass the Demux each message or give it the entire Stream.Message
	for message := range stream.Messages {
		demux.Handle(message)
	}

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")

	stream.Stop()
}



