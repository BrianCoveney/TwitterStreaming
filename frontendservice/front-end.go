package main

import (
	"fmt"
	"log"
	"os"


	"github.com/golang/protobuf/proto"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"time"
	"github.com/gorilla/mux"
	"sync"
	"github.com/nats-io/nats"
	"net/http"

	//"github.com/BrianCoveney/TwitterStreaming/twitter-route"
	//tr "github.com/BrianCoveney/TwitterStreaming/twitter-route"
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

	//myTweet := TransportTwitter.Tweet{Text: vars["Text"]}

	//myTweet := TransportTwitter.Tweet{}

	curTime := tr.Time{}


	myUser := tr.User{Id: vars["id"]}
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
			myUserWithName := tr.User{}
			err := proto.Unmarshal(msg.Data, &myUserWithName)
			if err == nil {
				myUser = myUserWithName
				log.Print("My User", myUser)
			}
		}
		wg.Done()
	}()


	//go func() {
	//	data, err := proto.Marshal(&myTweet)
	//	if err != nil || len(myTweet.Text) == 0 {
	//		fmt.Println(err)
	//		w.WriteHeader(500)
	//		fmt.Println("Problem with parsing the tweet Text.")
	//		return
	//	}
	//
	//	msg, err := nc.Request("Tweet test", data, 100*time.Millisecond)
	//	if err == nil && msg != nil {
	//		myTweetWithText := tr.Tweet{}
	//		err := proto.Unmarshal(msg.Data, &myTweetWithText)
	//		if err == nil {
	//			myTweet = myTweetWithText
	//			log.Print("!!!!!!!!!!!!!!!!! My tweet", myTweet)
	//		}
	//	}
	//	wg.Done()
	//}()


	//go func() {
	//	msg, err := nc.Request("Tweet test", nil, 100*time.Millisecond)
	//	if err == nil && msg != nil {
	//		receivedTweet := TransportTwitter.Tweet{}
	//		err := proto.Unmarshal(msg.Data, &receivedTweet)
	//		if err == nil {
	//			myTweet = receivedTweet
	//			log.Print("!!!!!!!!!!!!!!!!! My tweet", myTweet)
	//		}
	//	}
	//	wg.Done()
	//}()


	//go func() {
	//	msg, err := nc.Request("TimeTeller", nil, 100*time.Millisecond)
	//	if err == nil && msg != nil {
	//		receivedTime := TransportTwitter.Tweet{}
	//		err := proto.Unmarshal(msg.Data, &receivedTime)
	//		if err == nil {
	//			curTime = receivedTime
	//			log.Print("!!!!!!!!!!!!!!!!! My tweet", curTime)
	//		}
	//	}
	//	wg.Done()
	//}()

	go func() {
		msg, err := nc.Request("TimeTeller", nil, 100*time.Millisecond)
		if err == nil && msg != nil {
			receivedTime := tr.Time{}
			err := proto.Unmarshal(msg.Data, &receivedTime)
			if err == nil {
				curTime = receivedTime
			}
		}
		wg.Done()
	}()



	wg.Wait()


	fmt.Fprintln(w, "Hello ", myUser.Name, " with id ", myUser.Id, ", the tweet is ", curTime.Time, "day ", curTime.Day)



}
