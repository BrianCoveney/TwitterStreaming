package main

import (
	"fmt"
	"log"
	"os"

	tr "github.com/BrianCoveney/TwitterStreaming/twitter-route"
	transport "github.com/BrianCoveney/TwitterStreaming/transport"
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

	m := mux.NewRouter()
	m.HandleFunc("/{id}", handleTwitterUser)

	http.ListenAndServe(":3000", m)

}

func handleTwitterUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	myUser := transport.User{Id: vars["id"]}
	curTweet := tr.Tweet{}
	wg := sync.WaitGroup{}
	wg.Add(2)

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
			myUserWithName := transport.User{}
			err := proto.Unmarshal(msg.Data, &myUserWithName)
			if err == nil {
				myUser = myUserWithName
				log.Print("My User", myUser)
			}
		}
		wg.Done()
	}()


	go func() {
		msg, err := nc.Request("TweetTeller", nil, 100*time.Millisecond)
		if err == nil && msg != nil {
			receivedTweet := tr.Tweet{}
			err := proto.Unmarshal(msg.Data, &receivedTweet)
			if err == nil {
				curTweet = receivedTweet
				log.Print("!!!!!!!!!!!!!!!!! My tweet*", curTweet)
			}
		}
		wg.Done()
	}()



	wg.Wait()


	fmt.Fprintln(w, "Hello ", myUser.Name, " with id ", myUser.Id,
		", the tweet is ", curTweet.Text)



}
