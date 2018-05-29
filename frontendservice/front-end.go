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
	"io/ioutil"
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

	myTweetSlice := tr.TweetTwitter{}

	mySentiment := tr.Sentiment{}

	data, _ := ioutil.ReadAll(r.Body)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		msg, err := nc.Request("TwitterByText", data, 3000*time.Millisecond)
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
		wg.Done()
	}()

	go func() {
		msg, err := nc.Request("SentimentByText", data, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong. Waiting 2 seconds before retrying:", err)
			return
		}

		receivedSentiment := tr.Sentiment{}
		err = proto.Unmarshal(msg.Data, &receivedSentiment)
		if err != nil {
			log.Print("ERROR ", err)
		}

		mySentiment = receivedSentiment
		log.Print("My Sentiment received: ", mySentiment)
		wg.Done()
	}()

	// Blocks until the above 2 goroutines have completed
	wg.Wait()

	// Create map to hold variables to pass into html template
	m := map[string]interface{}{
		"MyTweets": myTweetSlice.TweetText,
		"MyScore":  mySentiment.Score,
	}

	t, _ := template.ParseFiles("view.html")
	t.Execute(w, m)
}

