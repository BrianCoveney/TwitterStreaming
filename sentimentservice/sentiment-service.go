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
	log.Print("my_sentiment tweets ", myTweetSlice.TweetText)
	fmt.Println("my_sentiment tweets ", myTweetSlice.TweetText)



	// Test positive returns 1
	//w := []string{
	//	"I had an awesome time watching this movie",
	//	"Sometimes I like to say hello to strangers and it's fun",
	//}


	score := getSentimentScore(myTweetSlice.TweetText)
	fmt.Println("my_sentiment score ", score)

	sent := tr.Sentiment{Score: int32(score)}

	data, err := proto.Marshal(&sent)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)
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

