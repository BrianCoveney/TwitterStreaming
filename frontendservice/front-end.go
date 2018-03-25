package main

import (
	"fmt"
	"log"
	"os"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"time"
	"github.com/gorilla/mux"
	"sync"
	"github.com/nats-io/nats"
	"net/http"

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
	router.HandleFunc("/{id}", handleTwitterUser)
	return router
}


func handleTwitterUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	myUser := tr.User{Id: vars["id"]}

	myTweet := tr.Tweet{}

	mySentiment := tr.Sentiment{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		data, err := proto.Marshal(&myUser)
		if err != nil || len(myUser.Id) == 0 {
			fmt.Println(err)
			w.WriteHeader(500)
			fmt.Println("Problem with parsing the user Id .")
			return
		}

		msg, err := nc.Request("UserNameById", data, 100*time.Millisecond)
		if err == nil && msg != nil {
			myUserWithName := tr.User{}
			err := proto.Unmarshal(msg.Data, &myUserWithName)
			if err == nil {
				myUser = myUserWithName
				log.Print("My User", myUser)
			}
		}
		wg.Done()
	}()


	// We need to increase the timeout to 3 seconds, so our subscriber has a chance to receive the message.
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


	go func() {
		msg, err := nc.Request("SentimentByText", nil, 3000*time.Millisecond)
		if msg == nil {
			log.Println("Error on msg {Sentiment} nil: %v", err)
		}
		if err != nil {
			log.Println("Error {Sentiment} on err not nil: %v", err)
		}
		if err == nil  {
			receivedSentiment := tr.Sentiment{}
			err := proto.Unmarshal(msg.Data, &receivedSentiment)
			if err == nil {
				mySentiment = receivedSentiment
			}
		}
		log.Print("My Sentiment", mySentiment)

		wg.Done()
	}()


	wg.Wait()

	if myTweet.Text == "" {
		fmt.Fprintln(w, "No tweets since the last page refresh. Try again in one minute")
	} else {
		fmt.Fprintln(w, "The the tweet is: \n\t ", myTweet.Text,
			"\n\nWith a sentiment score of: \n\t ", mySentiment.Score)
	}

}