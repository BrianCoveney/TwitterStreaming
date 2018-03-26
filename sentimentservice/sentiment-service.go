package main

import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/cdipaolo/sentiment"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/nats"
	"log"
	"os"
	"sync"
	"time"
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

	nc.QueueSubscribe("SentimentByText", "SentimentTeller", replyWithSentiment)
	select {} // Block forever
}

func replyWithSentiment(m *nats.Msg) {

	myTweet := tr.Tweet{}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		msg, err := nc.Request("TwitterByText", nil, 3000*time.Millisecond)

		if msg == nil || err != nil {
			log.Println("Error {Twitter} on msg nil or err: %v", err)
		} else {
			receivedTweet := tr.Tweet{}
			err := proto.Unmarshal(msg.Data, &receivedTweet)
			if err == nil {
				myTweet = receivedTweet
			}
		}
		//log.Print("My tweet ", myTweet)

		wg.Done()
	}()
	wg.Wait()

	tweetScore := getSentimentScore(myTweet.Text)

	sent := tr.Sentiment{Score: int32(tweetScore)}

	data, err := proto.Marshal(&sent)
	if err != nil {
		fmt.Println(err)
		return
	}
	nc.Publish(m.Reply, data)

	TestNegativeSentenceSentimentShouldReturnZero()
	TestPositiveSentenceSentimentShouldReturnOne()
}

func getSentimentScore(sentence string) uint8 {
	score, err := getSentimentAnalysis(sentence)
	if err != nil {
		log.Println("There was a problem getting the score")
	}
	return score
}

// Return just the score from "github.com/cdipaolo/sentiment"
func getSentimentAnalysis(tweet string) (uint8, error) {
	model, err := sentiment.Restore()
	if err != nil {
		return 0, err
	}
	analysis := model.SentimentAnalysis(tweet, sentiment.English)

	return analysis.Score, nil
}

/** Test methods **/
func TestNegativeSentenceSentimentShouldReturnZero() uint8 {
	score, _ := getSentimentAnalysis("I had an terrible time at a bad game of football")
	fmt.Println("Negative Sentence Sentiment returns ", score)
	return score
}

func TestPositiveSentenceSentimentShouldReturnOne() uint8 {
	score, _ := getSentimentAnalysis("I had an awesome time watching this movie")
	fmt.Println("Positive Sentence Sentiment returns ", score)
	return score
}
