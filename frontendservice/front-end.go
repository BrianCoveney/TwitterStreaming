package main

import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"html/template"
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

	server := &http.Server{
		Addr:    ":3000",
		Handler: initRoutes(),
	}

	server.ListenAndServe()
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{sentiment}", handleTwitterUser)
	return router
}

func handleTwitterUser(w http.ResponseWriter, r *http.Request) {

	tweetSlice := tr.TweetTwitter{}
	tweetSentiment := tr.Sentiment{}

	hackerNewsSlice := tr.HackerNews{}
	hackerNewsSentiment := tr.Sentiment{}

	//data, _ := ioutil.ReadAll(r.Body)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		msg, err := nc.Request("TwitterByText", nil, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong with TwitterByText. Waiting 2 seconds before retrying:", err)
			return
		}

		receivedTweetSlice := tr.TweetTwitter{}
		err = proto.Unmarshal(msg.Data, &receivedTweetSlice)
		if err != nil {
			log.Print("ERROR ", err)
		}

		tweetSlice = receivedTweetSlice
		wg.Done()
	}()

	go func() {
		msg, err := nc.Request("TwitterSentimentByText", nil, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong with Twitter SentimentByText. Waiting 2 seconds before retrying:", err)
			return
		}

		receivedSentiment := tr.Sentiment{}
		err = proto.Unmarshal(msg.Data, &receivedSentiment)
		if err != nil {
			log.Print("ERROR ", err)
		}

		tweetSentiment = receivedSentiment
		wg.Done()
	}()

	go func() {
		msg, err := nc.Request("HackerNewsByText", nil, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong with HackerNewsByText. Waiting 2 seconds before retrying:", err)
			return
		}

		receivedHackerNewsSlice := tr.HackerNews{}

		//fmt.Println("RECEIVED HN: ", receivedHackerNewsSlice)

		err = proto.Unmarshal(msg.Data, &receivedHackerNewsSlice)
		if err != nil {
			log.Print("ERROR ", err)
		}

		hackerNewsSlice = receivedHackerNewsSlice

		wg.Done()
	}()

	go func() {
		msg, err := nc.Request("HackerNewsSentimentByText", nil, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong with HNews SentimentByText. Waiting 2 seconds before retrying:", err)
			return
		}

		receivedSentiment := tr.Sentiment{}
		err = proto.Unmarshal(msg.Data, &receivedSentiment)
		if err != nil {
			log.Print("ERROR ", err)
		}

		hackerNewsSentiment = receivedSentiment
		wg.Done()
	}()

	// Blocks until the above 3 goroutines have completed
	wg.Wait()

	// Create map to hold variables to pass into html template
	m := map[string]interface{}{
		"Tweets": tweetSlice.TweetText,
		"TweetScore":  tweetSentiment.Score,
		"News":   hackerNewsSlice.News,
		"NewsScore": hackerNewsSentiment.Score,
	}

	t, _ := template.ParseFiles("view.html")
	t.Execute(w, m)
}

