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


	// We need to increase the timeout to 3 seconds, so our subscriber has a chance to receive the message.
	go func() {
		msg, err := nc.Request("TwitterByText", nil, 3000*time.Millisecond)
		if msg == nil || err != nil  {
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


	go func() {
		msg, err := nc.Request("SentimentByText", nil, 3000*time.Millisecond)
		if msg == nil || err != nil {
			log.Println("Error on msg {Sentiment} nil or err: %v", err)
		} else {
			receivedSentiment := tr.Sentiment{}
			err := proto.Unmarshal(msg.Data, &receivedSentiment)
			if err == nil {
				mySentiment = receivedSentiment
			}
		}
		//log.Print("My Sentiment", mySentiment)
		wg.Done()
	}()


	// Not relevant to this project, but I left this in because it starts the frontend and streaming, with the route
	// http://localhost:3000/<insert_anything>
	// Reason being, I could focus on other parts of the project.
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
				//log.Print("My User", myUser)
			}
		}
		wg.Done()
	}()

	// Blocks until the above 3 goroutines have completed
	wg.Wait()

	if myTweet.Text == "" {
		fmt.Fprintln(w, "Please run 'docker-compose up' again")
	} else {
		fmt.Fprintln(w, "The the tweet is: \n\t ", myTweet.Text,
			"\n\nWith a sentiment score of: \n\t ", mySentiment.Score)
	}

}