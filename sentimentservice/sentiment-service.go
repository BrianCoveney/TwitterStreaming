package main

import (
	"github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/nats"
	"os"
	"fmt"
	"github.com/cdipaolo/sentiment"
	"log"
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

	log.Print("Sentiment Log message")
	//negSentence := "Your mother is an awful lady"

	posSentence := "I had an awesome time watching this movie"
	posScore := PositiveSentenceSentimentShouldReturn1(posSentence)

	sent := Transport.Sentiment{Score: int32(posScore)}

	data, err := proto.Marshal(&sent)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)

}


func PositiveSentenceSentimentShouldReturn1(sentence string) uint8 {
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
