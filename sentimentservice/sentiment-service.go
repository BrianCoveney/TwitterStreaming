package main

import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/cdipaolo/sentiment"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/nats"
	"log"
	"os"
	"time"
	"sync"
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

	nc.QueueSubscribe("TwitterSentimentByText", "TwitterSentiment", replyWithSentimentForTwitter)
	nc.QueueSubscribe("HackerNewsSentimentByText", "HackerNewsSentiment", replyWithSentimentForHackerNews)

	select {} // Block forever
}

func replyWithSentimentForTwitter(m *nats.Msg) {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		myTweetSlice := tr.TweetTwitter{}
		msg, err := nc.Request("TwitterByText", nil, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong. Waiting 2 seconds before retrying:", err)
			return
		}

		receivedTweetSlice := tr.TweetTwitter{}
		err = proto.Unmarshal(msg.Data, &receivedTweetSlice)
		if err != nil {
			log.Print("ERROR ", err)
		}
		myTweetSlice = receivedTweetSlice

		score := getSentimentScore(myTweetSlice.TweetText)
		twSent := tr.Sentiment{Score: int32(score)}

		data, err := proto.Marshal(&twSent)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Replying to ", m.Reply)
		nc.Publish(m.Reply, data)
		nc.Flush()
		wg.Done()
	}()
	wg.Wait()
}

func replyWithSentimentForHackerNews(m *nats.Msg) {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		myHackerNewsSlice := tr.HackerNews{}
		msg, err := nc.Request("HackerNewsByText", nil, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong with HackerNewsByText. Waiting 2 seconds before retrying:", err)
			return
		}

		receivedHackerNewsSlice := tr.HackerNews{}
		err = proto.Unmarshal(msg.Data, &receivedHackerNewsSlice)
		if err != nil {
			log.Print("ERROR ", err)
		}
		myHackerNewsSlice = receivedHackerNewsSlice

		score := getSentimentScore(myHackerNewsSlice.News)
		nhSent := tr.Sentiment{Score: int32(score)}

		data, err := proto.Marshal(&nhSent)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Replying to ", m.Reply)
		nc.Publish(m.Reply, data)
		nc.Flush()
		wg.Done()
	}()
	wg.Wait()
}


func getSentimentScore(tweets []string) uint8 {
	score, err := getSentimentAnalysis(tweets)
	if err != nil {
		log.Println("There was a problem getting the score")
	}
	return score
}

// Return just the score from "github.com/cdipaolo/sentiment"
func getSentimentAnalysis(tweets []string) (uint8, error) {
	model, err := sentiment.Restore()
	if err != nil {
		return 0, err
	}

	var analysis uint8
	for _, sentence := range tweets {
		analysis := model.SentimentAnalysis(sentence, sentiment.English)
		return analysis.Score, nil
	}

	return analysis, nil
}

