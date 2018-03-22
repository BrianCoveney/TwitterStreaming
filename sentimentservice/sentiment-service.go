package main

import (
	"github.com/nats-io/nats"

	"os"
	"fmt"
	"github.com/cdipaolo/sentiment"
	"log"
)

var nc *nats.Conn
var model sentiment.Models


func main() {

	uri := os.Getenv("NATS_URI")

	var err error

	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to NATS server " + uri)

	nc.QueueSubscribe("SentimentByText", "SentimentTeller", TestPositiveSentenceSentimentShouldPass)
	select {} // Block forever
}




func TestPositiveSentenceSentimentShouldPass(m *nats.Msg) {

	log.Print("Server message")

	w := []string{
		"I had an awesome time watching this movie",
		"Sometimes I like to say hello to strangers and it's fun",
		"America needs to support the middle class",
		"Harry Potter is a great movie!",
		"The quest for love is a long one, but it ends in happiness",
		"You are a great person",
		"I love the way you can't talk",
		"You are a caring person",
		"I'm quite ambitious, and this job would be a great opportunity for me!",
		"I'm pretty easy-going.",
		"I find it easy to get along with people",
		"I am very hard-working",
		"I'm very methodical and take care over my work",
	}



	for _, sentence := range w {
		s := model.SentimentAnalysis(sentence, sentiment.English)
		if s.Score != uint8(1) {
			fmt.Print("Sentiment of sentence < %v > (returned %v) should be greater than 0.5!\n", sentence, s)
		} else {
			fmt.Print("Sentiment of sentence < %v > is valid.\n\tReturned %v\n", sentence, s)
		}
	}
}

func TestNegativeWordSentimentShouldPass1() {
	w := []string{"not", "resent", "deplorable", "bad", "terrible", "hate", "scary", "terrible", "concerned", "wrong", "rude!!!", "sad", "horrible", "unimpressed", "useless", "offended", "disrespectful"}
	for _, word := range w {
		s := model.SentimentAnalysis(word, sentiment.English)
		if s.Score != uint8(0) {
			fmt.Print("Sentiment of < %v > (returned %v) should be less than 0.5!\n", word, s)
		} else {


			fmt.Print("Sentiment of < %v > valid\n\tReturned %v\n", word, s)
		}
	}
}


func init() {
	var err error

	//model, err = Train()

	model, err = sentiment.Restore()
	if err != nil {
		panic(err.Error())
	}

}

