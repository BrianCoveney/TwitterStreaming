package main

import (
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/nats"
	"os"
	"fmt"
	"github.com/cdipaolo/sentiment"
	"log"
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

	nc.QueueSubscribe("SentimentByText", "SentimentTeller", replyWithSentiment)
	select {} // Block forever
}


func replyWithSentiment(m *nats.Msg) {

	myTweet := tr.Tweet{}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		msg, err := nc.Request("TwitterByText", nil, 3000*time.Millisecond)
		if msg == nil {
			log.Println("Error {Twitter} on msg nil: %v", err)
		}
		if err != nil {
			log.Println("Error {Twitter} on err not nil: %v", err)
		}
		if err == nil  {
			receivedTweet := tr.Tweet{}
			err := proto.Unmarshal(msg.Data, &receivedTweet)
			if err == nil {
				myTweet = receivedTweet
			}
		}
		log.Print("My tweet ", myTweet)

		wg.Done()
	}()
	wg.Wait()

	fmt.Println("FROM SENTTTTTTTTTTTTTTt ", myTweet)



	//posSentence := "I had an awesome time watching this movie"
	//posScore := GetSentimentScore(posSentence)

	tweetScore := GetSentimentScore(myTweet.Text)

	sent := tr.Sentiment{Score: int32(tweetScore)}

	data, err := proto.Marshal(&sent)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)

}


func GetSentimentScore(sentence string) uint8 {
	score, err := RunSentimentAnalysis(sentence)
	if err != nil {
		log.Println("There was a problem getting the score")
	}
	return score
}


//func NegativeSentenceSentimentShouldReturn0(sentence string) uint8 {
//
//	score, err := RunSentimentAnalysis(sentence)
//	if err != nil {
//		fmt.Println("There was a problem getting the score")
//	}
//	return score
//}


// Return just the score from "github.com/cdipaolo/sentiment"
func RunSentimentAnalysis(tweet string) (uint8, error) {
	model, err := sentiment.Restore()
	if err != nil {
		return 0, err
	}
	analysis := model.SentimentAnalysis(tweet, sentiment.English)

	return analysis.Score, nil
}
