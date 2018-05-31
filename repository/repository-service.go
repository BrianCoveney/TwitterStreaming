package main

import (
	"encoding/json"
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"time"
	"io/ioutil"
	"github.com/nats-io/nats"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"sync"
)

const (
	hosts      = "sentiment-mongodb:27017"
	database   = "sentiment_db"
	username   = ""
	password   = ""
	collection = "sentiment"
)

type Sent struct {
	Score	string `json:"score"`
}

type MongoStore struct {
	session *mgo.Session
}

var mongoStore = MongoStore{}

var nc *nats.Conn

func main() {

	//Create MongoDB session
	session := initialiseMongo()
	mongoStore.session = session

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/sentiment", sentGetHandler).Methods("GET")
	router.HandleFunc("/sentiment", sentPostHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":9090", router))

}

func initialiseMongo() (session *mgo.Session){

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	return

}

func sentGetHandler(w http.ResponseWriter, r *http.Request) {

	col := mongoStore.session.DB(database).C(collection)

	results := []Sent{}
	col.Find(bson.M{"score": bson.RegEx{"", ""}}).All(&results)
	jsonString, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(jsonString))

}

func sentPostHandler(w http.ResponseWriter, r *http.Request) {

	col := mongoStore.session.DB(database).C(collection)

	//Retrieve body from http request
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}


	wg := sync.WaitGroup{}
	wg.Add(1)

	tweetSentiment := tr.Sentiment{}

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

	wg.Wait()


	fmt.Println("repository sentiment", tweetSentiment)

}