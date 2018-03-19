package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"
	"github.com/BrianCoveney/TwitterStreaming/transport"
	"time"
	"github.com/gorilla/mux"
	"sync"
	"github.com/nats-io/nats"
	"net/http"
	pb "github.com/BrianCoveney/TwitterStreaming/twitter-route"
	"google.golang.org/grpc"
	"context"
	"io"
	"github.com/BrianCoveney/TwitterStreaming/sentiment"
)

var nc *nats.Conn
var client pb.TwitterRouteClient


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



	myUser := Transport.User{Id: vars["id"]}
	curTime := Transport.Time{}
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
			myUserWithName := Transport.User{}
			err := proto.Unmarshal(msg.Data, &myUserWithName)
			if err == nil {
				myUser = myUserWithName
			}
		}
		wg.Done()
	}()

	go func() {
		msg, err := nc.Request("TimeTeller", nil, 100*time.Millisecond)
		if err == nil && msg != nil {
			receivedTime := Transport.Time{}
			err := proto.Unmarshal(msg.Data, &receivedTime)
			if err == nil {
				curTime = receivedTime
			}
		}
		wg.Done()
	}()


	go func() {
		msg, err := nc.Request("TimeTeller", nil, 100*time.Millisecond)
		if err == nil && msg != nil {

			conn, _ := grpc.Dial("localhost:5253", grpc.WithInsecure())
			client := pb.NewTwitterRouteClient(conn)

			params := &pb.Params{
				Track:         []string{"Bitcoin"},
				Language:      []string{"en"},
				StallWarnings: false,
				Maxcount:      100,
			}

			stream, _ := client.GetTweets(context.Background(), params)

			// Wait for SIGINT and SIGTERM (HIT CTRL-C)
			ch := make(chan os.Signal)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			log.Println(<-ch)

			fmt.Println("Stopping Stream...")


			for {
				tweet, err := stream.Recv()
				if err == io.EOF {
					break
				}

				score, _ := sentiment.Run(tweet.Text)

				tweet.Score = int32(score)
			}


		}
		wg.Done()
	}()


	wg.Wait()


	fmt.Fprintln(w, "Hello ", myUser.Name, " with id ", myUser.Id, ", the time is ")

}
