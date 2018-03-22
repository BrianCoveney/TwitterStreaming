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

	m := mux.NewRouter()
	m.HandleFunc("/{id}", handleTwitterUser)

	http.ListenAndServe(":3000", m)

}


func handleTwitterUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	myUser := tr.User{Id: vars["id"]}

	myTweet := tr.Tweet{}

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


	go func() {
		msg, err := nc.Request("TwitterByText", nil, 100*time.Millisecond)
		if err == nil && msg != nil {
			receivedTweet := tr.Tweet{}
			err := proto.Unmarshal(msg.Data, &receivedTweet)
			log.Print("!!!!!!!!!!!!!!!!! Hello", receivedTweet)
			if err == nil {
				myTweet = receivedTweet
				log.Print("!!!!!!!!!!!!!!!!! My tweet*", myTweet)
			}
		}
		log.Print("!!!!!!!!!!!!!!!!! Hello nill pointer")

		wg.Done()
	}()



	wg.Wait()


	fmt.Fprintln(w, "Hello ", myUser.Name, " with id ", myUser.Id,
		", the tweet is ", myTweet.Text)

	//go func() {
	//	data, err := proto.Marshal(&myTweet)
	//	if err != nil {
	//		fmt.Println(err)
	//		w.WriteHeader(500)
	//		fmt.Println("Problem with parsing the tweet Id .")
	//		return
	//	}
	//
	//	msg, err := nc.Request("UserNameById", data, 100*time.Millisecond)
	//	if err == nil && msg != nil {
	//		myTwitterWithText := tr.Tweet{}
	//		err := proto.Unmarshal(msg.Data, &myTwitterWithText)
	//		if err == nil {
	//			myTweet = myTwitterWithText
	//			log.Print("My TTTTTTTTTT", myTweet)
	//		}
	//	}
	//	wg.Done()
	//}()
}
